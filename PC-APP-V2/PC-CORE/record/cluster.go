package record

import (
    "github.com/jinzhu/gorm"
    "github.com/pborman/uuid"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/utils"
)

type ClusterMeta struct {
    gorm.Model

    // this is short id
    ClusterID      string

    // this is for teleport and other things
    ClusterUUID    string
}

func NewClusterMeta() (*ClusterMeta) {
    return &ClusterMeta{
        ClusterID:    utils.NewRandomString(32),
        ClusterUUID:  uuid.New(),
    }
}

func FindClusterMeta() ([]*ClusterMeta, error) {
    var (
        meta []*ClusterMeta = nil
        err error = nil
    )
    SharedRecordGate().Session().Find(&meta)
    if len(meta) == 0 {
        return nil, NoItemFound
    }
    return meta, err
}

func UpsertClusterMeta(meta *ClusterMeta) (error) {
    if meta == nil {
        return errors.Errorf("[ERR] invalid null cluster meta")
    }
    if len(meta.ClusterID) == 0 {
        return errors.Errorf("[ERR] invalid cluster group id")
    }
    if len(meta.ClusterUUID) == 0 {
        return errors.Errorf("[ERR] invalid cluster UUID")
    }
    SharedRecordGate().Session().Create(meta)
    return nil
}
