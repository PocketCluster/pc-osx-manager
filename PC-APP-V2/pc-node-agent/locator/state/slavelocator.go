package state

import (
    "time"
    "github.com/stkim1/pc-core/msagent"
)

type slaveLocator struct {
    state       LocatorState
}

func (sl *slaveLocator) CurrentState() SlaveLocatingState {

}

func (sl *slaveLocator) TranstionWithTimestamp(timestamp time.Time) error {

}

func (sl *slaveLocator) TranstionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) error {

}
