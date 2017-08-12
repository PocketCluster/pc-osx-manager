package model

import (
    "sync"
    "time"

    "github.com/jinzhu/gorm"
    "github.com/pkg/errors"
)

const (
    CoreNodeName         string = "pc-core"
    coreNodeTable        string = `pc_corenode`
    CoreNodeModelVersion string = "0.1.4"
)

type CoreNode struct {
    gorm.Model

    //Id              uint         `gorm:"column:id;type:INTEGER;primary_key;AUTO_INCREMENT"`
    Joined          time.Time    `gorm:"column:joined;type:DATETIME"`
    Departed        time.Time    `gorm:"column:departed;type:DATETIME"`
    LastAlive       time.Time    `gorm:"column:last_alive;type:DATETIME"`

    ModelVersion    string       `gorm:"column:model_version;type:VARCHAR(8)"`
    // machine uuid from Vbox
    SlaveID         string       `gorm:"column:slave_id;type:VARCHAR(36)"`
    NodeName        string       `gorm:"column:node_name;type:VARCHAR(64)"`
    AuthToken       string       `gorm:"column:auth_token;type:VARCHAR(64)"`

    //slave node state : joined/ departed/ more in the future
    State           string       `gorm:"column:state;type:VARCHAR(32)"`

    //these two are last known elements of core node
    IP4Address      string       `gorm:"column:ip4_address;type:VARCHAR(32)"`
    IP4Gateway      string       `gorm:"column:ip4_gateway;type:VARCHAR(32)"`

    UserMadeName    string       `gorm:"column:user_made_name;type:VARCHAR(256)"`
    PublicKey       []byte       `gorm:"column:public_key;type:BLOB"`
    PrivateKey      []byte       `gorm:"column:private_key;type:BLOB"`

    updateLock      sync.Mutex   `gorm:"-"`
}

func RetrieveCoreNode() *CoreNode {
    var (
        coreNodes []CoreNode = nil
        node *CoreNode = nil
    )
    SharedRecordGate().Session().Find(&coreNodes)
    if len(coreNodes) == 0 {
        // when core is generated, new UUID should be assigned from teleport admin.
        return &CoreNode{
            ModelVersion:    CoreNodeModelVersion,
            // there always is only one core node and it's name is "pc-core"
            NodeName:        CoreNodeName,
            State:           SNMStateInit,
        }
    }

    node = &(coreNodes[0])
    return node
}

// instance methods
func (CoreNode) TableName() string {
    return coreNodeTable
}

// returns IP4 string part only
func (c *CoreNode) IP4AddrString() (string, error) {
    c.updateLock.Lock()
    defer c.updateLock.Unlock()

    return IP4AddrToString(c.IP4Address)
}

// when a slavenode is spawn, an UUID is given, but it could be changed for other joiner (e.g. teleport)
// when the slavenode has not been persisted.
func (c *CoreNode) SetAuthToken(authToken string) error {
    if len(authToken) == 0 {
        return errors.Errorf("[ERR] core auth token should be in appropriate lengh")
    }
    if c.State != SNMStateInit {
        return errors.Errorf("[ERR] cannot modify core auth token when core is not in initial state")
    }
    c.AuthToken = authToken
    return nil
}

func (c *CoreNode) GetAuthToken() (string, error) {
    if len(c.AuthToken) == 0 {
        return "", errors.Errorf("[ERR] invalid core auth token")
    }
    return c.AuthToken, nil
}

func (c *CoreNode) CreateCore() error {
    if c.ModelVersion != CoreNodeModelVersion {
        return errors.Errorf("[ERR] incorrect core model version")
    }
    if c.NodeName != CoreNodeName {
        return errors.Errorf("[ERR] incorrect node name")
    }
    if c.State != SNMStateInit {
        return errors.Errorf("[ERR] cannot join corenode when core isn't init state")
    }
    // TODO : check token format
    if len(c.AuthToken) == 0 {
        return errors.Errorf("[ERR] incorrect uuid length")
    }
    // TODO : check key format
    if len(c.PublicKey) == 0 {
        return errors.Errorf("[ERR] incorrect core public key")
    }
    if len(c.PrivateKey) == 0 {
        return errors.Errorf("[ERR] incorrect core private key")
    }
    c.State = SNMStateJoined
    ts := time.Now()
    c.CreatedAt = ts
    c.Joined = ts
    SharedRecordGate().Session().Create(c)
    return nil
}

func (c *CoreNode) Update() error {
    if c.State != SNMStateJoined {
        return errors.Errorf("[ERR] core node state is not SNMStateJoined : " + c.State)
    }
    if c.ModelVersion != CoreNodeModelVersion {
        return errors.Errorf("[ERR] incorrect core model version")
    }
    if c.NodeName != CoreNodeName {
        return errors.Errorf("[ERR] incorrect node name")
    }
    // TODO : check token format
    if len(c.AuthToken) == 0 {
        return errors.Errorf("[ERR] incorrect uuid length")
    }
    // TODO : check key format
    if len(c.PublicKey) == 0 {
        return errors.Errorf("[ERR] incorrect core public key")
    }
    if len(c.PrivateKey) == 0 {
        return errors.Errorf("[ERR] incorrect core private key")
    }
    c.LastAlive = time.Now()
    SharedRecordGate().Session().Save(c)
    return nil
}

func (c *CoreNode) UpdateIPv4WithGW(ipv4, gwv4 string) error {
    c.updateLock.Lock()
    defer c.updateLock.Unlock()

    // TODO add ip address format check
    if len(ipv4) == 0 {
        return errors.Errorf("[ERR] invalid IPv4 address format")
    }
    if len(gwv4) == 0 {
        return errors.Errorf("[ERR] invalid gateway address format")
    }

    c.IP4Address = ipv4
    c.IP4Gateway = gwv4
    return nil
}