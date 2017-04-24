package crash

import "fmt"

type CrashType uint32

const (
    CrashEmergentExit   CrashType = iota
)

func (c CrashType) String() string {
    switch c {
    case CrashEmergentExit:
        return "CrashEmergentExit"
    default:
        return fmt.Sprintf("lifecycle.Stage(%d)", c)
    }
}

type Crash struct {
    Reason    CrashType
}

func (c Crash) String() string {
    return "Crash Reason : " + c.Reason.String()
}