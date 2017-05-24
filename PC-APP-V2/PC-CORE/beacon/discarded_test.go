package beacon

import (
    "testing"
    "time"

    "github.com/stkim1/pc-core/model"
)

func Test_Discard_Shutdown(t *testing.T) {
    setUp()
    defer tearDown()

    // --- VARIABLE PREP ---
    var (
        debugComm CommChannel = &DebugCommChannel{}
        debugEvent BeaconOnTransitionEvent = &DebugTransitionEventReceiver{}
        masterTS  = time.Now()
        mb, err   = NewMasterBeacon(MasterInit, model.NewSlaveNode(slaveSanitizer), debugComm, debugEvent)
    )
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterInit {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- TX ACTION FAIL ---
    for i := 0; i <= int(TxActionLimit); i++ {
        masterTS = masterTS.Add(time.Millisecond + UnboundedTimeout)
        err = mb.TransitionWithTimestamp(masterTS)
        if err != nil {
            t.Log(err.Error())
        }
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    // check if mb shutdown works ok
    mb.Shutdown()
    if mb.SlaveNode() != nil {
        t.Error("[ERR] slavenode should be null after discard")
    }
}