package model

import (
//    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/jinzhu/gorm"
)

const (
    packageTable        string = `pc_package`
    PackageModelVersion string = "0.1.4"
)

type Package struct {
    gorm.Model                `json:"-"`
    // Application specific ID
    AppVer          string    `gorm:"column:app_ver;type:VARCHAR(16)"           json:"app-ver"`
    // Package unique id
    PkgID           string    `gorm:"column:pkg_id;type:VARCHAR(36) UNIQUE" sql:"index" json:"pkg-id"`
    // package name
    Name            string    `gorm:"column:name;type:VARCHAR(255)"             json:"name"`
    // Package Family
    Family          string    `gorm:"column:family;type:VARCHAR(255)"           json:"family"`
    // Description
    Description     string    `gorm:"column:description;type:VARCHAR(255)"      json:"description"`
    // Package revision
    PkgVer          string    `gorm:"column:pkg_ver;type:VARCHAR(32)"           json:"pkg-ver"`

    // Package Meta URL
    MetaURL         string    `gorm:"column:meta_url;type:VARCHAR(255)"         json:"meta-url"`
    // Package Meta Checksum
    MetaChksum      string    `gorm:"column:meta_chksum;type:VARCHAR(32)"       json:"meta-chksum"`

    // Core Node architecture
    CoreArch        string    `gorm:"column:core_arch;type:VARCHAR(32)"         json:"core-arch"`
    // Core Image Checksum
    CoreImageChksum string    `gorm:"column:core_image_chksum;type:VARCHAR(32)" json:"core-image-chksum"`
    // Core Image Sync
    CoreImageSync   string    `gorm:"column:core_image_sync;type:VARCHAR(255)"  json:"core-image-sync"`
    // Core Image URL
    CoreImageURL    string    `gorm:"column:core_image_url;type:VARCHAR(255)"   json:"core-image-url"`

    // Node Architecture
    NodeArch        string    `gorm:"column:node_arch;type:VARCHAR(32)"         json:"node-arch"`
    // Node Image Checksum
    NodeImageChksum string    `gorm:"column:node_image_chksum;type:VARCHAR(32)" json:"node-image-chksum"`
    // Node Image Sync
    NodeImageSync   string    `gorm:"column:node_image_sync;type:VARCHAR(255)"  json:"node-image-sync"`
    // Node Image URL
    NodeImageURL    string    `gorm:"column:node_image_url;type:VARCHAR(255)"   json:"node-image-url"`
}

// instance methods
func (Package) TableName() string {
    return packageTable
}

func FindPackage(query interface{}, args ...interface{}) ([]*Package, error) {
    var pkgs []*Package = nil
    SharedRecordGate().Session().Where(query, args).Find(&pkgs)
    if len(pkgs) == 0 {
        return nil, NoItemFound
    }
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
                // make this gorm.Model w/ PID identical to update otherwise update will fail
                pkgs[i].Model = ppkgs[p].Model
                SharedRecordGate().Session().Save(pkgs[i])
                continue updatelp
            }
        }
        SharedRecordGate().Session().Create(pkgs[i])
    }
    return nil
}
