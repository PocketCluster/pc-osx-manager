//
//  PCConstants.h
//  manager
//
//  Created by Almighty Kim on 11/1/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#pragma once

#define MAX_TRIAL_RASP_NODE_COUNT (4)

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

//NSUserDefault KEYS
extern NSString * const kPCPrefDefaultTerm;

//Notfication Key
extern NSString * const kPOCKET_CLUSTER_UPDATE_AVAILABLE;
extern NSString * const kPOCKET_CLUSTER_UPDATE_VALUE;
extern NSString * const kPOCKET_CLUSTER_NODE_COUNT;
extern NSString * const kPOCKET_CLUSTER_LIVE_NODE_COUNT;

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
