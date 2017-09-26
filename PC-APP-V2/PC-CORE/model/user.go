package model

import (
    "github.com/pkg/errors"
    "github.com/jinzhu/gorm"
    "github.com/pborman/uuid"

    "github.com/stkim1/pc-core/utils/randstr"
)

const (
    userMetaTable string = `pc_usermeta`
)

type UserMeta struct {
    gorm.Model
    // this is short user id
    UserID        string    `gorm:"column:user_id;type:VARCHAR(36)"`
    // short user login name
    Login         string    `gorm:"column:login;type:VARCHAR(36)"`
    // this is for teleport and other things
    Password      string    `gorm:"column:password;type:VARCHAR(8)"`
}

// instance methods
func (UserMeta) TableName() string {
    return userMetaTable
}

func NewUserMeta(login string) (*UserMeta) {
    return &UserMeta{
        UserID:    uuid.New(),
        Login:     login,
        Password:  randstr.NewRandomString(8),
    }
}

func UpsertUserMeta(meta *UserMeta) (error) {
    if meta == nil {
        return errors.Errorf("[ERR] invalid null user meta")
    }
    if len(meta.UserID) == 0 {
        return errors.Errorf("[ERR] invalid user uuid")
    }
    if len(meta.Login) == 0 {
        return errors.Errorf("[ERR] invalid user login name")
    }
    if len(meta.Password) != 8 {
        return errors.Errorf("[ERR] invalid password length")
    }
    SharedRecordGate().Session().Create(meta)
    return nil
}

func FindUserMetaWithLogin(login string) ([]*UserMeta, error) {
    var (
        meta []*UserMeta = nil
        err error = nil
    )
    SharedRecordGate().Session().Where("login = ?", login).Find(&meta)
    if len(meta) == 0 {
        return nil, NoItemFound
    }
    return meta, err
}
