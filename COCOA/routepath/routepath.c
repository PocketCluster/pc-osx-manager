//
//  route_path.c
//  static-core
//
//  Created by Almighty Kim on 8/13/17.
//  Copyright © 2017 PocketCluster. All rights reserved.
//

#include "routepath.h"

const char* RPATH_SYSTEM_READINESS          = "/v1/inquiry/system/readiness";
const char* RPATH_APP_EXPIRED               = "/v1/inquiry/app/expired";
const char* RPATH_USER_AUTHED               = "/v1/inquiry/user/authed";
const char* RPATH_SYSTEM_IS_FIRST_RUN       = "/v1/inquiry/system/is-first-run";
const char* RPATH_PACKAGE_LIST              = "/v1/inquiry/package/list";
const char* RPATH_PACKAGE_INSTALL           = "/v1/cmd/package/install";
const char* RPATH_PACKAGE_INSTALL_PROGRESS  = "/v1/monitor/package/install";
const char* RPATH_PACKAGE_STARTUP           = "/v1/cmd/package/startup";
const char* RPATH_PACKAGE_KILL              = "/v1/cmd/package/kill";

const char* RPATH_MONITOR_PACKAGE_PROCESS   = "/v1/monitor/package/process";
const char* RPATH_MONITOR_SERVICE_STATUS    = "/v1/monitor/service/status";
const char* RPATH_MONITOR_NODE_STATUS       = "/v1/monitor/node/status";

const char* RPATH_NOTI_SRVC_ONLINE_TIMEUP  = "/v1/noti/srvc/online-timeup";
const char* RPATH_NOTI_NODE_ONLINE_TIMEUP  = "/v1/noti/node/online-timeup";