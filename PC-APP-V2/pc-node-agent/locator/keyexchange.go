package locator

import (
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/slagent"
)

func newKeyexchangeState(searchComm SearchTx, beaconComm BeaconTx, event LocatorOnTransitionEvent) LocatorState {
    ks := &keyexchange{}

    ks.constState                   = SlaveKeyExchange

    ks.constTransitionFailureLimit  = TransitionFailureLimit
    ks.constTransitionTimout        = UnboundedTimeout * time.Duration(TxActionLimit)
    ks.constTxActionLimit           = TxActionLimit
    ks.constTxTimeWindow            = UnboundedTimeout

    ks.lastTransitionTS             = time.Now()

    ks.timestampTransition          = ks.transitionActionWithTimestamp
    ks.masterMetaTransition         = ks.transitionWithMasterMeta

    ks.LocatorOnTransitionEvent     = event
    ks.searchComm                   = searchComm
    ks.beaconComm                   = beaconComm
    return ks
}

type keyexchange struct{
    locatorState
}

func (ls *keyexchange) transitionActionWithTimestamp(slaveTimestamp time.Time) error {
    slctx := slcontext.SharedSlaveContext()

    masterAgentName, err := slctx.GetMasterAgent()
    if err != nil {
        return err
    }
    agent, err := slagent.KeyExchangeStatus(slaveTimestamp)
    if err != nil {
        return err
    }
    sm, err := slagent.KeyExchangeMeta(masterAgentName, agent, slctx.GetPublicKey())
    if err != nil {
        return err
    }
    pm, err := slagent.PackedSlaveMeta(sm)
    if err != nil {
        return err
    }
    ma, err := slcontext.SharedSlaveContext().GetMasterIP4Address()
    if err != nil {
        return err
    }
    if ls.beaconComm == nil {
        return errors.Errorf("[ERR] Comm Channel is nil")
    }
    return ls.beaconComm.UcastSend(ma, pm)
}

func (ls *keyexchange) transitionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (SlaveLocatingTransition, error) {
    if meta == nil || meta.MetaVersion != msagent.MASTER_META_VERSION {
        // if master is wrong version, It's perhaps from different master. we'll skip and wait for another time
        return SlaveTransitionIdle, errors.Errorf("[ERR] Null or incorrect version of master meta")
    }
    if len(meta.RsaCryptoSignature) == 0 {
        return SlaveTransitionFail, errors.Errorf("[ERR] Null or incorrect RSA signature from Master command")
    }
    if len(meta.EncryptedAESKey) == 0 {
        return SlaveTransitionFail, errors.Errorf("[ERR] Null or incorrect AES key from Master command")
    }
    if len(meta.EncryptedMasterCommand) == 0 {
        return SlaveTransitionFail, errors.Errorf("[ERR] Null or incorrect encrypted master command")
    }
    if len(meta.EncryptedSlaveStatus) == 0 {
        return SlaveTransitionFail, errors.Errorf("[ERR] Null or incorrect slave status from master command")
    }

    msAgent, err := slcontext.SharedSlaveContext().GetMasterAgent()
    if err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }
    if meta.MasterBoundAgent != msAgent {
        return SlaveTransitionFail, errors.Errorf("[ERR] Master bound agent is different than current one %s", msAgent)
    }

    aeskey, err := slcontext.SharedSlaveContext().DecryptByRSA(meta.EncryptedAESKey, meta.RsaCryptoSignature)
    if err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }
    slcontext.SharedSlaveContext().SetAESKey(aeskey)

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
        return SlaveTransitionFail, errors.Errorf("[ERR] Incorrect version of master command")
    }
    // if command is not for exchange key, just ignore
    if msCmd.MasterCommandType != msagent.COMMAND_EXCHANGE_CRPTKEY {
        return SlaveTransitionIdle, nil
    }
    // set slave node name
    idPack, err := slcontext.SharedSlaveContext().DecryptByAES(meta.EncryptedSlaveStatus)
    if err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }
    nodeIdentity, err := slagent.UnpackedPocketSlaveIdentity(idPack)
    if err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }
    if len(nodeIdentity.SlaveNodeName) == 0 {
        return SlaveTransitionFail, errors.Errorf("[ERR] invalid slave node name")
    }
    if len(nodeIdentity.SlaveUUID) == 0 {
        return SlaveTransitionFail, errors.Errorf("[ERR] invalid slave node UUID")
    }
    err = slcontext.SharedSlaveContext().SetSlaveNodeName(nodeIdentity.SlaveNodeName)
    if err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }
    err = slcontext.SharedSlaveContext().SetSlaveNodeUUID(nodeIdentity.SlaveUUID)
    if err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }

    return SlaveTransitionOk, nil
}
