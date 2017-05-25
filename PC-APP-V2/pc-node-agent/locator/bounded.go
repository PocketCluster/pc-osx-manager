package locator

import (
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-node-agent/slcontext"
)

func newBoundedState(searchComm SearchTx, beaconComm BeaconTx, event LocatorOnTransitionEvent) LocatorState {
    bs := &bounded{}

    bs.constState                   = SlaveBounded

    bs.constTransitionFailureLimit  = TransitionFailureLimit
    bs.constTransitionTimout        = BoundedTimeout * time.Duration(TxActionLimit)
    bs.constTxActionLimit           = TxActionLimit
    bs.constTxTimeWindow            = BoundedTimeout

    bs.lastTransitionTS             = time.Now()

    bs.timestampTransition          = bs.transitionActionWithTimestamp
    bs.masterMetaTransition         = bs.transitionWithMasterMeta
    bs.onTransitionSuccess          = bs.onStateTranstionSuccess
    bs.onTransitionFailure          = bs.onStateTranstionFailure

    bs.LocatorOnTransitionEvent     = event
    bs.searchComm                   = searchComm
    bs.beaconComm                   = beaconComm
    return bs
}

type bounded struct{
    locatorState
}

func (ls *bounded) transitionActionWithTimestamp(slaveTimestamp time.Time) error {
    slctx := slcontext.SharedSlaveContext()

    masterAgentName, err := slctx.GetMasterAgent()
    if err != nil {
        return errors.WithStack(err)
    }
    slaveAgentName, err := slctx.GetSlaveNodeName()
    if err != nil {
        return errors.WithStack(err)
    }
    slaveUUID, err := slctx.GetSlaveNodeUUID()
    if err != nil {
        return errors.WithStack(err)
    }
    aesCryptor, err := slctx.AESCryptor()
    if err != nil {
        return errors.WithStack(err)
    }
    sa, err := slagent.SlaveBoundedStatus(slaveAgentName, slaveUUID, slaveTimestamp)
    if err != nil {
        return errors.WithStack(err)
    }
    sm, err := slagent.SlaveBoundedMeta(masterAgentName, sa, aesCryptor)
    if err != nil {
        return errors.WithStack(err)
    }
    pm, err := slagent.PackedSlaveMeta(sm)
    if err != nil {
        return errors.WithStack(err)
    }
    ma, err := slcontext.SharedSlaveContext().GetMasterIP4Address()
    if err != nil {
        return errors.WithStack(err)
    }
    if ls.beaconComm == nil {
        return errors.Errorf("[ERR] Comm Channel is nil")
    }
    return ls.beaconComm.UcastSend(ma, pm)
}

func (ls *bounded) transitionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (SlaveLocatingTransition, error) {
    if meta == nil || meta.MetaVersion != msagent.MASTER_META_VERSION {
        // if master is wrong version, It's perhaps from different master. we'll skip and wait for another time
        return SlaveTransitionIdle, errors.Errorf("[ERR] Null or incorrect version of master meta")
    }

    msAgent, err := slcontext.SharedSlaveContext().GetMasterAgent()
    if err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }
    // The return value should be SlaveTransitionFail.
    // But, that would lead to bind break with a malicious attack as MasterBoundAgent is exposed.
    // So, we'll relax a bit here to accomodate attach and have room to accept correct input
    if meta.MasterBoundAgent != msAgent {
        return SlaveTransitionIdle, errors.Errorf("[ERR] master bound agent is different than commissioned one %s", msAgent)
    }
    if len(meta.EncryptedMasterCommand) == 0 {
        return SlaveTransitionFail, errors.Errorf("[ERR] null or incorrect encrypted master command")
    }
    // aes decryption of command
    pckedCmd, err := slcontext.SharedSlaveContext().DecryptByAES(meta.EncryptedMasterCommand)
    if err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }
    msCmd, err := msagent.UnpackedMasterCommand(pckedCmd)
    if err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }
    if msCmd.Version != msagent.MASTER_COMMAND_VERSION {
        return SlaveTransitionFail, errors.Errorf("[ERR] incorrect version of master command")
    }
    // if command is not for exchange key, proceed to fail
    if msCmd.MasterCommandType != msagent.COMMAND_SLAVE_ACK {
        return SlaveTransitionFail, errors.Errorf("[ERR] invalid master command type")
    }

    // We'll reset TX action count to 0 and now so successful tx action can happen infinitely
    // We need to reset the counter here when correct master meta comes in
    // It is b/c when succeeded in confirming with master, we should be able to keep receiving master meta
    ls.txActionCount = 0

    // we do not reply here so there will not be an endless master <-> slave loop across network.
    // In fact, one second delayed respose might increase a window of opportunity to get unbounded,
    // but it would not create a chance of overflowing network.

    return SlaveTransitionOk, nil
}

func (ls *bounded) onStateTranstionSuccess(slaveTimestamp time.Time) error {
    return nil
}

func (ls *bounded) onStateTranstionFailure(slaveTimestamp time.Time) error {
    slcontext.SharedSlaveContext().DiscardAESKey()
    return nil
}