//
//  Bookmark.m
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "Raspberry.h"

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
