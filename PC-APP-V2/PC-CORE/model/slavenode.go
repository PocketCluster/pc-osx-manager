package model

import (
    "fmt"
    "strconv"
    "time"

    "github.com/jinzhu/gorm"
)

const slaveNodeCreation string = `CREATE TABLE IF NOT EXISTS slavenode (id INTEGER PRIMARY KEY ASC, joined DATETIME, departed DATETIME, last_alive DATETIME, mac_address VARCHAR(32), arch VARCHAR(32), node_name VARCHAR(64), state VARCHAR(32), ip4_address VARCHAR(32), ip4_gateway VARCHAR(32), ip4_netmask VARCHAR(32), public_key TEXT, private_key TEXT);`

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

    // slave node       s tate : joined/ departed/ more in the future
    State           string       `gorm:"column:state;type:VARCHAR(32)"`

    IP4Address      string       `gorm:"column:ip4_address;type:VARCHAR(32)"`
    IP4Gateway      string       `gorm:"column:ip4_gateway;type:VARCHAR(32)"`
    IP4Netmask      string       `gorm:"column:ip4_netmask;type:VARCHAR(32)"`

    UserMadeName    string       `gorm:"column:user_made_name;type:VARCHAR(256)"`
    PublicKey       []byte       `gorm:"column:public_key;type:BLOB"`
    PrivateKey      []byte       `gorm:"column:private_key;type:BLOB"`
}

func (SlaveNode) TableName() string {
    return slaveNodeTable
}

func NewSlaveNode() *SlaveNode {
    return &SlaveNode{
        ModelVersion: SlaveNodeModelVersion,
        State       : SNMStateJoined,
    }
}

func InsertSlaveNode(slave *SlaveNode) (err error) {
    if slave == nil {
        return fmt.Errorf("[ERR] Slave node is null")
    }
    if slave.State != SNMStateJoined {
        return fmt.Errorf("[ERR] Slave node state is not SNMStateJoined : " + slave.State)
    }
    if slave.ModelVersion != SlaveNodeModelVersion {
        return fmt.Errorf("[ERR] ")
    }
    db, err := SharedModelRepoInstance().Session()
    if err != nil {
        return
    }
    db.Create(slave)
    return
}

func FindAllSlaveNode() (nodes []SlaveNode, err error) {
    db, err := SharedModelRepoInstance().Session()
    if err != nil {
        return
    }
    db.Find(&nodes)
    return
}

func FindSlaveNode(query interface{}, args ...interface{}) (nodes []SlaveNode, err error) {
    db, err := SharedModelRepoInstance().Session()
    if err != nil {
        return
    }
    db.Where(query, args).Find(&nodes)
    return
}

func FindSlaveNameCandiate() (string, error) {
    db, err := SharedModelRepoInstance().Session()
    if err != nil {
        return "", err
    }
    var nodes []SlaveNode
    db.Where(string(SNMFieldState + " = ?"), SNMStateJoined).Find(&nodes)
    return "pc-node" + strconv.Itoa(len(nodes) + 1), nil
}

func UpdateSlaveNode(slave *SlaveNode) (err error) {
    db, err := SharedModelRepoInstance().Session()
    if err != nil {
        return
    }
    db.Save(slave)
    return
}

func DeleteAllSlaveNode() (err error) {
    db, err := SharedModelRepoInstance().Session()
    if err != nil {
        return
    }
    db.Delete(SlaveNode{})
    return
}
