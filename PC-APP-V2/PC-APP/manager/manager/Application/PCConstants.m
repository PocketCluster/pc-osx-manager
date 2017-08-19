//
//  PCConstants.m
//  manager
//
//  Created by Almighty Kim on 11/1/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCConstants.h"

#pragma mark - PROTOCOL KEYS

//# ------ VERSION ------
NSString * const PC_PROTO            = @"pc_ver";
NSString * const VERSION             = @"1.0.0";

//------ network interfaces ------
NSString * const ADDRESS             = @"address";
NSString * const NETMASK             = @"netmask";
NSString * const BROADCS             = @"broadcast";
NSString * const GATEWAY             = @"gateway";
NSString * const NAMESRV             = @"dns-nameservers";

//------ protocol definitions ------
NSString * const MASTER_COMMAND_TYPE = @"pc_ma_ct";
NSString * const COMMAND_FIX_BOUND   = @"ct_fix_bound";

//------ MASTER SECTION ------
NSString * const MASTER_SECTION      = @"master";

//bound-id. The master's agent name
NSString * const MASTER_BOUND_AGENT  = @"pc_ma_ba";
//master ip4 / ip6
NSString * const MASTER_IP4_ADDRESS  = @"pc_ma_i4";
NSString * const MASTER_IP6_ADDRESS  = @"pc_ma_i6";
//master hostname
NSString * const MASTER_HOSTNAME     = @"pc_ma_hn";
// master datetime
NSString * const MASTER_DATETIME     = @"pc_ma_dt";
// master timezone
NSString * const MASTER_TIMEZONE     = @"pc_ma_tz";

//------ SLAVE SECTION ------
NSString * const SLAVE_SECTION       = @"slave";

//node looks for agent
NSString * const SLAVE_LOOKUP_AGENT  = @"pc_sl_la";
NSString * const SLAVE_NODE_MACADDR  = @"pc_sl_ma";
NSString * const SLAVE_NODE_NAME     = @"pc_sl_nm";
NSString * const SLAVE_TIMEZONE      = @"pc_sl_tz";
NSString * const SLAVE_CLUSTER_MEMBERS = @"pc_sl_cl";

//------ network configuration ------
NSString * const POCKETCAST_GROUP    = @"239.193.127.127";
const NSUInteger PAGENT_SEND_PORT    = 10060;
const NSUInteger PAGENT_RECV_PORT    = 10061;
const double UNBOUNDED_TIMEOUT       = 3.0;
const double BOUNDED_TIMEOUT         = 10.0;

//------ Webserver path
NSString * const WEB_SERVER_ROOT_PATH = @"/bigpkg/archive";
const int WEB_SERVER_PORT             = 10120;

#pragma mark - NSUserDefaults KEYS
//Raspberry collection
NSString * const kPCPrefDefaultTerm     = @"default_terminal";
NSString * const kPCClusterType         = @"pocketcluster.cluster-type";

NSString * const kRaspberryCollection   = @"raspberries";
NSString * const kRaspberryClusterArray = @"raspberryclusterArray";
NSString * const kRaspberryClusterId    = @"raspberryclusterId";
NSString * const kRaspberryClusterTitle = @"raspberryclusterTitle";
NSString * const kRaspberryClusterSwapOn = @"raspberryclusterSwapOn";

NSString * const kPCInstalledPackageCollection = @"installed-package-collection";

#pragma mark - NOTIFICATION KEY (APPLICATION)
NSString * const kPOCKET_CLUSTER_UPDATE_AVAILABLE                    = @"pocketcluster.update-available";
NSString * const kPOCKET_CLUSTER_UPDATE_VALUE                        = @"pocketcluster.is_update_available";
NSString * const kPOCKET_CLUSTER_NODE_COUNT                          = @"pocketcluster.node-count";
NSString * const kPOCKET_CLUSTER_LIVE_NODE_COUNT                     = @"pocketcluster.live-node-count";

#pragma mark - NOTIFICATION KEY (RASPBERRY)
NSString * const kRASPBERRY_MANAGER_NOTIFICATION_PREFERENCE_CHANGED  = @"raspberry-manager.notification-preference-changed";
NSString * const kRASPBERRY_MANAGER_REFRESHING_STARTED               = @"raspberry-manager.refreshing-started";
NSString * const kRASPBERRY_MANAGER_REFRESHING_ENDED                 = @"raspberry-manager.refreshing-ended";
NSString * const kRASPBERRY_MANAGER_UPDATE_LIVE_NODE_COUNT           = @"raspberry-manager.update-running-node-count";
NSString * const kRASPBERRY_MANAGER_UPDATE_NODE_COUNT                = @"raspberry-manager.update-node-count";
NSString * const kRASPBERRY_MANAGER_NODE                             = @"raspberry-manager.node";
NSString * const kRASPBERRY_MANAGER_NODE_ADDED                       = @"raspberry-manager.node-added";
NSString * const kRASPBERRY_MANAGER_NODE_REMOVED                     = @"raspberry-manager.node-removed";
NSString * const kRASPBERRY_MANAGER_NODE_UPDATED                     = @"raspberry-manager.node-updated";


// Package Process Notification
NSString * const kPOCKET_CLUSTER_PACKAGE_PROCESS_STATUS              = @"pocketcluster.package-process-status";
NSString * const kPOCKET_CLUSTER_PACKAGE_PROCESS_ISALIVE             = @"pocketcluster.package-process-isalive";
NSString * const kPOCKET_CLUSTER_PACKAGE_IDENTIFICATION              = @"pocketcluster.package-identification";

NSString * const kPOCKET_CLUSTER_SALT_STATE_PATH                     = @"/pocket/salt/states";
