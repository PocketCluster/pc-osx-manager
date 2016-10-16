package model

import (
    "sync"
    "fmt"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

var repository *modelRepo
var once sync.Once

func ModelRepoInstance() (repo ModelRepo) {
    repo = singletonModelRepoInstance()
    return
}

func CloseModelRepo() {
    singletonModelRepoInstance().close()
    repository = nil
}

func singletonModelRepoInstance() (repo *modelRepo) {
    once.Do(func() {
        repository = &modelRepo{}
        initializeModelRepo(repository)
    })
    return repository
}

func initializeModelRepo(mr *modelRepo) {
    // TODO : need a path to save all this filess
    sess, err := gorm.Open("sqlite3", "pc-core.db")
    if err != nil {
        // TODO : capture this error
        return
    }

    if !sess.HasTable(&SlaveNode{}) {
        sess.CreateTable(&SlaveNode{})
    } else {
        // Migrate the schema
        sess.AutoMigrate(&SlaveNode{});
    }
    mr.session = sess
}

type modelRepo struct {
    session         *gorm.DB
}

// Close closes the currently active connection to the database and clears caches.
func (mr *modelRepo) close() (err error) {
    if mr.session == nil {
        err = fmt.Errorf("[ERR] Null session cannot be closed")
        return
    }
    err = mr.session.Close()
    mr.session = nil
    return
}

// Collection returns a collection reference given a table name.
func (mr *modelRepo) Session() (sess *gorm.DB, err error) {
    if mr.session == nil {
        err = fmt.Errorf("[ERR] Null session cannot be queried")
        return
    }
    sess = mr.session
    return
}

func (mr *modelRepo) Where(query interface{}, args ...interface{}) (*gorm.DB) {
    if mr.session == nil {
        return nil
    }
    return mr.session.Where(query, args)
}
