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
        // (03/25/2017)
        // cluster id length is 16 for now. It should suffice to count all the cluster in the world.
        // Later, if it's necessary, we'll increase the length to cover
        ClusterID:    utils.NewRandomString(16),
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
