package model

import (
    "fmt"

    "github.com/jinzhu/gorm"
    "github.com/pborman/uuid"
    "github.com/pkg/errors"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-core/utils/randstr"
)

const (
    clusterMetaTable        string = `pc_clustermeta`
)

type ClusterMeta struct {
    gorm.Model
    // this is short id
    ClusterID        string    `gorm:"column:cluster_id;type:VARCHAR(16)"`
    // this is for teleport and other things
    ClusterUUID      string    `gorm:"column:cluster_uuid;type:VARCHAR(36)"`
    // Cluster Domain name
    ClusterDomain    string    `gorm:"column:cluster_domain;type:VARCHAR(42)"`
    // User defined Name
    UserMadeName     string    `gorm:"column:user_made_name;type:VARCHAR(255)"`
}

// instance methods
func (ClusterMeta) TableName() string {
    return clusterMetaTable
}

func NewClusterMeta() (*ClusterMeta) {
    var (
        cid string = randstr.NewRandomString(16)
        domain string = fmt.Sprintf(pcrypto.FormFQDNClusterID, cid)
    )
    return &ClusterMeta{
        // (03/25/2017)
        // cluster id length is 16 for now. It should suffice to count all the cluster in the world.
        // Later, if it's necessary, we'll increase the length to cover
        ClusterID:        cid,
        ClusterUUID:      uuid.New(),
        ClusterDomain:    domain,
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
    if len(meta.ClusterDomain) == 0 {
        return errors.Errorf("[ERR] invalid cluster domain name")
    }
    SharedRecordGate().Session().Create(meta)
    return nil
}

func (c *ClusterMeta) Update() (error) {
    if len(c.ClusterID) == 0 {
        return errors.Errorf("[ERR] invalid cluster group id")
    }
    if len(c.ClusterUUID) == 0 {
        return errors.Errorf("[ERR] invalid cluster UUID")
    }
    if len(c.ClusterDomain) == 0 {
        return errors.Errorf("[ERR] invalid cluster domain name")
    }
    SharedRecordGate().Session().Save(c)
    return nil
}
