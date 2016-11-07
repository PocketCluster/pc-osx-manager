package locator

import (
    "testing"
    "time"
    "bytes"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/msagent"
)

var masterAgentName string
var slaveNodeName string
var initSendTimestmap time.Time

func setUp() {
    masterAgentName, _ = context.DebugContextPrepare().MasterAgentName()
    slaveNodeName = "pc-node1"
    initSendTimestmap, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    slcontext.DebugSlcontextPrepare()
}

func tearDown() {
    slcontext.DebugSlcontextDestroy()
    context.DebugContextDestroy()
}

func TestUnboundedState_InquiredTransition(t *testing.T) {
    setUp()
    defer tearDown()

    meta, err := msagent.TestMasterInquireSlaveRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }

    ssd := NewSlaveDiscovery()
    err = ssd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 2))
    if err != nil {
        t.Error(err.Error())
        return
    }

    if ssd.CurrentState() != SlaveInquired {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", ssd.CurrentState().String())
        return
    }
}

func TestInquired_KeyExchangeTransition(t *testing.T) {
    setUp()
    defer tearDown()

    meta, endTime, err := msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterPublicKey(), initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    // set to slave discovery state to "Inquired"
    sd := NewSlaveDiscovery()
    sd.(*slaveLocator).locatingState = SlaveInquired

    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, endTime.Add(time.Second)); err != nil {
        t.Errorf(err.Error())
        return
    }
    if sd.CurrentState() != SlaveKeyExchange {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }
}

func TestKeyExchange_CryptoCheckTransition(t *testing.T) {
    setUp()
    defer tearDown()

    meta, masterTS, err := msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterPublicKey(), initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    // set to slave discovery state to "Inquired"
    context := slcontext.SharedSlaveContext()
    sd := NewSlaveDiscovery()
    sd.(*slaveLocator).locatingState = SlaveInquired

    // execute state transition
    slaveTS := masterTS.Add(time.Second)
    if err = sd.TranstionWithMasterMeta(meta, slaveTS); err != nil {
        t.Errorf(err.Error())
        return
    }
    if sd.CurrentState() != SlaveKeyExchange {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }

    // get master meta with aeskey
    masterTS = slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterKeyExchangeCommand(masterAgentName, slaveNodeName, pcrypto.TestSlavePublicKey(), pcrypto.TestAESKey, pcrypto.TestAESCryptor, pcrypto.TestMasterRSAEncryptor, masterTS)
    if err != nil {
        t.Error(err.Error())
        return
    }

    // execute state transition
    slaveTS = masterTS.Add(time.Second)
    if err = sd.TranstionWithMasterMeta(meta, slaveTS); err != nil {
        t.Error(err.Error())
        return
    }
    if sd.CurrentState() != SlaveCryptoCheck {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }


    // Verification
    if msName, _ := context.GetMasterAgent(); msName != masterAgentName {
        t.Errorf("[ERR] master node name is setup inappropriately | Current : %s\n", msName)
        return
    }
    if snName, _ := context.GetSlaveNodeName(); snName != slaveNodeName {
        t.Errorf("[ERR] slave node name is setup inappropriately | Current : %s\n", snName)
        return
    }
    if bytes.Compare(context.GetAESKey(), pcrypto.TestAESKey) != 0 {
        t.Errorf("[ERR] slave aes key is setup inappropriately")
        return
    }
}

func Test_Unbounded_Bounded_Onepass(t *testing.T) {
    setUp()
    defer tearDown()

    meta, masterTS, err := msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterPublicKey(), initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    // set to slave discovery state to "Inquired"
    context := slcontext.SharedSlaveContext()
    sd := NewSlaveDiscovery()
    sd.(*slaveLocator).locatingState = SlaveInquired

    // execute state transition
    slaveTS := masterTS.Add(time.Second)
    if err = sd.TranstionWithMasterMeta(meta, slaveTS); err != nil {
        t.Errorf(err.Error())
        return
    }

    // get master meta with aeskey
    masterTS = slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterKeyExchangeCommand(masterAgentName, slaveNodeName, pcrypto.TestSlavePublicKey(), pcrypto.TestAESKey, pcrypto.TestAESCryptor, pcrypto.TestMasterRSAEncryptor, masterTS)

    if err != nil {
        t.Error(err.Error())
        return
    }

    // execute state transition
    slaveTS = masterTS.Add(time.Second)
    if err = sd.TranstionWithMasterMeta(meta, slaveTS); err != nil {
        t.Error(err.Error())
        return
    }
    if sd.CurrentState() != SlaveCryptoCheck {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }

    // get master bind ready
    masterTS = slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterCheckCryptoCommand(masterAgentName, slaveNodeName, pcrypto.TestAESCryptor, masterTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // execute state transition
    slaveTS = masterTS.Add(time.Second)
    if err = sd.TranstionWithMasterMeta(meta, slaveTS); err != nil {
        t.Error(err.Error())
        return
    }

    if sd.CurrentState() != SlaveBounded {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }

    meta, err = msagent.TestMasterBrokenBindRecoveryCommand(masterAgentName, pcrypto.TestAESKey, pcrypto.TestAESCryptor, pcrypto.TestMasterRSAEncryptor)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // execute state transition
    slaveTS = slaveTS.Add(time.Second)
    if err = sd.TranstionWithMasterMeta(meta, slaveTS); err != nil {
        t.Error(err.Error())
        return
    }
    // now broken bind is recovered
    if sd.CurrentState() != SlaveBounded {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }

    // Verification
    if msName, _ := context.GetMasterAgent(); msName != masterAgentName {
        t.Errorf("[ERR] master node name is setup inappropriately | Current : %s\n", msName)
        return
    }
    if snName, _ := context.GetSlaveNodeName(); snName != slaveNodeName {
        t.Errorf("[ERR] slave node name is setup inappropriately | Current : %s\n", snName)
        return
    }
    if bytes.Compare(context.GetAESKey(), pcrypto.TestAESKey) != 0 {
        t.Errorf("[ERR] slave aes key is setup inappropriately")
        return
    }
}