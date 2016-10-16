package model

import "github.com/jinzhu/gorm"

type ModelRepo interface {
    // Get the session to database
    Session() (sess *gorm.DB, err error)

    // Find according to where conditions
    Where(query interface{}, args ...interface{}) (*gorm.DB)
}