package model

import (
    "time"

    "github.com/jinzhu/gorm"
)

const slaveNodeCreation string = `CREATE TABLE IF NOT EXISTS slavenode (id INTEGER PRIMARY KEY ASC, joined DATETIME, departed DATETIME, last_alive DATETIME, mac_address VARCHAR(32), arch VARCHAR(32), node_name VARCHAR(64), state VARCHAR(32), ip4_address VARCHAR(32), ip4_gateway VARCHAR(32), ip4_netmask VARCHAR(32), public_key TEXT, private_key TEXT);`

const slaveNodeTable string = `slavenode`

type SlaveNode struct {
    gorm.Model

    //Id              uint         `gorm:"column:id;type:INTEGER;primary_key;AUTO_INCREMENT"`
    Joined          time.Time    `gorm:"column:joined;type:DATETIME"`
    Departed        time.Time    `gorm:"column:departed;type:DATETIME"`
    LastAlive       time.Time    `gorm:"column:last_alive;type:DATETIME"`

    MacAddress      string       `gorm:"column:mac_address;type:VARCHAR(32)"`
    Arch            string       `gorm:"column:arch;type:VARCHAR(32)"`
    NodeName        string       `gorm:"column:node_name;type:VARCHAR(64)"`

    // slave node s tate : joined/ departed/ more in the future
    State           string       `gorm:"column:state;type:VARCHAR(32)"`

    IP4Address      string       `gorm:"column:ip4_address;type:VARCHAR(32)"`
    IP4Gateway      string       `gorm:"column:ip4_gateway;type:VARCHAR(32)"`
    IP4Netmask      string       `gorm:"column:ip4_netmask;type:VARCHAR(32)"`

    PublicKey       []byte       `gorm:"column:public_key;type:BLOB"`
    PrivateKey      []byte       `gorm:"column:private_key;type:BLOB"`
}

const (
    Id              = "id"
    Joined          = "joined"
    Departed        = "departed"
    LastAlive       = "last_alive"
    MacAddress      = "mac_address"
    Arch            = "arch"
    NodeName        = "node_name"
    State           = "state"
    IP4Address      = "ip4_address"
    IP4Gateway      = "ip4_gateway"
    IP4Netmask      = "ip4_netmask"
    PublicKey       = "public_key"
    PrivateKey      = "private_key"
)

func (SlaveNode) TableName() string {
    return slaveNodeTable
}

func InsertSlaveNode(slave *SlaveNode) (err error) {
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
