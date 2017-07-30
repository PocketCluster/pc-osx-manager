package context

import (
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/model"
)

type HostContextClusterMeta interface {
    SetClusterMeta(meta *model.ClusterMeta) error
    MasterAgentName() (string, error)
    GetClusterUUID() (string, error)
    GetClusterDomain() (string, error)
}

func (c *hostContext) SetClusterMeta(meta *model.ClusterMeta) error {
    c.Lock()
    defer c.Unlock()

    if meta == nil {
        return errors.Errorf("[ERR] cluster metadata is null")
    }
    c.ClusterMeta = meta
    return nil
}

//TODO rename this function to ClusterID
func (c *hostContext) MasterAgentName() (string, error) {
    if len(c.ClusterID) == 0 {
        return "", errors.Errorf("[ERR] invalid cluster id")
    }
    return c.ClusterID, nil
}

func (c *hostContext) GetClusterUUID() (string, error) {
    if len(c.ClusterUUID) == 0 {
        return "", errors.Errorf("[ERR] invalid cluster uuid")
    }
    return c.ClusterUUID, nil
}

func (c *hostContext) GetClusterDomain() (string, error) {
    if len(c.ClusterDomain) == 0 {
        return "", errors.Errorf("[ERR] invalid cluster domain")
    }
    return c.ClusterDomain, nil
}
