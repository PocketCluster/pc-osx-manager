package msagent

type PocketMasterAgentMeta struct {
    MetaVersion           MetaProtocol                       `msgpack:"pc_ms_pm"`
    DiscoveryResponder    *PocketMasterDiscoveryResponder    `msgpack:"pc_ms_dr", inline, omitempty`
    StatusCollector       *PocketMasterStatusCommander       `msgpack:"pc_ms_sc", inline, omitempty`
    EncryptedCollector    []byte                             `msgpack:"pc_ms_ec", omitempty`
    MasterPubkey          []byte                             `msgpack:"pc_ms_pk", omitempty`
    EncryptedAESKey       []byte                             `msgpack:"pc_ms_ak", omitempty`
}

func UnboundedInqueryMeta() (meta *PocketMasterAgentMeta, err error) {
    return
}

func InqueredIdentityMeta() (meta *PocketMasterAgentMeta, err error) {
    return
}

func KeyExchangeSendMeta() (meta *PocketMasterAgentMeta, err error) {
    return
}

func CryptoCheckSendMeta() (meta *PocketMasterAgentMeta, err error) {
    return
}

func BoundedStatusMeta() (meta *PocketMasterAgentMeta, err error) {
    return
}

func BindBrokenStatuMeta() (meta *PocketMasterAgentMeta, err error) {
    return
}

