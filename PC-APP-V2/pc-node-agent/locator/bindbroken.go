package locator

import (
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/slagent"
)

func newBindbrokenState(searchComm SearchTx, beaconComm BeaconTx, event LocatorOnTransitionEvent) LocatorState {
    bs := &bindbroken{}

    bs.constState                   = SlaveBindBroken

    bs.constTransitionFailureLimit  = TransitionFailureLimit
    bs.constTransitionTimout        = UnboundedTimeout * time.Duration(TxActionLimit)
    bs.constTxActionLimit           = TxActionLimit
    bs.constTxTimeWindow            = UnboundedTimeout

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

type bindbroken struct{
    locatorState
}

func (ls *bindbroken) transitionActionWithTimestamp(slaveTimestamp time.Time) error {
    // we'll reset TX action count to 0 and now so successful tx action can happen infinitely until we confirm with master
    // we need to reset the counter here than receiver
    ls.txActionCount = 0

    slctx := slcontext.SharedSlaveContext()
    masterAgentName, err := slctx.GetClusterID()
    if err != nil {
        return errors.WithStack(err)
    }
    sm, err := slagent.BrokenBindMeta(masterAgentName)
    if err != nil {
        return errors.WithStack(err)
    }
    pm, err := slagent.PackedSlaveMeta(sm)
    if err != nil {
        return errors.WithStack(err)
    }
    if ls.searchComm == nil {
        return errors.Errorf("[ERR] Comm Channel is nil")
    }
    return ls.searchComm.McastSend(pm)
}

func (ls *bindbroken) transitionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (SlaveLocatingTransition, error) {
    if meta == nil || meta.MetaVersion != msagent.MASTER_META_VERSION {
        // if master is wrong version, It's perhaps from different master. we'll skip and wait for another time
        return SlaveTransitionIdle, errors.Errorf("[ERR] Null or incorrect version of master meta")
    }
    msAgent, err := slcontext.SharedSlaveContext().GetClusterID()
    if err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }
    if meta.MasterBoundAgent != msAgent {
        return SlaveTransitionFail, errors.Errorf("[ERR] Master bound agent is different than commissioned one %s", msAgent)
    }
    if len(meta.EncryptedMasterRespond) == 0 {
        return SlaveTransitionFail, errors.Errorf("[ERR] Null or incorrect encrypted master respond")
    }
    if len(meta.EncryptedAESKey) == 0 {
        return SlaveTransitionFail, errors.Errorf("[ERR] Null or incorrect AES key from Master command")
    }
    if len(meta.RsaCryptoSignature) == 0 {
        return SlaveTransitionFail, errors.Errorf("[ERR] Null or incorrect RSA signature from Master command")
    }

    aeskey, err := slcontext.SharedSlaveContext().DecryptByRSA(meta.EncryptedAESKey, meta.RsaCryptoSignature)
    if err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }
    slcontext.SharedSlaveContext().SetAESKey(aeskey)

    // aes decryption of command
    pckedRsp, err := slcontext.SharedSlaveContext().DecryptByAES(meta.EncryptedMasterRespond)
    if err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }
    msRsp, err := msagent.UnpackedMasterRespond(pckedRsp)
    if err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }

    if msRsp.Version != msagent.MASTER_RESPOND_VERSION {
        return SlaveTransitionFail, errors.Errorf("[ERR] Null or incorrect version of master meta")
    }
    // if command is not for exchange key, just ignore
    if msRsp.MasterCommandType != msagent.COMMAND_RECOVER_BIND {
        return SlaveTransitionIdle, nil
    }

    // set the master ip address
    if len(msRsp.MasterAddress) == 0 {
        return SlaveTransitionFail, errors.Errorf("[ERR] Null or incorrect master address")
    }
    slcontext.SharedSlaveContext().SetMasterIP4Address(msRsp.MasterAddress)

    return SlaveTransitionOk, nil
}

func (ls *bindbroken) onStateTranstionSuccess(slaveTimestamp time.Time) error {
    return slcontext.SharedSlaveContext().SyncAll()
}

func (ls *bindbroken) onStateTranstionFailure(slaveTimestamp time.Time) error {
    slcontext.SharedSlaveContext().DiscardMasterSession()
    return nil
}
