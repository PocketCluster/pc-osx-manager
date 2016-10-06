package msagent

type PocketMasterReponderMeta struct {
    MetaVersion           MetaProtocol                       `msgpack:"pc_ms_pm"`
    StatusCollector       []byte                             `msgpack:"pc_ms_sc", omitempty`
    DiscoveryResponder    *PocketMasterDiscoveryResponder    `msgpack:"pc_ms_dr", inline, omitempty`
}
