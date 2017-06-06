package model

import (
    "strconv"
    "time"

    "github.com/pkg/errors"
    "github.com/jinzhu/gorm"
    "strings"
    "github.com/pborman/uuid"
)

const slaveNodeTable string = `slavenode`

const SlaveNodeModelVersion = "0.1.4"

const (
    SNMFieldId              = "id"
    SNMFieldJoined          = "joined"
    SNMFieldDeparted        = "departed"
    SNMFieldLastAlive       = "last_alive"
    SNMFieldSlaveID         = "slave_id"
    SNMFieldHardware        = "hardware"
    SNMFieldNodeName        = "node_name"
    SNMFieldState           = "state"
    SNMFieldAuthToken       = "auth_token"
    SNMFieldIP4Address      = "ip4_address"
    SNMFieldIP4Gateway      = "ip4_gateway"
    SNMFieldIP4Netmask      = "ip4_netmask"
    SNMFieldUserMadeName    = "user_made_name"
    SNMFieldPublicKey       = "public_key"
    SNMFieldPrivateKey      = "private_key"
)

const (
    SNMStateInit            = "node_init"
    SNMStateJoined          = "node_joined"
    SNMStateDeparted        = "node_departed"
)

// the purpose of node sanitization is 1) to give name, and 2) to make it aligned with other node data
type NodeSanitizerFunc func(s *SlaveNode) error

func (n NodeSanitizerFunc) Sanitize(s *SlaveNode) error {
    return n(s)
}

type NodeSanitizer interface {
    Sanitize(s *SlaveNode) error
}

type SlaveNode struct {
    gorm.Model

    //Id              uint         `gorm:"column:id;type:INTEGER;primary_key;AUTO_INCREMENT"`
    Joined          time.Time    `gorm:"column:joined;type:DATETIME"`
    Departed        time.Time    `gorm:"column:departed;type:DATETIME"`
    LastAlive       time.Time    `gorm:"column:last_alive;type:DATETIME"`

    ModelVersion    string       `gorm:"column:model_version;type:VARCHAR(8)"`
    SlaveID         string       `gorm:"column:slave_id;type:VARCHAR(32)"`
    Hardware        string       `gorm:"column:hardware;type:VARCHAR(32)"`
    NodeName        string       `gorm:"column:node_name;type:VARCHAR(64)"`
    AuthToken       string       `gorm:"column:auth_token;type:VARCHAR(64)"`

    // slave node       s tate : joined/ departed/ more in the future
    State           string       `gorm:"column:state;type:VARCHAR(32)"`

    // these two are last known addresses
    IP4Address      string       `gorm:"column:ip4_address;type:VARCHAR(32)"`
    IP4Gateway      string       `gorm:"column:ip4_gateway;type:VARCHAR(32)"`

    UserMadeName    string       `gorm:"column:user_made_name;type:VARCHAR(256)"`
    PublicKey       []byte       `gorm:"column:public_key;type:BLOB"`
    PrivateKey      []byte       `gorm:"column:private_key;type:BLOB"`

    sanitizer       NodeSanitizer   `gorm:"-"`
}

func IP4AddrToString(ip4Addr string) (string, error) {
    if len(ip4Addr) == 0 {
        return "", errors.Errorf("[ERR] empty address")
    }
    addrform := strings.Split(ip4Addr, "/")
    if !strings.Contains(ip4Addr, "/") || len(addrform) != 2 {
        return "", errors.Errorf("[ERR] invalid ip4 + subnet format")
    }
    return addrform[0], nil
}

func NewSlaveNode(ns NodeSanitizer) *SlaveNode {
    return &SlaveNode {
        ModelVersion:    SlaveNodeModelVersion,
        // whenever slave is generated, new UUID should be assigned to it.
        AuthToken:       uuid.New(),
        State:           SNMStateInit,
        sanitizer:       ns,
    }
}

