//
//  PCConstants.h
//  manager
//
//  Created by Almighty Kim on 11/1/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#pragma once

#define MAX_TRIAL_NODE_COUNT (6)

// --- MONITORING MANAGEMENT
typedef enum PCClusterType {
    PC_CLUSTER_NONE = 0
    ,PC_CLUTER_VAGRANT
    ,PC_CLUSTER_RASPBERRY
    ,PC_CLUSTER_TYPE_SIZE
} PCClusterType;
