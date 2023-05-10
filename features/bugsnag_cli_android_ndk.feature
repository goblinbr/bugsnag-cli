Feature: Bugsnag CLI Android NDK behavior
  Scenario: Starting bugsnag-cli upload android-ndk on mac without an API Key
    When I run bugsnag-cli with upload android-ndk
    Then I should see the API Key error

  Scenario: Starting bugsnag-cli upload android-ndk on mac without a path
    When I run bugsnag-cli with upload android-ndk --api-key=1234567890ABCDEF1234567890ABCDEF
    Then I should see the missing path error

  Scenario: Starting bugsnag-cli upload android-ndk with an invalid path
    When I run bugsnag-cli with upload android-ndk --api-key=1234567890ABCDEF1234567890ABCDEF /path/to/no/file
    Then I should see the no such file or directory error
    