package model

import (
    "github.com/jinzhu/gorm"
    "github.com/pkg/errors"
)

const (
    templateMetaTable string = `pc_template`
)

type TemplateMeta struct {
    gorm.Model
    // Package unique id
    PkgID    string    `gorm:"column:pkg_id;type:VARCHAR(36) UNIQUE" sql:"index"`
    // tempalte body. SQLite3 TEXT type can hold text as large as 2^31-1 characters. Hope that's enough ( https://sqlite.org/limits.html )
    Body     []byte    `gorm:"column:body;type:BLOB"`
}

func (TemplateMeta) TableName() string {
    return templateMetaTable
}

func FindTemplateWithPackageID(pkgID string) (*TemplateMeta, error) {
    var tmpl []*TemplateMeta = nil
    if len(pkgID) != 36 {
        return nil, errors.Errorf("invalid package id to search template")
    }
    SharedRecordGate().Session().Where("pkg_id = ?", pkgID).Find(&tmpl)
    if len(tmpl) == 0 {
        return nil, NoItemFound
    }
    return tmpl[0], nil
}

func NewTemplateMeta() *TemplateMeta {
    return &TemplateMeta{}
}

func (t *TemplateMeta) Update() error {
    if len(t.PkgID) != 36 {
        return errors.Errorf("invalid package id to update")
    }
    if len(t.Body) == 0 {
        return errors.Errorf("invalid template body to update")
    }

    _, err := FindTemplateWithPackageID(t.PkgID)
    if err != nil {
        if err == NoItemFound {
            SharedRecordGate().Session().Create(t)
            return nil
        } else {
            return errors.WithStack(err)
        }
    }

    SharedRecordGate().Session().Save(t)
    return nil
}