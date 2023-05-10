package android

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// Objcopy - Processes files using objcopy
func Objcopy(objcopyPath string, file string, outputPath string) (string, error) {

	objcopyLocation, err := exec.LookPath(objcopyPath)

	if err != nil {
		return "", err
	}

	outputFile := filepath.Join(outputPath, filepath.Base(file))
	outputFile = strings.ReplaceAll(outputFile, filepath.Ext(outputFile), ".so.sym")

	cmd := exec.Command(objcopyLocation, "--compress-debug-sections=zlib", "--only-keep-debug", file, outputFile)

	_, err = cmd.CombinedOutput()

	return outputFile, nil
}
