package android

import (
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"path/filepath"
	"strings"
)

func UploadAndroidNdk(
	fileList []string,
	apiKey string,
	applicationId string,
	versionName string,
	versionCode string,
	projectRoot string,
	overwrite bool,
	endpoint string,
	timeout int,
	retries int,
	dryRun bool,
) error {
	fileFieldData := make(map[string]string)

	numberOfFiles := len(fileList)

	if numberOfFiles < 1 {
		log.Info("No NDK files found to process")
		return nil
	}

	for _, file := range fileList {
		uploadOptions, err := utils.BuildAndroidNDKUploadOptions(apiKey, applicationId, versionName, versionCode, projectRoot, filepath.Base(file), overwrite)

		if err != nil {
			return err
		}

		fileFieldData["soFile"] = file

		err = server.ProcessFileRequest(endpoint+"/ndk-symbol", uploadOptions, fileFieldData, timeout, retries, file, dryRun)

		if err != nil {
			if strings.Contains(err.Error(), "409") && strings.Contains(err.Error(), "duplicate") {
				log.Warn("Duplicate file detected, skipping upload of " + filepath.Base(file))
			} else {
				return err
			}
		} else {
			log.Success("Uploaded " + filepath.Base(file))
		}
	}

	return nil
}
