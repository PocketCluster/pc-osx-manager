//
//  Bookmark.m
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "Raspberry.h"

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


@implementation Raspberry

- (instancetype)initWithDictionary:(NSDictionary *)aDict {
    
    self = [super init];
    if(self){
        
        self.masterBoundAgent   = [aDict objectForKey:MASTER_BOUND_AGENT];
        self.slaveNodeName      = [aDict objectForKey:SLAVE_NODE_NAME];
        self.slaveNodeMacAddr   = [aDict objectForKey:SLAVE_NODE_MACADDR];
        self.address            = [aDict objectForKey:ADDRESS];
        
    }
    
    return self;
}

- (instancetype)initWithCoder:(NSCoder *)aDecoder {

    self = [super init];
    
    if (self){
        
        __attribute__((unused)) NSString *pcVer = [aDecoder decodeObjectForKey:PC_PROTO];
        self.masterBoundAgent   = [aDecoder decodeObjectForKey:MASTER_BOUND_AGENT];
        self.slaveNodeName      = [aDecoder decodeObjectForKey:SLAVE_NODE_NAME];
        self.slaveNodeMacAddr   = [aDecoder decodeObjectForKey:SLAVE_NODE_MACADDR];
        self.address            = [aDecoder decodeObjectForKey:ADDRESS];

    }
    
    return self;
}

- (void)encodeWithCoder:(NSCoder *)anEncoder {
    
    [anEncoder encodeObject:VERSION forKey:PC_PROTO];
    [anEncoder encodeObject:self.masterBoundAgent forKey:MASTER_BOUND_AGENT];
    [anEncoder encodeObject:self.slaveNodeName forKey:SLAVE_NODE_NAME];
    [anEncoder encodeObject:self.slaveNodeMacAddr forKey:SLAVE_NODE_MACADDR];
    [anEncoder encodeObject:self.address forKey:ADDRESS];
    
}

- (id)copyWithZone:(NSZone*)zone {
    Raspberry *rpi = [[[self class] allocWithZone:zone] init];
    
    if(rpi) {
        rpi.masterBoundAgent = self.masterBoundAgent;
        rpi.slaveNodeName = self.slaveNodeName;
        rpi.slaveNodeMacAddr = self.slaveNodeMacAddr;
        rpi.address = self.address;
    }

    return rpi;
}

-(NSString *)description {
    NSString *sd = [super description];
    return [NSString stringWithFormat:@"%@ - %@ (%@) <%@> <- %@",sd, self.slaveNodeName, self.slaveNodeMacAddr, self.address, self.masterBoundAgent];
}


@end
