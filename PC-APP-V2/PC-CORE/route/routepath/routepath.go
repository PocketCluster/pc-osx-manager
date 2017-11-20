// +build darwin
package routepath

/*

#include "routepath.h"

*/
import "C"

func RpathSystemContextInit() string {
    return C.GoString(C.RPATH_CONTEXT_INIT)
}

func RpathSystemNetworkInit() string {
    return C.GoString(C.RPATH_NETWORK_INIT)
}

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

func RpathPackageListAvailable() string {
    return C.GoString(C.RPATH_PACKAGE_LIST_AVAILABLE)
}

func RpathPackageListInstalled() string {
    return C.GoString(C.RPATH_PACKAGE_LIST_INSTALLED)
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

func RpathNodeRegStart() string {
    return C.GoString(C.RPATH_NODE_REG_START)
}

func RpathNodeUnregList() string {
    return C.GoString(C.RPATH_NODE_UNREG_LIST)
}

func RpathNodeRegCandiate() string {
    return C.GoString(C.RPATH_NODE_REG_CANDIDATE)
}

func RpathNodeRegConfirm() string {
    return C.GoString(C.RPATH_NODE_REG_CONFIRM)
}

func RpathNodeRegStop() string {
    return C.GoString(C.RPATH_NODE_REG_STOP)
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

func RpathAppPrepShutdown() string {
    return C.GoString(C.RPATH_APP_SHUTDOWN_READY)
}
