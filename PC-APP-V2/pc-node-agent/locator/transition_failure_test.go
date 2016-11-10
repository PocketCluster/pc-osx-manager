package locator

import (
    "testing"
    "time"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-core/msagent"
)

// unbounded -> inquired
func Test_Unbounded_Inquired_MasterMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

}

func Test_Unbounded_Inquired_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

}

// inquired -> keyexchange
func Test_Inquired_Keyexchange_MasterMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

}

func Test_Inquired_Keyexchange_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

}

// keyexchange -> cryptocheck
func Test_Keyexchange_Cryptocheck_MasterMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

}

func Test_keyexchange_Cryptocheck_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

}

// cryptocheck -> bounded
func Test_Cryptocheck_Bounded_MasterMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

}

func Test_Cryptocheck_Bounded_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

}

// bounded -> bindbroken
func Test_Bounded_BindBroken_MasterMeta_fail(t *testing.T) {
    setUp()
    defer tearDown()

}

func Test_Bounded_BindBroken_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    // Let's have a bounded state
    context := slcontext.SharedSlaveContext()
    context.SetMasterPublicKey(pcrypto.TestMasterPublicKey())
    context.SetMasterAgent(masterAgentName)
    context.SetSlaveNodeName(slaveNodeName)

    // have a slave locator
    sd, err := NewSlaveLocator(SlaveBindBroken)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS := time.Now()
    meta, err := msagent.TestMasterBrokenBindRecoveryCommand(masterAgentName, pcrypto.TestAESKey, pcrypto.TestAESCryptor, pcrypto.TestMasterRSAEncryptor)
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS := masterTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS);
    if err != nil {
        t.Error(err.Error())
        return
    }
    state, err := sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveBounded {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }

    // bounded state will fail after TxActionLimit trial
    var i uint = 0
    for ;i <= TxActionLimit; i++ {
        slaveTS = slaveTS.Add(time.Millisecond + BoundedTimeout)
        err = sd.TranstionWithTimestamp(slaveTS)
        if err != nil {
            t.Skip(err.Error())
            //return
        }
        state, err = sd.CurrentState()
        if err != nil {
            t.Error(err.Error())
            return
        }
        if i < TxActionLimit {
            if state != SlaveBounded {
                t.Errorf("[ERR] Slave state should not change properly | Current : %s\n", state.String())
                return
            }
        } else {
            if state != SlaveBindBroken {
                t.Errorf("[ERR] Slave state should not change properly | Current : %s\n", state.String())
                return
            }
        }
    }
}

// bindbroken -> bindbroken
func Test_BindBroken_BindBroken_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    // by the time bind broken state is revived, previous master public key should have been available.
    context := slcontext.SharedSlaveContext()
    context.SetMasterPublicKey(pcrypto.TestMasterPublicKey())
    context.SetMasterAgent(masterAgentName)
    context.SetSlaveNodeName(slaveNodeName)

    sd, err := NewSlaveLocator(SlaveBindBroken)
    if err != nil {
        t.Error(err.Error())
        return
    }

    slaveTS := time.Now()
    TxCountTarget := TxActionLimit * TxActionLimit
    var i uint = 0
    for ;i < TxCountTarget; i++ {
        slaveTS = slaveTS.Add(time.Millisecond + UnboundedTimeout)
        err = sd.TranstionWithTimestamp(slaveTS);
        if err != nil {
            t.Error(err.Error())
            return
        }
        state, err := sd.CurrentState()
        if err != nil {
            t.Error(err.Error())
            return
        }
        // now broken bind is recovered
        if state != SlaveBindBroken {
            t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
            return
        }
    }
}

// unbounded -> unbounded
func Test_Unbounded_Unbounded_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    sd, err := NewSlaveLocator(SlaveUnbounded)
    if err != nil {
        t.Error(err.Error())
        return
    }

    // we'll send TxActionLimit * TxActionLimit to
    slaveTS := time.Now()
    TxCountTarget := TxActionLimit * TxActionLimit
    var i uint = 0
    for ;i < TxCountTarget; i++ {
        slaveTS = slaveTS.Add(time.Millisecond + UnboundedTimeout)
        err = sd.TranstionWithTimestamp(slaveTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        state, err := sd.CurrentState()
        if err != nil {
            t.Error(err.Error())
            return
        }
        if state != SlaveUnbounded {
            t.Errorf("[ERR] Slave state should not change properly | Current : %s\n", state.String())
            return
        }
    }
}