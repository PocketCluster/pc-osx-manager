//
//  route_path.c
//  static-core
//
//  Created by Almighty Kim on 8/13/17.
//  Copyright © 2017 PocketCluster. All rights reserved.
//

#include "routepath.h"

const char* RPATH_SYSTEM_READINESS       = "/v1/inquiry/system/readiness";
const char* RPATH_APP_EXPIRED            = "/v1/inquiry/app/expired";
const char* RPATH_USER_AUTHED            = "/v1/inquiry/user/authed";
const char* RPATH_SYSTEM_IS_FIRST_RUN    = "/v1/inquiry/system/is-first-run";
const char* RPATH_PACKAGE_LIST           = "/v1/inquiry/package/list";

const char* RPATH_MONITOR_NODE_REGISTERED   = "/v1/monitor/node/registered";
const char* RPATH_MONITOR_NODE_UNREGISTERED = "/v1/monitor/node/unregistered";
const char* RPATH_MONITOR_SERVICE_STATUS    = "/v1/monitor/service/status";