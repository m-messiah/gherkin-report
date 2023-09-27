# gherkin report

The report util is used to collect features and generate the table of what scenarios are enabled and what environments they use.

Checks two directories: `./features` and `./smoketest_features`. (see examples in there)
Generates two csv files `./features_result.csv` and `./smoketest_features_result.csv`.

## How to run

1. Change the current path to the `resources` directory that contains the `./features` and `./smoketest_features` directories.
2. Run the tool with `go run github.com/m-messiah/gherkin-report@latest` or with `gherkin-report` CLI installed

## Customisation

Currently, the tool looks for features and creates report files in the same directory where it is invoked.
It also has a hardcoded list of environments and emulators, which should be extended in code if needed.

## Build from source

1. `go build`
2. `go install`
