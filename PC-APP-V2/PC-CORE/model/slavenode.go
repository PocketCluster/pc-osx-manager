package model

import (
    "strconv"
    "time"

    "github.com/pkg/errors"
    "github.com/jinzhu/gorm"
)

const slaveNodeTable string = `slavenode`

const SlaveNodeModelVersion = "0.1.4"

const (
    SNMFieldId              = "id"
    SNMFieldJoined          = "joined"
    SNMFieldDeparted        = "departed"
    SNMFieldLastAlive       = "last_alive"
    SNMFieldMacAddress      = "mac_address"
    SNMFieldArch            = "arch"
    SNMFieldNodeName        = "node_name"
    SNMFieldState           = "state"
    SNMFieldIP4Address      = "ip4_address"
    SNMFieldIP4Gateway      = "ip4_gateway"
    SNMFieldIP4Netmask      = "ip4_netmask"
    SNMFieldUserMadeName    = "user_made_name"
    SNMFieldPublicKey       = "public_key"
    SNMFieldPrivateKey      = "private_key"
)

const (
    SNMStateJoined          = "node_joined"
    SNMStateDeparted        = "node_departed"
)

type SlaveNode struct {
    gorm.Model

    //Id              uint         `gorm:"column:id;type:INTEGER;primary_key;AUTO_INCREMENT"`
    Joined          time.Time    `gorm:"column:joined;type:DATETIME"`
    Departed        time.Time    `gorm:"column:departed;type:DATETIME"`
    LastAlive       time.Time    `gorm:"column:last_alive;type:DATETIME"`

    ModelVersion    string       `gorm:"column:model_version;type:VARCHAR(8)"`
    MacAddress      string       `gorm:"column:mac_address;type:VARCHAR(32)"`
    Arch            string       `gorm:"column:arch;type:VARCHAR(32)"`
    NodeName        string       `gorm:"column:node_name;type:VARCHAR(64)"`
    SlaveID         string       `gorm:"column:slave_id;type:VARCHAR(64)"`

    // slave node       s tate : joined/ departed/ more in the future
    State           string       `gorm:"column:state;type:VARCHAR(32)"`

    // these two are last known addresses
    IP4Address      string       `gorm:"column:ip4_address;type:VARCHAR(32)"`
    IP4Gateway      string       `gorm:"column:ip4_gateway;type:VARCHAR(32)"`

    UserMadeName    string       `gorm:"column:user_made_name;type:VARCHAR(256)"`
    PublicKey       []byte       `gorm:"column:public_key;type:BLOB"`
    PrivateKey      []byte       `gorm:"column:private_key;type:BLOB"`
}

func (SlaveNode) TableName() string {
    return slaveNodeTable
}

func NewSlaveNode() *SlaveNode {
    return &SlaveNode {
        ModelVersion:    SlaveNodeModelVersion,
        State:           SNMStateJoined,
    }
}

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
    var (
        nodes []SlaveNode = nil
        err error = nil
    )
    SharedRecordGate().Session().Find(&nodes)
    if len(nodes) == 0 {
        return nil, NoItemFound
    }
    return nodes, err
}

func FindSlaveNode(query interface{}, args ...interface{}) ([]SlaveNode, error) {
    var (
        nodes []SlaveNode = nil
        err error = nil
    )
    SharedRecordGate().Session().Where(query, args).Find(&nodes)
    return nodes, err
}

func FindSlaveNameCandiate() (string, error) {
    var (
        nodes []SlaveNode = nil
    )
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
