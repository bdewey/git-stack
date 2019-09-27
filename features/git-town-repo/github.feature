Feature: git-repo when origin is on GitHub

  Scenario Outline: result
    Given my repo's remote origin is <ORIGIN>
    And I have "open" installed
    When I run `git-town repo`
    Then I see the GitHub homepage of the "Originate/git-town" repository in my browser

    Examples:
      | ORIGIN                                      |
      | http://github.com/bdewey/git-stack.git    |
      | http://github.com/bdewey/git-stack        |
      | https://github.com/bdewey/git-stack.git   |
      | https://github.com/bdewey/git-stack       |
      | git@github.com:Originate/git-town.git       |
      | git@github.com:Originate/git-town           |
      | ssh://git@github.com/bdewey/git-stack.git |
      | ssh://git@github.com/bdewey/git-stack     |
