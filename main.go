package main

import (
	"github.com/alecthomas/kong"
	"github.com/bugsnag/bugsnag-cli/pkg/build"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
)

func main() {
	var commands struct {
		UploadAPIRootUrl  string `help:"Bugsnag On-Premise upload server URL. Can contain port number" default:"https://upload.bugsnag.com"`
		BuildApiRootUrl   string `help:"Bugsnag On-Premise build server URL. Can contain port number" default:"https://build.bugsnag.com"`
		Port              int    `help:"Port number for the upload server" default:"443"`
		ApiKey            string `help:"(required) Bugsnag integration API key for this application"`
		FailOnUploadError bool   `help:"Stops the upload when a mapping file fails to upload to Bugsnag successfully" default:false`
		AppVersion        string `help:"The version of the application."`
		AppVersionCode    string `help:"The version code for the application (Android only)."`
		AppBundleVersion  string `help:"The bundle version for the application (iOS only)."`
		Upload            struct {

			// shared options
			Overwrite bool `help:"Whether to overwrite any existing symbol file with a matching ID"`
			Timeout   int  `help:"Number of seconds to wait before failing an upload request" default:"300"`
			Retries   int  `help:"Number of retry attempts before failing an upload request" default:"0"`

			// required options
			AndroidAab      upload.AndroidAabMapping      `cmd:"" help:"Process and upload application bundle files for Android"`
			All             upload.DiscoverAndUploadAny   `cmd:"" help:"Upload any symbol/mapping files"`
			AndroidNdk      upload.AndroidNdkMapping      `cmd:"" help:"Process and upload Proguard mapping files for Android"`
			AndroidProguard upload.AndroidProguardMapping `cmd:"" help:"Process and upload NDK symbol files for Android"`
			DartSymbol      upload.DartSymbol             `cmd:"" help:"Process and upload symbol files for Flutter" name:"dart"`
		} `cmd:"" help:"Upload symbol/mapping files"`
		CreateBuild build.CreateBuild `cmd:"" help:"Provide extra information whenever you build, release, or deploy your application"`
	}

	// If running without any extra arguments, default to the --help flag
	// https://github.com/alecthomas/kong/issues/33#issuecomment-1207365879
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	ctx := kong.Parse(&commands)

	// Check if we have an apiKey in the request
	if commands.ApiKey == "" {
		log.Error("no API key provided", 1)
	}

	// Build connection URI
	endpoint, err := utils.BuildEndpointUrl(commands.UploadAPIRootUrl, commands.Port)

	if err != nil {
		log.Error("Failed to build upload url: "+err.Error(), 1)
	}

	switch ctx.Command() {

	// Upload command
	case "upload all <path>":
		log.Info("Uploading files to: " + endpoint)

		err := upload.All(
			commands.Upload.All.Path,
			commands.Upload.All.UploadOptions,
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.ApiKey,
			commands.FailOnUploadError)

		if err != nil {
			log.Error(err.Error(), 1)
		}

		log.Success("Upload(s) completed")

	case "upload android-ndk <path>":
		endpoint = endpoint + "/ndk-symbol"
		log.Info("Uploading files to: " + endpoint)
		err := upload.ProcessAndroidNDK(
			commands.Upload.AndroidNdk.Path,
			commands.Upload.AndroidNdk.AndroidNdkRoot,
			commands.Upload.AndroidNdk.AppManifestPath,
			commands.Upload.AndroidNdk.Configuration,
			commands.Upload.AndroidNdk.ProjectRoot,
			commands.Upload.AndroidNdk.VersionCode,
			commands.Upload.AndroidNdk.VersionName,
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.ApiKey,
			commands.FailOnUploadError)

		if err != nil {
			log.Error(err.Error(), 1)
		}

		log.Success("Upload(s) completed")

	case "upload android-proguard <path>":
		endpoint = endpoint + "/ndk-symbol"
		log.Info("Uploading files to: " + endpoint)
		err := upload.ProcessAndroidProguard(
			commands.Upload.AndroidProguard.Path,
			commands.Upload.AndroidProguard.ApplicationId,
			commands.Upload.AndroidProguard.AppManifestPath,
			commands.Upload.AndroidProguard.BuildUuid,
			commands.Upload.AndroidProguard.Configuration,
			commands.Upload.AndroidProguard.VersionCode,
			commands.Upload.AndroidProguard.VersionName,
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.ApiKey,
			commands.FailOnUploadError,
			commands.Upload.AndroidProguard.DryRun)

		if err != nil {
			log.Error(err.Error(), 1)
		}

		log.Success("Upload(s) completed")
	case "upload dart <path>":
		endpoint = endpoint + "/dart-symbol"
		log.Info("Uploading files to: " + endpoint)
		err := upload.Dart(commands.Upload.DartSymbol.Path,
			commands.AppVersion,
			commands.AppVersionCode,
			commands.AppBundleVersion,
			commands.Upload.DartSymbol.IosAppPath,
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.ApiKey,
			commands.FailOnUploadError)

		if err != nil {
			log.Error(err.Error(), 1)
		}

		log.Success("Upload(s) completed")

	case "create-build":
		// Build connection URI
		endpoint, err := utils.BuildEndpointUrl(commands.BuildApiRootUrl, commands.Port)

		if err != nil {
			log.Error("Failed to build upload url: "+err.Error(), 1)
		}

		log.Info("Creating build on: " + endpoint)
		buildUploadError := build.ProcessBuildRequest(commands.ApiKey,
			commands.CreateBuild.BuilderName,
			commands.CreateBuild.ReleaseStage,
			commands.CreateBuild.Provider,
			commands.CreateBuild.Repository,
			commands.CreateBuild.Revision,
			commands.AppVersion,
			commands.AppVersionCode,
			commands.AppBundleVersion,
			commands.CreateBuild.Metadata,
			endpoint)
		if buildUploadError != nil {
			log.Error(buildUploadError.Error(), 1)
		}

		log.Success("Build created")
	default:
		println(ctx.Command())
	}
}
