//
//  PCConstants.m
//  manager
//
//  Created by Almighty Kim on 11/1/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCConstants.h"

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


#pragma mark - NSUserDefaults KEYS
//Raspberry collection
NSString * const kRaspberryCollection = @"raspberries";

NSString * const kPCPrefDefaultTerm  = @"default_terminal";