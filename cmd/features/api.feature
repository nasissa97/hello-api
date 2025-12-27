Feature: Translate API
  Users should be able to submit a word to translate

  @smoke-test
  Scenario: Translation
    Given the word "hello"
    When I translate it to "german"
    Then the response should be "hallo"

  @smoke-test
  Scenario: Translation unknown
    Given the word "goodbye"
    When I translate it to "german"
    Then the response should be ""
  Scenario: Translation unknown
    Given the word "hello"
    When I translate it to "bulgarian"
    Then the response should be ""
