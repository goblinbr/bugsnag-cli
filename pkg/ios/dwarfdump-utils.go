package ios

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// DwarfInfo stores the UUID, architecture and name of a dwarf file
type DwarfInfo struct {
	UUID string
	Arch string
	Name string
}

// isDwarfDumpInstalled checks if dwarfdump is installed by checking if there is a path returned for it
func isDwarfDumpInstalled() bool {
	return utils.LocationOf(utils.DWARFDUMP) != ""
}

// GetDsymsForUpload returns information on the valid dSYM files found in a given path
func GetDsymsForUpload(path string) (*[]*DwarfInfo, error) {
	filesFound, _ := os.ReadDir(path)
	var tempDir string
	var dsymFiles []*DwarfInfo

	switch len(filesFound) {
	case 0:
		return nil, errors.Errorf("No dSYM files found in expected location '%s'", path)
	default:
		if isDwarfDumpInstalled() {
			for _, file := range filesFound {
				if strings.HasSuffix(file.Name(), ".zip") {
					log.Info("Attempting to unzip " + file.Name() + " before proceeding to upload")
					tempDir, _ = utils.ExtractFile(filepath.Join(path, file.Name()), "dsym")

					if tempDir != "" {
						log.Info("Unzipped " + file.Name() + " to " + tempDir + " for uploading")
						path = tempDir
					} else {
						log.Warn("Could not unzip " + file.Name() + " to a temporary directory, skipping")
						// Silently remove the temp dir if one was created before continuing
						removeTempDir(tempDir)
						continue
					}
				}
				dsymFiles = getDwarfFileInfo(path, file.Name())
			}
			removeTempDir(tempDir)
		} else {
			return nil, errors.New("Unable to locate dwarfdump on this system.")
		}
	}

	return &dsymFiles, nil
}

// getDwarfFileInfo parses dwarfdump output to easier to manage/parsable DwarfInfo structs
func getDwarfFileInfo(path, fileName string) []*DwarfInfo {
	var dwarfInfo []*DwarfInfo

	cmd := exec.Command(utils.DWARFDUMP, "-u", strings.TrimSuffix(fileName, ".zip"))
	cmd.Dir = path

	output, _ := cmd.Output()
	if len(output) > 0 {
		outputStr := string(output)

		outputStr = strings.TrimSuffix(outputStr, "\n")
		outputStr = strings.ReplaceAll(outputStr, "(", "")
		outputStr = strings.ReplaceAll(outputStr, ")", "")

		outputSlice := strings.Split(outputStr, "\n")

		for _, str := range outputSlice {
			if strings.Contains(str, "UUID: ") {
				rawDwarfInfo := strings.Split(str, " ")
				if len(rawDwarfInfo) == 4 {
					dwarf := &DwarfInfo{}
					dwarf.UUID = rawDwarfInfo[1]
					dwarf.Arch = rawDwarfInfo[2]
					dwarf.Name = rawDwarfInfo[3]
					dwarfInfo = append(dwarfInfo, dwarf)
				}
			}
		}
	} else {
		log.Info("Skipping file without UUID: " + fileName)
	}

	return dwarfInfo
}

// removeTempDir removes a temporary directory and disregards any errors
func removeTempDir(tempDir string) {
	_ = os.RemoveAll(tempDir)
}
