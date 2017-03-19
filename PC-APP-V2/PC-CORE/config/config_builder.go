package config

import (
    "database/sql"
    "os"
    "path/filepath"

    _ "github.com/mattn/go-sqlite3"
    "github.com/pkg/errors"
    "github.com/gravitational/teleport/lib/defaults"
    "github.com/stkim1/pc-core/context"
)

func SetupBaseConfigPath(ctx context.HostContext) error {
    dataDir, err := ctx.ApplicationUserDataDirectory()
    if err != nil {
        return errors.Wrap(err, "Failed to determine hostname")
    }
    // check if the path exists and make it if absent
    if _, err := os.Stat(dataDir); err != nil {
        if os.IsNotExist(err) {
            os.MkdirAll(dataDir, os.ModeDir | 0700);
        }
    }

    // TODO create container repository
    return nil
}

func OpenStorageInstance(ctx context.HostContext) (*sql.DB, error) {
    // This is a good place to migrate & load encryption plugin

    // data directory
    dataDir, err := ctx.ApplicationUserDataDirectory()
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // database file path
    dbPath := filepath.Join(dataDir, defaults.CoreKeysSqliteFile)
    path, err := filepath.Abs(dbPath)
    if err != nil {
        return nil, errors.Wrap(err, "Failed to convert path")
    }

    // check if path is ok to use
    dir := filepath.Dir(path)
    s, err := os.Stat(dir)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    if !s.IsDir() {
        return nil, errors.Errorf("Path '%v' should be a valid directory", dir)
    }

    // create database
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return db, nil
}

func CloseStorageInstance(db *sql.DB) error {
    if db != nil {
        return errors.Errorf("Unable to close null db")
    }
    return errors.WithStack(db.Close())
}