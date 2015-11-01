//
//  Bookmark.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "VirtualMachineInfo.h"
#import "PCConstants.h"

@interface Raspberry : NSObject <NSCoding, NSCopying>

-(instancetype)initWithDictionary:(NSDictionary *)aDict;

@property (strong, nonatomic) NSString *masterBoundAgent;
@property (strong, nonatomic) NSString *slaveNodeName;
@property (strong, nonatomic) NSString *slaveNodeMacAddr;
@property (strong, nonatomic) NSString *address;
@property (readwrite, getter=isAlive) BOOL alive;

@end
