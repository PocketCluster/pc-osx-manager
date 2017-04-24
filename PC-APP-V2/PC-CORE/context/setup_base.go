package context

import (
    "os"

    "github.com/pkg/errors"
)

func SetupBasePath() error {
    var (
        dataDir string = ""
        err error = nil
    )
    dataDir, err = SharedHostContext().ApplicationUserDataDirectory()
    if err != nil {
        return errors.WithStack(err)
    }

    // check if the path exists and make it if absent
    _, err = os.Stat(dataDir)
    if os.IsNotExist(err) {
        err = os.MkdirAll(dataDir, os.ModeDir | 0700)
    }
    if err != nil {
        return errors.WithStack(err)
    }

    // TODO create container repository
    return nil
}
