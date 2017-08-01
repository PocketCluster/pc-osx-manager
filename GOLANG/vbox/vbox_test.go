package vbox

import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestAppVersion(t *testing.T) {
    if AppVersion <= 4003000 {
        t.Error("AppVersion below 4.3: ", AppVersion)
    }
}

func TestGetRevision(t *testing.T) {
    revision, err := GetRevision()
    if err != nil {
        t.Fatal(err)
    }
    if revision <= 100000 {
        t.Error("Revision below 100000: ", revision)
    }
}

func TestComposeMachineFilename(t *testing.T) {
    tempPath := os.TempDir()
    vmpath, err := ComposeMachineFilename("TestVM", "", os.TempDir())
    if err != nil {
        t.Fatal(err)
    }
    if vmpath != filepath.Join(tempPath, "TestVM.vbox") {
        t.Error("Wrong VM filename when given baseFolder: ", vmpath)
    }

    vmpath, err = ComposeMachineFilename("TestVM", "", "")
    if err != nil {
        t.Fatal(err)
    }
    if !strings.Contains(vmpath, "VirtualBox") {
        t.Error("VM filename without baseFolder doesn't have VirtualBox: ", vmpath)
    }
    if !strings.Contains(vmpath, "TestVM.vbox") {
        t.Error("VM filename without baseFolder doesn't have TestVM.vbox: ", vmpath)
    }
}
