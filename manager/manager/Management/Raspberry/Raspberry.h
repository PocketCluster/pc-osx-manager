//
//  Bookmark.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "VirtualMachineInfo.h"
#import "PCConstants.h"
#include <time.h>

#define HEARTBEAT_CHECK_INTERVAL (30.0)


@interface Raspberry : NSObject <NSCoding, NSCopying>

@property (strong, nonatomic) NSString *masterBoundAgent;
@property (strong, nonatomic) NSString *slaveNodeName;
@property (strong, nonatomic) NSString *slaveNodeMacAddr;
@property (strong, nonatomic) NSString *address;
@property (nonatomic, readonly) BOOL isAlive;
@property (nonatomic, assign) struct timeval heartbeat;	// seconds and microseconds

-(instancetype)initWithDictionary:(NSDictionary *)aDict;

@end
