// +build darwin
package routepath

/*

#include "routepath.h"

*/
import "C"

func RpathSystemReadiness() string {
    return C.GoString(C.RPATH_SYSTEM_READINESS)
}

func RpathAppExpired() string {
    return C.GoString(C.RPATH_APP_EXPIRED)
}

func RpathUserAuthed() string {
    return C.GoString(C.RPATH_USER_AUTHED)
}

func RpathSystemIsFirstRun() string {
    return C.GoString(C.RPATH_SYSTEM_IS_FIRST_RUN)
}

func RpathCmdServiceStart() string {
    return C.GoString(C.RPATH_CMD_SERVICE_START)
}

func RpathMonitorNodeBounded() string {
    return C.GoString(C.RPATH_MONITOR_NODE_BOUNDED)
}

func RpathMonitorNodeUnbounded() string {
    return C.GoString(C.RPATH_MONITOR_NODE_UNBOUNDED)
}