// TODO : this should be deprecated
func InsertSlaveNode(slave *SlaveNode) (error) {
    if slave == nil {
        return errors.Errorf("[ERR] Slave node is null")
    }
    if slave.State != SNMStateJoined {
        return errors.Errorf("[ERR] Slave node state is not SNMStateJoined : " + slave.State)
    }
    if slave.ModelVersion != SlaveNodeModelVersion {
        return errors.Errorf("[ERR] incorrect slave model version")
    }
    SharedRecordGate().Session().Create(slave)
    return nil
}

func FindAllSlaveNode() ([]SlaveNode, error) {
    var nodes []SlaveNode = nil
    SharedRecordGate().Session().Find(&nodes)
    if len(nodes) == 0 {
        return nil, NoItemFound
    }
    return nodes, nil
}

func FindSlaveNode(query interface{}, args ...interface{}) ([]SlaveNode, error) {
    var nodes []SlaveNode = nil
    SharedRecordGate().Session().Where(query, args).Find(&nodes)
    return nodes, nil
}

func FindSlaveNameCandiate() (string, error) {
    return "", errors.Errorf("[ERR] this function is deprecated")
    var nodes []SlaveNode = nil
    SharedRecordGate().Session().Where(string(SNMFieldState + " = ?"), SNMStateJoined).Find(&nodes)
    return "pc-node" + strconv.Itoa(len(nodes) + 1), nil
}

func UpdateSlaveNode(slave *SlaveNode) (error) {
    SharedRecordGate().Session().Save(slave)
    return nil
}

func DeleteAllSlaveNode() (error) {
    SharedRecordGate().Session().Delete(SlaveNode{})
    return nil
}

// instance methods
func (SlaveNode) TableName() string {
    return slaveNodeTable
}

// returns IP4 string part only
func (s *SlaveNode) IP4AddrString() (string, error) {
    return IP4AddrToString(s.IP4Address)
}

func (s *SlaveNode) SanitizeSlave() error {
    if s.sanitizer == nil {
        return errors.Errorf("[ERR] node sanitizer is nil")
    }
    return s.sanitizer.Sanitize(s)
}

// when a slavenode is spawn, an UUID is given, but it could be changed for other joiner (e.g. teleport)
// when the slavenode has not been persisted.
func (s *SlaveNode) SetSlaveID(id string) error {
    if len(id) == 0 {
        return errors.Errorf("[ERR] slave id should be in appropriate lengh")
    }
    if s.State != SNMStateInit {
        return errors.Errorf("[ERR] cannot modify slave id when slave is not in SNMStateInit")
    }
    s.AuthToken = id
    return nil
}

func (s *SlaveNode) JoinSlave() error {
    if s.ModelVersion != SlaveNodeModelVersion {
        return errors.Errorf("[ERR] incorrect slave model version")
    }
    // TODO : check token format
    if len(s.AuthToken) == 0 {
        return errors.Errorf("[ERR] incorrect uuid length")
    }
    // TODO : check node name formet
    if len(s.NodeName) == 0 {
        return errors.Errorf("[ERR] incorrect node name")
    }
    // TODO : check key format
    if len(s.PublicKey) == 0 {
        return errors.Errorf("[ERR] incorrect slave Public key")
    }
    ts := time.Now()
    s.State = SNMStateJoined
    s.Joined = ts
    s.LastAlive = ts
    SharedRecordGate().Session().Create(s)
    return nil
}

func (s *SlaveNode) Update() error {
    if s.State != SNMStateJoined {
        return errors.Errorf("[ERR] Slave node state is not SNMStateJoined : " + s.State)
    }
    if s.ModelVersion != SlaveNodeModelVersion {
        return errors.Errorf("[ERR] incorrect slave model version")
    }
    // TODO : check token format
    if len(s.AuthToken) == 0 {
        return errors.Errorf("[ERR] incorrect uuid length")
    }
    // TODO : check node name formet
    if len(s.NodeName) == 0 {
        return errors.Errorf("[ERR] incorrect node name")
    }
    // TODO : check key format
    if len(s.PublicKey) == 0 {
        return errors.Errorf("[ERR] incorrect slave Public key")
    }
    s.LastAlive = time.Now()
    SharedRecordGate().Session().Save(s)
    return nil
}

func (s *SlaveNode) RemoveSanitizer() {
    s.sanitizer = nil
}