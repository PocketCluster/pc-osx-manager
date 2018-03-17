package model

import (
    "github.com/pkg/errors"
    "github.com/jinzhu/gorm"
)

const (
    packageTable        string = `pc_package`
    PackageModelVersion string = "0.1.4"
)

// this is a model that reflects what's available in api backend
type Package struct {
    gorm.Model                `json:"-"`
    // Application specific ID
    AppVer          string    `gorm:"column:app_ver;type:VARCHAR(16)"           json:"app-ver"`
    // Package unique id
    PkgID           string    `gorm:"column:pkg_id;type:VARCHAR(36) UNIQUE" sql:"index" json:"pkg-id"`
    // Package revision
    PkgVer          string    `gorm:"column:pkg_ver;type:VARCHAR(32)"           json:"pkg-ver"`
    // Package checksum
    PkgChksum       string    `gorm:"column:pkg_chksum;type:VARCHAR(32)"        json:"pkg-chksum"`

    // package name
    Name            string    `gorm:"column:name;type:VARCHAR(255)"             json:"name"`
    // Package Family
    Family          string    `gorm:"column:family;type:VARCHAR(255)"           json:"family"`
    // Menu name
    MenuName        string    `gorm:"column:menu_name;type:VARCHAR(255)"        json:"menu-name"`
    // Description
    Description     string    `gorm:"column:description;type:VARCHAR(255)"      json:"description"`
    // web ports
    WebPorts        string    `gorm:"column:web_ports;type:VARCHAR(255)"        json:"web-ports"`

    // Package Meta URL
    MetaURL         string    `gorm:"column:meta_url;type:VARCHAR(255)"         json:"meta-url"`
    // Package Meta Checksum
    MetaChksum      string    `gorm:"column:meta_chksum;type:VARCHAR(32)"       json:"meta-chksum"`

    // Core Node architecture
    CoreArch        string    `gorm:"column:core_arch;type:VARCHAR(32)"         json:"core-arch"`
    // Core Image Name
    CoreImageName   string    `gorm:"column:core_image_name;type:VARCHAR(255)"  json:"core-image-name"`
    // Core Image Size
    CoreImageSize   string    `gorm:"column:core_image_size;type:VARCHAR(255)"  json:"core-image-size"`
    // Core Image Checksum
    CoreImageChksum string    `gorm:"column:core_image_chksum;type:VARCHAR(32)" json:"core-image-chksum"`
    // Core Image Sync
    CoreImageSync   string    `gorm:"column:core_image_sync;type:VARCHAR(255)"  json:"core-image-sync"`
    // Core Image URL
    CoreImageURL    string    `gorm:"column:core_image_url;type:VARCHAR(255)"   json:"core-image-url"`
    // Core Data path to setup
    CoreDataPath    string    `gorm:"column:core_data_path;type:VARCHAR(255)"   json:"core-data-path"`

    // Node Architecture
    NodeArch        string    `gorm:"column:node_arch;type:VARCHAR(32)"         json:"node-arch"`
    // Node Image Name
    NodeImageName   string    `gorm:"column:node_image_name;type:VARCHAR(255)"  json:"node-image-name"`
    // Node Image Size
    NodeImageSize   string    `gorm:"column:node_image_size;type:VARCHAR(255)"  json:"node-image-size"`
    // Node Image Checksum
    NodeImageChksum string    `gorm:"column:node_image_chksum;type:VARCHAR(32)" json:"node-image-chksum"`
    // Node Image Sync
    NodeImageSync   string    `gorm:"column:node_image_sync;type:VARCHAR(255)"  json:"node-image-sync"`
    // Node Image URL
    NodeImageURL    string    `gorm:"column:node_image_url;type:VARCHAR(255)"   json:"node-image-url"`
    // Node Data path to setup
    NodeDataPath    string    `gorm:"column:node_data_path;type:VARCHAR(255)"   json:"node-data-path"`
}

// class methods
func (Package) TableName() string {
    return packageTable
}

func AllPackages() ([]*Package, error) {
    var pkgs []*Package = nil
    SharedRecordGate().Session().Find(&pkgs)
    if len(pkgs) == 0 {
        return nil, NoItemFound
    }
    return pkgs, nil
}

func FindPackage(query interface{}, args ...interface{}) ([]*Package, error) {
    var pkgs []*Package = nil
    SharedRecordGate().Session().Where(query, args).Find(&pkgs)
    if len(pkgs) == 0 {
        return nil, NoItemFound
    }
    return pkgs, nil
}

func UpsertPackages(pkgs []*Package) (error) {
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
