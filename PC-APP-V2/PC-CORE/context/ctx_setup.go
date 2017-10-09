package context

import (
    "os"
    "path/filepath"

    "github.com/pkg/errors"
)

func SetupBasePath() error {
    // when this returns error, it's from PosixHomeDirectory and there is nothing you can fix it. just return error
    dataDir, err := SharedHostContext().ApplicationUserDataDirectory()
    if err != nil {
        return errors.WithStack(err)
    }
    if _, err = os.Stat(dataDir); os.IsNotExist(err) {
        err = os.MkdirAll(dataDir, os.ModeDir|0700)
        if err != nil {
            return errors.WithStack(err)
        }
    }

    // registry configuration
    repoPath, err := SharedHostContext().ApplicationRepositoryDirectory()
    if err != nil {
        return errors.WithStack(err)
    }
    if _, err := os.Stat(repoPath); os.IsNotExist(err) {
        err = os.MkdirAll(filepath.Join(repoPath, "docker/registry/v2/repositories"), os.ModeDir|0700)
        if err != nil {
            return errors.WithStack(err)
        }
        err = os.MkdirAll(filepath.Join(repoPath, "docker/registry/v2/blobs"), os.ModeDir|0700)
        if err != nil {
            return errors.WithStack(err)
        }
    }

    //etcd configuration
    storagePath, err := SharedHostContext().ApplicationStorageDirectory()
    if err != nil {
        return errors.WithStack(err)
    }
    if _, err := os.Stat(storagePath); os.IsNotExist(err) {
        err = os.MkdirAll(storagePath, os.ModeDir|0700)
        if err != nil {
            return errors.WithStack(err)
        }
    }

    // virtual machine configuration
    vmPath, err := SharedHostContext().ApplicationVirtualMachineDirectory()
    if err != nil {
        return errors.WithStack(err)
    }
    if _, err := os.Stat(vmPath); os.IsNotExist(err) {
        err = os.MkdirAll(vmPath, os.ModeDir|0700)
        if err != nil {
            return errors.WithStack(err)
        }
    }

    // pocket core data directory
    cdata, err := SharedHostContext().ApplicationPocketCoreDataDirectory()
    if err != nil {
        return errors.WithStack(err)
    }
    if _, err := os.Stat(cdata); os.IsNotExist(err) {
        err = os.MkdirAll(cdata, os.ModeDir|0700)
        if err != nil {
            return errors.WithStack(err)
        }
    }

    // pocket core input directory
    cinput, err := SharedHostContext().ApplicationPocketCoreInputDirectory()
    if err != nil {
        return errors.WithStack(err)
    }
    if _, err := os.Stat(cinput); os.IsNotExist(err) {
        err = os.MkdirAll(cinput, os.ModeDir|0700)
        if err != nil {
            return errors.WithStack(err)
        }
    }

    return nil
}
