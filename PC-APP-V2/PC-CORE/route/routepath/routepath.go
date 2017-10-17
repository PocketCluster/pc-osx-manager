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

func RpathPackageList() string {
    return C.GoString(C.RPATH_PACKAGE_LIST)
}

func RpathPackageInstall() string {
    return C.GoString(C.RPATH_PACKAGE_INSTALL)
}

func RpathPackageInstallProgress() string {
    return C.GoString(C.RPATH_PACKAGE_INSTALL_PROGRESS)
}

func RpathPackageStartup() string {
    return C.GoString(C.RPATH_PACKAGE_STARTUP)
}

func RpathPackageKill() string {
    return C.GoString(C.RPATH_PACKAGE_KILL)
}

func RpathPackageProcess() string {
    return C.GoString(C.RPATH_MONITOR_PACKAGE_PROCESS)
}

func RpathMonitorServiceStatus() string {
    return C.GoString(C.RPATH_MONITOR_SERVICE_STATUS)
}

func RpathMonitorNodeStatus() string {
    return C.GoString(C.RPATH_MONITOR_NODE_STATUS)
}

func RpathNotiSrvcOnlineTimeup() string {
    return C.GoString(C.RPATH_NOTI_SRVC_ONLINE_TIMEUP)
}

func RpathNotiNodeOnlineTimeup() string {
    return C.GoString(C.RPATH_NOTI_NODE_ONLINE_TIMEUP)
}
