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
const char* RPATH_CMD_SERVICE_START      = "/v1/cmd/service-start";

const char* RPATH_MONITOR_NODE_BOUNDED   = "/v1/monitor/node/bounded";
const char* RPATH_MONITOR_NODE_UNBOUNDED = "/v1/monitor/node/unbounded";