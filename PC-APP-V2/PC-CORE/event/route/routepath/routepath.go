// +build darwin
package routepath

/*

#include "routepath.h"

*/
import "C"

var (
    RpathSystemReadiness            string = C.GoString(C.RPATH_SYSTEM_READINESS)
    RpathAppExpired                 string = C.GoString(C.RPATH_APP_EXPIRED)
    RpathUserAuthed                 string = C.GoString(C.RPATH_USER_AUTHED)
    RpathSystemIsFirstRun           string = C.GoString(C.RPATH_SYSTEM_IS_FIRST_RUN)
    RpathCmdServiceStart            string = C.GoString(C.RPATH_CMD_SERVICE_START)

    RpathMonitorNodeBounded         string = C.GoString(C.RPATH_MONITOR_NODE_BOUNDED)
    RpathMonitorNodeUnbounded       string = C.GoString(C.RPATH_MONITOR_NODE_UNBOUNDED)
)