package locator

import (
    "time"
    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/locator/state"
)

type SlaveLocator interface {
    CurrentState() state.SlaveLocatingState
    TranstionWithTimestamp(timestamp time.Time) error
    TranstionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) error
    Close() error
}