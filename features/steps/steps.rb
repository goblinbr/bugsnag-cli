require 'rbconfig'
require 'etc'

os = RbConfig::CONFIG['host_os']
arch = RbConfig::CONFIG['host_cpu']

case
when os.downcase.include?('windows_nt'), ENV['WSL_DISTRO_NAME'] != nil
  os = 'windows'
  binary = 'bugsnag-cli.exe'
when os.downcase.include?('linux')
  os = 'linux'
  binary = 'bugsnag-cli'
when os.downcase.include?('darwin')
  os = 'macos'
  binary = 'bugsnag-cli'
end

When('I run bugsnag-cli') do
  @output = `bin/#{arch}-#{os}-#{binary} 2>&1`
end

When(/^I run bugsnag-cli with (.*)$/) do |flags|
  @output = `bin/#{arch}-#{os}-#{binary} #{flags} 2>&1`
  puts @output
end

Then('I should see the help banner') do
  Maze.check.include(run_output, "Usage: #{arch}-#{os}-#{binary} <command>")
end

Then('I should see the API Key error') do
  Maze.check.include(run_output, "[ERROR] missing api key, please specify using `--api-key`")
end

Then('I should see the Project Root error') do
  Maze.check.include(run_output, "[ERROR] --project-root is required when uploading dSYMs from a directory that is not an Xcode project or workspace")
end


Then('I should see the missing path error') do
  Maze.check.include(run_output, "error: expected \"<path>\"")
end

Then('I should see the missing app version error') do
  Maze.check.include(run_output, "[ERROR] missing app version, please specify using `--app-version`")
end

Then('I should see the no such file or directory error') do
  Maze.check.include(run_output, "error: <path>: stat /path/to/no/file: no such file or directory")
end

Then('the sourcemap is valid for the Proguard Build API') do
  steps %(
    Then the sourcemap is valid for the Android Build API
  )
end

Then('the sourcemap is valid for the NDK Build API') do
  steps %(
    Then the sourcemap is valid for the Android Build API
  )
end

Then('the sourcemap is valid for the Dart Build API') do
  steps %(
    And the sourcemap payload field "apiKey" equals "#{$api_key}"
    And the sourcemap payload field "buildId" is not null
  )
end

Then('the sourcemap is valid for the React Native Build API') do
  steps %(
    And the sourcemap payload field "apiKey" equals "#{$api_key}"
    And the sourcemap payload field "appVersion" is not null
  )
end

Then('the sourcemap is valid for the dSYM Build API') do
  steps %(
    And the sourcemap payload field "apiKey" equals "#{$api_key}"
    And the sourcemap payload field "appVersion" is not null
  )
end

Then('the sourcemap is valid for the Android Build API') do
  steps %(
    And the sourcemap payload field "apiKey" equals "#{$api_key}"
    And the sourcemap payload field "appId" is not null
  )
end

Then('the build is valid for the Builds API') do
  steps %(
    And the build payload field "apiKey" equals "#{$api_key}"
    And the build payload field "appVersion" is not null
  )
end

Then('the sourcemaps Content-Type header is valid multipart form-data') do
  expected = /^multipart\/form-data; boundary=([^;]+)/
  actual = Maze::Server.sourcemaps.current[:request]['content-type']
  Maze.check.match(expected, actual)
end

Then('{string} should be used as {string}') do |value, field|
  Maze.check.include(run_output, "Using #{value} as #{field} from")
end

def get_version_number(file_path)
  package_version = nil

  file_content = File.read(file_path)

  file_content.each_line do |line|
    if line =~ /\bpackage_version\s*=\s*(['"])(.*?)\1/
      package_version = $2
      break
    end
  end

  package_version
end

Then(/^the version number should match the version set in main\.go$/) do
  version_number = get_version_number "main.go"
  Maze.check.include(run_output, version_number)
end

And(/^I wait for the build to succeed$/) do
  Maze.check.not_include(run_output, "Error 1")
end

When(/^I make the "([^"]*)"$/) do |arg|
  @output = `make #{arg} 2>&1`
  puts @output
end
