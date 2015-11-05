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

//------ SLAVE SECTION ------
NSString * const SLAVE_SECTION       = @"slave";

//node looks for agent
NSString * const SLAVE_LOOKUP_AGENT  = @"pc_sl_la";
NSString * const SLAVE_NODE_MACADDR  = @"pc_sl_ma";
NSString * const SLAVE_NODE_NAME     = @"pc_sl_nm";

//------ network configuration ------
NSString * const POCKETCAST_GROUP    = @"239.193.127.127";
const NSUInteger PAGENT_SEND_PORT    = 10060;
const NSUInteger PAGENT_RECV_PORT    = 10061;
const double UNBOUNDED_TIMEOUT       = 3.0;
const double BOUNDED_TIMEOUT         = 10.0;

//------ time settings ------
NSString * const MASTER_TIME_STAMP   = @"pc_ma_ts";
NSString * const MASTER_TIME_ZONE    = @"pc_ma_tz";

//------ Webserver path
NSString * const WEB_SERVER_ROOT_PATH = @"/bigpkg/archive";
const int WEB_SERVER_PORT            = 10120;

#pragma mark - NSUserDefaults KEYS
//Raspberry collection
NSString * const kPCPrefDefaultTerm     = @"default_terminal";
NSString * const kPCClusterType         = @"pocketcluster.cluster-type";

NSString * const kRaspberryCollection   = @"raspberries";
NSString * const kRaspberryClusterArray = @"raspberryclusterArray";
NSString * const kRaspberryClusterId    = @"raspberryclusterId";
NSString * const kRaspberryClusterTitle = @"raspberryclusterTitle";

NSString * const kPCVagrantNetInterface = @"vagrant-net-interface";

#pragma mark - NOTIFICATION KEY (APPLICATION)
NSString * const kPOCKET_CLUSTER_UPDATE_AVAILABLE                    = @"pocketcluster.update-available";
NSString * const kPOCKET_CLUSTER_UPDATE_VALUE                        = @"pocketcluster.is_update_available";

#pragma mark - NOTIFICATION KEY (VAGRANT)
NSString * const kVAGRANT_MANAGER_NOTIFICATION_PREFERENCE_CHANGED    = @"vagrant-manager.notification-preference-changed";
NSString * const kVAGRANT_MANAGER_INSTANCE                           = @"vagrant-manager.instance";
NSString * const kVAGRANT_MANAGER_INSTANCE_OLD                       = @"vagrant-manager.instance.old";
NSString * const kVAGRANT_MANAGER_INSTANCE_NEW                       = @"vagrant-manager.instance.new";
NSString * const kVAGRANT_MANAGER_INSTANCE_ADDED                     = @"vagrant-manager.instance-added";
NSString * const kVAGRANT_MANAGER_INSTANCE_REMOVED                   = @"vagrant-manager.instance-removed";
NSString * const kVAGRANT_MANAGER_INSTANCE_UPDATED                   = @"vagrant-manager.instance-updated";
NSString * const kVAGRANT_MANAGER_REFRESHING_STARTED                 = @"vagrant-manager.refreshing-started";
NSString * const kVAGRANT_MANAGER_REFRESHING_ENDED                   = @"vagrant-manager.refreshing-ended";
NSString * const kVAGRANT_MANAGER_UPDATE_RUNNING_VM_COUNT            = @"vagrant-manager.update-running-vm-count";
NSString * const kVAGRANT_MANAGER_UPDATE_INSTANCES_COUNT             = @"vagrant-manager.update-instances-count";

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

NSString * const kPOCKET_CLUSTER_SALT_STATE_PATH                     = @"/pocket/salt/states";
