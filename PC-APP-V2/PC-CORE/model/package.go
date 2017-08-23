package model

import (
//    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/jinzhu/gorm"
)

const packageTable string = `pc_package`

const PackageModelVersion = "0.1.4"

type Package struct {
    gorm.Model
    // Application specific ID
    AppVer          string    `gorm:"column:app_ver;type:VARCHAR(16)"           json:"app-ver"`
    // Package unique id
    PkgID           string    `gorm:"column:pkg_id;type:VARCHAR(36) UNIQUE" sql:"index" json:"pkg-id"`
    // Package revision
    PkgVer          int       `gorm:"column:pkg_ver;type:INT"                   json:"pkg-ver"`
    // package name
    Name            string    `gorm:"column:name;type:VARCHAR(255)"             json:"name"`
    // User defined Name
    Family          string    `gorm:"column:family;type:VARCHAR(255)"           json:"family"`
    // User defined Name
    Description     string    `gorm:"column:description;type:VARCHAR(255)"      json:"description"`
    // User defined Name
    MetaURL         string    `gorm:"column:meta_url;type:VARCHAR(255)"         json:"meta-url"`
    // User defined Name
    CoreArch        string    `gorm:"column:core_arch;type:VARCHAR(32)"         json:"core-arch"`
    // User defined Name
    CoreImageURL    string    `gorm:"column:core_image_url;type:VARCHAR(255)"   json:"core-image-url"`
    // User defined Name
    NodeArch        string    `gorm:"column:node_arch;type:VARCHAR(32)"         json:"node-arch"`
    // User defined Name
    NodeImageURL    string    `gorm:"column:node_image_url;type:VARCHAR(255)"   json:"node-image-url"`
}

// instance methods
func (Package) TableName() string {
    return packageTable
}

func FindPackage(query interface{}, args ...interface{}) ([]Package, error) {
    var pkgs []Package = nil
    SharedRecordGate().Session().Where(query, args).Find(&pkgs)
    return pkgs, nil
}

func UpdatePackages(pkgs []*Package) (error) {
    if pkgs == nil || len(pkgs) == 0 {
        return errors.Errorf("[ERR] no packages to update")
    }
    var ppkgs []*Package = nil
    SharedRecordGate().Session().Find(&ppkgs)

    updatelp: for i, _ := range pkgs {
        for p, _ := range ppkgs {
            if pkgs[i].PkgID == ppkgs[p].PkgID {
                SharedRecordGate().Session().Save(pkgs[i])
                continue updatelp
            }
        }
        SharedRecordGate().Session().Create(pkgs[i])
    }
    return nil
}
