//
//  PCConstants.h
//  manager
//
//  Created by Almighty Kim on 11/1/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#pragma once

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

//------ SLAVE SECTION ------
extern NSString * const SLAVE_SECTION;

//node looks for agent
extern NSString * const SLAVE_LOOKUP_AGENT;
extern NSString * const SLAVE_NODE_MACADDR;
extern NSString * const SLAVE_NODE_NAME;

//------ network configuration ------
extern NSString * const POCKETCAST_GROUP;
extern const NSUInteger PAGENT_SEND_PORT;
extern const NSUInteger PAGENT_RECV_PORT;

//------ time settings ------
extern NSString * const MASTER_TIME_STAMP;
extern NSString * const MASTER_TIME_ZONE;

//Raspberry collection
extern NSString * const kRaspberryCollection;
extern NSString * const kPCPrefDefaultTerm;

//Notfication Key
extern NSString * const kPOCKET_CLUSTER_UPDATE_AVAILABLE;
extern NSString * const kPOCKET_CLUSTER_UPDATE_VALUE;

extern NSString * const kVAGRANT_MANAGER_INSTANCE;
extern NSString * const kVAGRANT_MANAGER_INSTANCE_OLD;
extern NSString * const kVAGRANT_MANAGER_INSTANCE_NEW;
extern NSString * const kVAGRANT_MANAGER_NOTIFICATION_PREFERENCE_CHANGED;
extern NSString * const kVAGRANT_MANAGER_INSTANCE_ADDED;
extern NSString * const kVAGRANT_MANAGER_INSTANCE_REMOVED;
extern NSString * const kVAGRANT_MANAGER_INSTANCE_UPDATED;
extern NSString * const kVAGRANT_MANAGER_REFRESHING_STARTED;
extern NSString * const kVAGRANT_MANAGER_REFRESHING_ENDED;
extern NSString * const kVAGRANT_MANAGER_UPDATE_RUNNING_VM_COUNT;
extern NSString * const kVAGRANT_MANAGER_UPDATE_INSTANCES_COUNT;

extern NSString * const kRASPBERRY_MANAGER_NODE;
extern NSString * const kRASPBERRY_MANAGER_NOTIFICATION_PREFERENCE_CHANGED;
extern NSString * const kRASPBERRY_MANAGER_REFRESHING_STARTED;
extern NSString * const kRASPBERRY_MANAGER_REFRESHING_ENDED;
extern NSString * const kRASPBERRY_MANAGER_UPDATE_RUNNING_NODE_COUNT;
extern NSString * const kRASPBERRY_MANAGER_UPDATE_NODE_COUNT;
extern NSString * const kRASPBERRY_MANAGER_NODE_UP;
extern NSString * const kRASPBERRY_MANAGER_NODE_DOWN;
extern NSString * const kRASPBERRY_MANAGER_NODE_ADDED;
extern NSString * const kRASPBERRY_MANAGER_NODE_REMOVED;
extern NSString * const kRASPBERRY_MANAGER_NODE_UPDATED;