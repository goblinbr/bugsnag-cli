package upload_testing

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/stretchr/testify/assert"
)

func GetBasePath() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	sampleRegexp := regexp.MustCompile(`/[^/]*/[^/]*$`)
	basePath := sampleRegexp.ReplaceAllString(path, "")

	return basePath
}

func TestReadElfBuildId(t *testing.T) {
	t.Log("Testing getting a build ID from an ELF file")
	results, err := upload.ReadElfBuildId(filepath.Join(GetBasePath()+ "/test/testdata/dart/app-debug-info/app.android-arm64.symbols"))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, results, "07cc131ca803c124e93268ce19322737", "Build Ids should be the same")
}

func TestDwarfDumpUuid(t *testing.T) {
	t.Log("Testing getting a build ID from a Dwarf file")
	results, err := upload.DwarfDumpUuid(GetBasePath()+"/test/testdata/dart/app-debug-info/app.ios-arm64.symbols", GetBasePath()+"/test/testdata/dart/build/ios/iphoneos/Runner.app/Frameworks/App.framework/App")

	if err != nil {
		log.Println(err)
	}

	assert.Equal(t, results, "E30C1BE5-DEB6-373C-98B4-52D827B7FF0D", "UUID should match")
}

func TestGetIosAppPath(t *testing.T) {
	t.Log("Testing getting the IOS app path from a given symbols path")
	results, err := upload.GetIosAppPath(GetBasePath() + "/test/testdata/dart/app-debug-info/app.android-arm64.symbols")

	if err != nil {
		log.Println(err)
	}

	assert.Equal(t, results, GetBasePath()+"/test/testdata/dart/build/ios/iphoneos/Runner.app/Frameworks/App.framework/App", "They should match")
}
