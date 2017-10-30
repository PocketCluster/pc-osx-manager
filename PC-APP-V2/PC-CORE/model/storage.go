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

var (
    // ItemNotFound
    NoItemFound       = &ngError{"[ERR] no item found"}

    gate *dbGate      = nil
    gateLock          = sync.Mutex{}
)

type ngError struct {
    s string
}

func (n *ngError) Error() string {
    return n.s
}

type dbGate struct {
    database    *sql.DB
    certsess    certdb.Accessor
    ormsess     *gorm.DB
}

func (d *dbGate) DataBase() *sql.DB {
    return d.database
}

func (d *dbGate) Certdb() certdb.Accessor {
    return d.certsess
}

func (d *dbGate) Session() *gorm.DB {
    return d.ormsess
}

type RecordGate interface {
    // the database engine
    DataBase() (*sql.DB)

    // the cert database
    Certdb() (certdb.Accessor)

    // Get the session to database
    Session() (*gorm.DB)
}

func SharedRecordGate() (RecordGate) {
    gateLock.Lock()
    defer gateLock.Unlock()

    return gate
}

// This is where database instances are initiated
func OpenRecordGate(dataDir, recordFile string) (RecordGate, error) {
    gateLock.Lock()
    defer gateLock.Unlock()

    var (
        cmeta = &ClusterMeta{}
        umeta = &UserMeta{}
        snode = &SlaveNode{}
        cnode = &CoreNode{}
        pkgck = &Package{}
        precd = &PkgRecord{}
        tmplt = &TemplateMeta{}

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

    // TODO : disable when release
    //sess.LogMode(true)

    if !sess.HasTable(cmeta) {
        sess.CreateTable(cmeta)
    } else {
        sess.AutoMigrate(cmeta)
    }

    if !sess.HasTable(umeta) {
        sess.CreateTable(umeta)
    } else {
        sess.AutoMigrate(umeta)
    }

    if !sess.HasTable(snode) {
        sess.CreateTable(snode)
    } else {
        sess.AutoMigrate(snode)
    }

    if !sess.HasTable(cnode) {
        sess.CreateTable(cnode)
    } else {
        sess.AutoMigrate(cnode)
    }

    if !sess.HasTable(pkgck) {
        sess.CreateTable(pkgck)
    } else {
        sess.AutoMigrate(pkgck)
    }

    if !sess.HasTable(precd) {
        sess.CreateTable(precd)
    } else {
        sess.AutoMigrate(precd)
    }

    if !sess.HasTable(tmplt) {
        sess.CreateTable(tmplt)
    } else {
        sess.AutoMigrate(tmplt)
    }

    gate = &dbGate{
        database:    database,
        certsess:    cert,
        ormsess:     sess,
    }
    return gate, nil
}

func CloseRecordGate() error {
    gateLock.Lock()
    defer gateLock.Unlock()

    var db *sql.DB = gate.database

    // now, it's time to reset everything. You should not set once to a new instance as only one gate instance should
    // persist throughout the lifecycle of an application
    gate.database = nil
    gate.ormsess = nil
    gate.certsess = nil
    gate = nil

    return errors.WithStack(db.Close())
}
