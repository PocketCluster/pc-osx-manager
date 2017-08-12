// +build darwin
package defaults

/*
#include "routepath.h"
*/
import "C"

var (
    RPATH_SYSTEM_READINESS          string = C.GoString(C.ROUTE_SYSTEM_READINESS)
    RPATH_APP_EXPIRED               string = C.GoString(C.ROUTE_APP_EXPIRED)
    RPATH_USER_AUTHED               string = C.GoString(C.ROUTE_USER_AUTHED)
    RPATH_SYSTEM_IS_FIRST_RUN       string = C.GoString(C.ROUTE_SYSTEM_IS_FIRST_RUN)
    RPATH_CMD_SERVICE_START         string = C.GoString(C.ROUTE_CMD_SERVICE_START)

    RPATH_MONITOR_NODE_BOUNDED      string = C.GoString(C.RPATH_MONITOR_NODE_BOUNDED)
    RPATH_MONITOR_NODE_UNBOUNDED    string = C.GoString(C.RPATH_MONITOR_NODE_UNBOUNDED)
)