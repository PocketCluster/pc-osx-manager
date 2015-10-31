//
//  Bookmark.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "VirtualMachineInfo.h"

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

@interface Raspberry : NSObject <NSCoding, NSCopying>

-(instancetype)initWithDictionary:(NSDictionary *)aDict;

@property (strong, nonatomic) NSString *masterBoundAgent;
@property (strong, nonatomic) NSString *slaveNodeName;
@property (strong, nonatomic) NSString *slaveNodeMacAddr;
@property (strong, nonatomic) NSString *address;

@end
