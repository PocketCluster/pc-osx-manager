package model

import (
    "database/sql"
    "os"
    "path/filepath"
    "sync"

    _ "github.com/mattn/go-sqlite3"
    "github.com/cloudflare/cfssl/certdb"
    "github.com/jinzhu/gorm"
    "github.com/pkg/errors"

    "github.com/stkim1/pcrypto"
)

type RecordGate interface {
    // the database engine
    DataBase() (*sql.DB)

    // the cert database
    Certdb() (certdb.Accessor)

    // Get the session to database
    Session() (*gorm.DB)
}

type ngError struct {
    s string
}

func (n *ngError) Error() string {
    return n.s
}

var (
    // ItemNotFound
    NoItemFound       = &ngError{"[ERR] NotFound: No items are found"}

    gate *dbGate      = nil
    once sync.Once    = sync.Once{}
)

func SharedRecordGate() (RecordGate) {
    return gate
}

// This is where database instances are initiated
func OpenRecordGate(dataDir, recordFile string) (RecordGate, error) {
    var (
        err error = nil
    )
    // TODO : This is a good place to migrate & load encryption plugin

    // Once a gate is open, it should remain alive until closed
    if gate != nil {
        return gate, nil
    }
    // database file path
    dbPath := filepath.Join(dataDir, recordFile)
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
    database, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // cert engine
    cert, err := pcrypto.NewPocketCertStorage(database)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    sess, err := gorm.Open("sqlite3", database)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    if !sess.HasTable(&SlaveNode{}) {
        sess.CreateTable(&SlaveNode{}, &ClusterMeta{})
    } else {
        // Migrate the schema
        sess.AutoMigrate(&SlaveNode{}, &ClusterMeta{});
    }

    return openStorageGate(database, cert, sess), nil
}

func CloseRecordGate() error {
    var err error = gate.close()
    // now, it's time to reset everything. You should not set once to a new instance as only one gate instance should
    // persist throughout the lifecycle of an application
    gate = nil
    //once = sync.Once{}
    return err
}

type dbGate struct {
    database    *sql.DB
    certsess    certdb.Accessor
    ormsess     *gorm.DB
}

func openStorageGate(database *sql.DB, cert certdb.Accessor, sess *gorm.DB) (*dbGate) {
    once.Do(func() {
        gate = &dbGate{
            database:    database,
            certsess:    cert,
            ormsess:     sess,
        }
    })
    return gate
}

// Close closes the currently active connection to the database and clears caches.
func (d *dbGate) close() (error) {
    var (
        err error = nil
        db *sql.DB = nil
    )
    db = d.database

    d.database = nil
    d.ormsess = nil
    d.certsess = nil

    err = db.Close()
    return errors.WithStack(err)
}

func (d *dbGate) DataBase() (*sql.DB) {
    return d.database
}

func (d *dbGate) Certdb() (certdb.Accessor) {
    return d.certsess
}

func (d *dbGate) Session() (*gorm.DB) {
    return d.ormsess
}