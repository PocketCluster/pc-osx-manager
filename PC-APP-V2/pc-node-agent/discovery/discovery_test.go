package discovery

import (
    "testing"
    "time"
    "bytes"

    "github.com/stkim1/pc-node-agent/crypt"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-core/context"
)

var masterBoundAgentName string = ""
const slaveNodeName string = "pc-node1"
var initSendTimestmap, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")

//--

func setUp() {
    masterBoundAgentName, _ = context.DebugContextPrepared().MasterAgentName()
}

func tearDown() {
    context.DebugContextDestroyed()
}

func TestUnboundedState_InquiredTransition(t *testing.T) {
    setUp()
    defer tearDown()

    meta, err := MasterIdentityInqueryRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }

    context := slcontext.DebugSlaveContext(crypt.TestSlavePublicKey(), crypt.TestSlavePrivateKey())
    ssd := NewSlaveDiscovery(context)
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

    meta, err := MasterIdentityFixationRespond(initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    // set to slave discovery state to "Inquired"
    context := slcontext.DebugSlaveContext(crypt.TestSlavePublicKey(), crypt.TestSlavePrivateKey())
    sd := NewSlaveDiscovery(context)
    sd.(*slaveDiscovery).discoveryState = SlaveInquired

    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 2)); err != nil {
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

    meta, err := MasterIdentityFixationRespond(initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    // set to slave discovery state to "Inquired"
    context := slcontext.DebugSlaveContext(crypt.TestSlavePublicKey(), crypt.TestSlavePrivateKey())
    sd := NewSlaveDiscovery(context)
    sd.(*slaveDiscovery).discoveryState = SlaveInquired

    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 2)); err != nil {
        t.Errorf(err.Error())
        return
    }
    if sd.CurrentState() != SlaveKeyExchange {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }

    // get master meta with aeskey
    meta, err = MasterKeyExchangeCommand(masterBoundAgentName, slaveNodeName, initSendTimestmap.Add(time.Second * 3))
    if err != nil {
        t.Error(err.Error())
        return
    }

    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 4)); err != nil {
        t.Error(err.Error())
        return
    }
    if sd.CurrentState() != SlaveCryptoCheck {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }


    // Verification
    if msName, _ := context.GetMasterAgent(); msName != masterBoundAgentName {
        t.Errorf("[ERR] master node name is setup inappropriately | Current : %s\n", msName)
        return
    }
    if snName, _ := context.GetSlaveNodeName(); snName != slaveNodeName {
        t.Errorf("[ERR] slave node name is setup inappropriately | Current : %s\n", snName)
        return
    }
    if bytes.Compare(context.GetAESKey(), crypt.TestAESKey) != 0 {
        t.Errorf("[ERR] slave aes key is setup inappropriately")
        return
    }
}

func TestCryptoCheck_BoundedTransition(t *testing.T) {
    setUp()
    defer tearDown()

    meta, err := MasterIdentityFixationRespond(initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    // set to slave discovery state to "Inquired"
    context := slcontext.DebugSlaveContext(crypt.TestSlavePublicKey(), crypt.TestSlavePrivateKey())
    sd := NewSlaveDiscovery(context)
    sd.(*slaveDiscovery).discoveryState = SlaveInquired

    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 2)); err != nil {
        t.Errorf(err.Error())
        return
    }

    // get master meta with aeskey
    meta, err = MasterKeyExchangeCommand(masterBoundAgentName, slaveNodeName, initSendTimestmap.Add(time.Second * 3))
    if err != nil {
        t.Error(err.Error())
        return
    }

    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 4)); err != nil {
        t.Error(err.Error())
        return
    }
    if sd.CurrentState() != SlaveCryptoCheck {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }

    // get master bind ready
    meta, err = MasterCryptoCheckCommand(masterBoundAgentName, slaveNodeName, initSendTimestmap.Add(time.Second * 5))
    if err != nil {
        t.Error(err.Error())
        return
    }
    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 6)); err != nil {
        t.Error(err.Error())
        return
    }

    if sd.CurrentState() != SlaveBounded {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }

    meta, err = MasterBrokenBindRecoveryCommand(masterBoundAgentName, initSendTimestmap.Add(time.Second * 7))
    if err != nil {
        t.Error(err.Error())
        return
    }
    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 8)); err != nil {
        t.Error(err.Error())
        return
    }
    // now broken bind is recovered
    if sd.CurrentState() != SlaveBounded {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }

    // Verification
    if msName, _ := context.GetMasterAgent(); msName != masterBoundAgentName {
        t.Errorf("[ERR] master node name is setup inappropriately | Current : %s\n", msName)
        return
    }
    if snName, _ := context.GetSlaveNodeName(); snName != slaveNodeName {
        t.Errorf("[ERR] slave node name is setup inappropriately | Current : %s\n", snName)
        return
    }
    if bytes.Compare(context.GetAESKey(), crypt.TestAESKey) != 0 {
        t.Errorf("[ERR] slave aes key is setup inappropriately")
        return
    }
}