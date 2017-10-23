package model

import (
    "github.com/pkg/errors"
    "github.com/jinzhu/gorm"
)

const (
    recordTable string = "pc_pkgrecord"
)

// PkgRecord saves history of what package has been installed
type PkgRecord struct {
    gorm.Model `json:"-"`
    // Application specific ID
    AppVer          string    `gorm:"column:app_ver;type:VARCHAR(16)"           json:"app-ver"`
    // Package unique id
    PkgID           string    `gorm:"column:pkg_id;type:VARCHAR(36) UNIQUE" sql:"index" json:"pkg-id"`
    // Package revision
    PkgVer          string    `gorm:"column:pkg_ver;type:VARCHAR(32)"           json:"pkg-ver"`
    // Package checksum
    PkgChksum       string    `gorm:"column:pkg_chksum;type:VARCHAR(32)"        json:"pkg-chksum"`
}

// class methods
func (PkgRecord) TableName() string {
    return recordTable
}

func AllRecords() ([]*PkgRecord, error) {
    var recs []*PkgRecord = nil
    SharedRecordGate().Session().Find(&recs)
    if len(recs) == 0 {
        return nil, NoItemFound
    }
    return recs, nil
}

func FindRecord(query interface{}, args ...interface{}) ([]*PkgRecord, error) {
    var recs []*PkgRecord = nil
    SharedRecordGate().Session().Where(query, args).Find(&recs)
    if len(recs) == 0 {
        return nil, NoItemFound
    }
    return recs, nil
}

func UpsertRecords(nRecs []*PkgRecord) error {
    if nRecs == nil || len(nRecs) == 0 {
        return errors.Errorf("[ERR] no record to update")
    }
    var oRecs []*PkgRecord = nil
    SharedRecordGate().Session().Find(&oRecs)

    updatelp:
    for n, _ := range nRecs {
        for o, _ := range oRecs {
            if nRecs[n].PkgID == oRecs[o].PkgID {
                nRecs[n].Model = oRecs[o].Model
                SharedRecordGate().Session().Save(nRecs[n])
                continue updatelp
            }
        }
        SharedRecordGate().Session().Create(nRecs[n])
    }
    return nil
}
