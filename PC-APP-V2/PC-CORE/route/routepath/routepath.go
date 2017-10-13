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

func RpathMonitorNodeRegistered() string {
    return C.GoString(C.RPATH_MONITOR_NODE_REGISTERED)
}

func RpathMonitorNodeUnregistered() string {
    return C.GoString(C.RPATH_MONITOR_NODE_UNREGISTERED)
}

func RpathMonitorNodeBounded() string {
    return C.GoString(C.RPATH_MONITOR_NODE_BOUNDED)
}

func RpathMonitorNodePcssh() string {
    return C.GoString(C.RPATH_MONITOR_NODE_PCSSH)
}

func RpathMonitorNodeOrchst() string {
    return C.GoString(C.RPATH_MONITOR_NODE_ORCHST)
}

func RpathMonitorServiceStatus() string {
    return C.GoString(C.RPATH_MONITOR_SERVICE_STATUS)
}