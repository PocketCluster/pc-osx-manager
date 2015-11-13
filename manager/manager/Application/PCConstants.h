//
//  PCConstants.h
//  manager
//
//  Created by Almighty Kim on 11/1/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#pragma once

#define MAX_TRIAL_RASP_NODE_COUNT (6)

//# ------ VERSION ------
extern NSString * const PC_PROTO;
extern NSString * const VERSION;

//------ network interfaces ------
extern NSString * const ADDRESS;
extern NSString * const NETMASK;
extern NSString * const BROADCS;
extern NSString * const GATEWAY;
extern NSString * const NAMESRV;

//------ protocol definitions ------
extern NSString * const MASTER_COMMAND_TYPE;
extern NSString * const COMMAND_FIX_BOUND;

//------ MASTER SECTION ------
extern NSString * const MASTER_SECTION;

//bound-id
extern NSString * const MASTER_BOUND_AGENT;
//master ip4 / ip6
extern NSString * const MASTER_IP4_ADDRESS;
extern NSString * const MASTER_IP6_ADDRESS;
//master hostname
extern NSString * const MASTER_HOSTNAME;
// master datetime
extern NSString * const MASTER_DATETIME;
// master timezone
extern NSString * const MASTER_TIMEZONE;
//------ SLAVE SECTION ------
extern NSString * const SLAVE_SECTION;

//node looks for agent
extern NSString * const SLAVE_LOOKUP_AGENT;
extern NSString * const SLAVE_NODE_MACADDR;
extern NSString * const SLAVE_NODE_NAME;
extern NSString * const SLAVE_TIMEZONE;
extern NSString * const SLAVE_CLUSTER_MEMBERS;

//------ network configuration ------
extern NSString * const POCKETCAST_GROUP;
extern const NSUInteger PAGENT_SEND_PORT;
extern const NSUInteger PAGENT_RECV_PORT;
extern const double UNBOUNDED_TIMEOUT;
extern const double BOUNDED_TIMEOUT;

//------ Webserver path
extern NSString * const WEB_SERVER_ROOT_PATH;
extern const int WEB_SERVER_PORT;

//NSUserDefault KEYS
extern NSString * const kPCPrefDefaultTerm;
extern NSString * const kPCClusterType;

extern NSString * const kRaspberryCollection;
extern NSString * const kRaspberryClusterArray;
extern NSString * const kRaspberryClusterId;
extern NSString * const kRaspberryClusterTitle;
extern NSString * const kRaspberryClusterSwapOn;

extern NSString * const kPCVagrantNetInterface;
extern NSString * const kPCInstalledPackageCollection;

// --- MONITORING MANAGEMENT
typedef enum PCClusterType {
    PC_CLUSTER_NONE = 0
    ,PC_CLUTER_VAGRANT
    ,PC_CLUSTER_RASPBERRY
    ,PC_CLUSTER_TYPE_SIZE
} PCClusterType;

//Notfication Key
extern NSString * const kPOCKET_CLUSTER_UPDATE_AVAILABLE;
extern NSString * const kPOCKET_CLUSTER_UPDATE_VALUE;
extern NSString * const kPOCKET_CLUSTER_NODE_COUNT;
extern NSString * const kPOCKET_CLUSTER_LIVE_NODE_COUNT;

extern NSString * const kVAGRANT_MANAGER_NOTIFICATION_PREFERENCE_CHANGED;
extern NSString * const kVAGRANT_MANAGER_INSTANCE;
extern NSString * const kVAGRANT_MANAGER_INSTANCE_OLD;
extern NSString * const kVAGRANT_MANAGER_INSTANCE_NEW;
extern NSString * const kVAGRANT_MANAGER_INSTANCE_ADDED;
extern NSString * const kVAGRANT_MANAGER_INSTANCE_REMOVED;
extern NSString * const kVAGRANT_MANAGER_INSTANCE_UPDATED;
extern NSString * const kVAGRANT_MANAGER_REFRESHING_STARTED;
extern NSString * const kVAGRANT_MANAGER_REFRESHING_ENDED;
extern NSString * const kVAGRANT_MANAGER_UPDATE_RUNNING_VM_COUNT;
extern NSString * const kVAGRANT_MANAGER_UPDATE_INSTANCES_COUNT;

extern NSString * const kRASPBERRY_MANAGER_NOTIFICATION_PREFERENCE_CHANGED;
extern NSString * const kRASPBERRY_MANAGER_REFRESHING_STARTED;
extern NSString * const kRASPBERRY_MANAGER_REFRESHING_ENDED;
extern NSString * const kRASPBERRY_MANAGER_UPDATE_LIVE_NODE_COUNT;
extern NSString * const kRASPBERRY_MANAGER_UPDATE_NODE_COUNT;
extern NSString * const kRASPBERRY_MANAGER_NODE;
extern NSString * const kRASPBERRY_MANAGER_NODE_ADDED;
extern NSString * const kRASPBERRY_MANAGER_NODE_REMOVED;
extern NSString * const kRASPBERRY_MANAGER_NODE_UPDATED;

// Package Process Notification
extern NSString * const kPOCKET_CLUSTER_PACKAGE_PROCESS_STATUS;
extern NSString * const kPOCKET_CLUSTER_PACKAGE_PROCESS_ISALIVE;
extern NSString * const kPOCKET_CLUSTER_PACKAGE_IDENTIFICATION;

//Master Salt State Directory
extern NSString * const kPOCKET_CLUSTER_SALT_STATE_PATH;

typedef enum PCLibraryError {
    PC_LIB_JAVA = 25
    ,PC_LIB_BREW
    ,PC_LIB_VIRTUABOX
    ,PC_LIB_VAGRANT
} PCLibraryError;

#define PROCESS_REFRESH_TIME_INTERVAL (5.0)
