package model

import (
    "time"

    "github.com/pkg/errors"
    "github.com/jinzhu/gorm"
)

const (
    coreNodeName         string = "pc-core"
    coreNodeTable        string = `corenode`
    CoreNodeModelVersion string = "0.1.4"
)

type CoreNode struct {
    gorm.Model

    //Id              uint         `gorm:"column:id;type:INTEGER;primary_key;AUTO_INCREMENT"`
    Joined          time.Time    `gorm:"column:joined;type:DATETIME"`
    Departed        time.Time    `gorm:"column:departed;type:DATETIME"`
    LastAlive       time.Time    `gorm:"column:last_alive;type:DATETIME"`

    ModelVersion    string       `gorm:"column:model_version;type:VARCHAR(8)"`
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
            NodeName:        coreNodeName,
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
    if c.NodeName != coreNodeName {
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
    ts := time.Now()
    c.Joined = ts
    SharedRecordGate().Session().Create(c)
    return nil
}

func (c *CoreNode) JoinCore() error {
    if c.ModelVersion != CoreNodeModelVersion {
        return errors.Errorf("[ERR] incorrect core model version")
    }
    if c.NodeName != coreNodeName {
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
    ts := time.Now()
    c.State = SNMStateJoined
    c.Joined = ts
    c.LastAlive = ts
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
    if c.NodeName != coreNodeName {
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
    c.LastAlive = time.Now()
    SharedRecordGate().Session().Save(c)
    return nil
}
