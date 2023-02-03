package android_testing

import (
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/stretchr/testify/assert"
)

func TestBuildObjcopyPath(t *testing.T) {
	t.Log("Testing building Objcopy path")
	results, err := android.BuildObjcopyPath("/opt/homebrew/share/android-commandlinetools/ndk/24.0.8215888")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "/opt/homebrew/share/android-commandlinetools/ndk/24.0.8215888/toolchains/llvm/prebuilt/darwin-x86_64/bin/llvm-objcopy", results, "The paths should match")
}
