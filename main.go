package main

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	gherkin "github.com/cucumber/gherkin/go/v26"
	messages "github.com/cucumber/messages/go/v21"
)

const (
	ScenarioOutline = "Scenario Outline"
	defaultEmulator = "ft-emulator-MOCK"
	ignoreTag       = "@ignore"
)

var (
	emulators = map[string]string{
		"@ft-emulator_dev_stable": "ft-emulator-DEV",
		"@ft-emulator_beta":       "ft-emulator-BETA",
	}
	smokeTestEnvironments = map[string]struct{}{
		"@PROD":  {},
		"@BETA":  {},
		"@PROMO": {},
		"@DEV":   {},
	}
)

func fillFeature(filename, testType string) [][]string {
	var result [][]string
	f, err := os.Open(filename)
	if err != nil {
		return result
	}
	defer f.Close()
	gherkinDocument, err := gherkin.ParseGherkinDocument(f, (&messages.Incrementing{}).NewId)
	if gherkinDocument == nil {
		return result
	}
	feature := gherkinDocument.Feature
	for _, tag := range feature.Tags {
		if tag != nil && tag.Name == ignoreTag {
			fmt.Fprintf(os.Stderr, "feature %q is ignored\n", feature.Name)
			return result
		}
	}
	for _, child := range feature.Children {
		if child.Scenario == nil {
			continue
		}
		if child.Scenario.Keyword != ScenarioOutline {
			continue
		}
		var env string
		var ignore bool
		switch testType {
		case "features":
			env, ignore = findFeatureEnv(child.Scenario.Tags)
		case "smoketest_features":
			env, ignore = findSmoketestFeatureEnv(child.Scenario.Tags)
		default:
			continue
		}
		if !ignore {
			result = append(result, []string{feature.Name, stripScenario(child.Scenario.Name), env})
		}
	}
	return result
}

func stripScenario(name string) string {
	if i := strings.IndexRune(name, '['); i > -1 {
		name = name[:i]
	}
	return strings.TrimSpace(name)
}

func findFeatureEnv(tags []*messages.Tag) (string, bool) {
	result := defaultEmulator
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		if strings.HasPrefix(tag.Name, "@ft-emulator") {
			if value, ok := emulators[tag.Name]; ok {
				result = value
			}
		}
		if tag.Name == ignoreTag {
			return result, true
		}
	}
	return result, false
}

func findSmoketestFeatureEnv(tags []*messages.Tag) (string, bool) {
	var envs []string
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		if _, ok := smokeTestEnvironments[tag.Name]; ok {
			envs = append(envs, tag.Name)
		}
		if tag.Name == ignoreTag {
			return "", true
		}
	}
	sort.Strings(envs)
	return strings.Join(envs, " "), false
}

func fillAllFeatures(testType string) [][]string {
	var result [][]string
	filepath.WalkDir(testType, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		result = append(result, fillFeature(path, testType)...)
		return nil
	})
	sort.Slice(result, func(i, j int) bool {
		return result[i][0] < result[j][0]
	})
	return result
}

func parseFeatures(testType string) error {
	resultFile, err := os.Create(testType + "_result.csv")
	if err != nil {
		return err
	}
	defer resultFile.Close()
	output := csv.NewWriter(resultFile)
	output.Write([]string{"feature", "scenario", "env"})
	output.WriteAll(fillAllFeatures(testType))
	return nil
}

func main() {
	for _, testType := range []string{"features", "smoketest_features"} {
		err := parseFeatures(testType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "got an error for %s: %s", testType, err)
		}
	}
}
