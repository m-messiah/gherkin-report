@Login @Regression @RegressionPhase1
Feature: Login

  @ignore
    @LoginFieldsViaHeaderBar
  Scenario Outline: Checking login fields via the header bar [<localization>]
    Given old registration is enabled
    And I open homepage for localization <localization>
    When I click login link
    Then I see login popup
    And I don't see clear buttons
    And I check translations on login pop-up
    And I check that T&C link open in new tab

    @Default
    Examples:
      | localization |
      | EN_UK    |

    Examples:
      | localization    |
      | UK              |
      | EN_UK           |
      | EN_US           |

  @LoginFieldsViaHeaderBarNewRegistration
  Scenario Outline: Checking login fields via the header bar for new registration flow [<localization>]
    Given new registration is enabled
    And I open homepage for localization <localization>
    When I click login link
    Then I see login popup
    And I don't see clear buttons
    And I check translations on login pop-up with new registration
    And I check that T&C link open in new tab

    @Default
    Examples:
      | localization  |
      | EN_UK         |

    Examples:
      | localization  |
      | EN_UK         |
