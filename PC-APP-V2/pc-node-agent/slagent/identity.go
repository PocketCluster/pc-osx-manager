package slagent

import (
    "github.com/pkg/errors"
    "gopkg.in/vmihailenco/msgpack.v2"
)

type PocketSlaveIdentity struct {
    Version          StatusProtocol    `msgpack:"s_ps"`
    // slave nodename
    SlaveNodeName    string            `msgpack:"s_nm,omitempty"`
    // slave UUID
    SlaveUUID        string            `msgpack:"s_uu,omitempty"`
}

func NewPocketSlaveIdentity(nodename, uuid string) *PocketSlaveIdentity {
    return &PocketSlaveIdentity {
        Version:          SLAVE_STATUS_VERSION,
        SlaveNodeName:    nodename,
        SlaveUUID:        uuid,
    }
}

func PackPocketSlaveIdentity(si *PocketSlaveIdentity) ([]byte, error) {
    var (
        pi []byte
        err error = nil
    )
    pi, err = msgpack.Marshal(si)
    return pi, errors.WithStack(err)
}

func UnpackedPocketSlaveIdentity(data []byte) (*PocketSlaveIdentity, error) {
    var (
        si *PocketSlaveIdentity
        err error = nil
    )
    err = errors.WithStack(msgpack.Unmarshal(data, &si))
    return si, err
}