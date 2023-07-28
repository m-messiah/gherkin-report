@SmokeLogin @SmokeTest @no_money @SmokePhase1
Feature: Login Smoketest

  @PROMO @BETA @PROD @DEV
    @SmokeLoginHeader @mobile @tablet @default_breakpoint
  Scenario Outline: To check header in logged in state [<localization>]
    Given I have a test user for <localization> localization
    And I open homepage for localization <localization>
    And I accept cookie policy
    When I login via header bar
    Then I check header layout after login for compliance with the HeaderAfterLogin.gspec file:
      | MOBILE_PORTRAIT | TABLET_PORTRAIT | TABLET_LANDSCAPE |

    Examples:
      | localization |
      | EN_UK        |
