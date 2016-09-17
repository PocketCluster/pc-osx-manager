//
//  RaspberryCluster.m
//  manager
//
//  Created by Almighty Kim on 11/2/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "RaspberryCluster.h"
#import "DeviceSerialNumber.h"
#import "NullStringChecker.h"
#import "PCPackageManager.h"

@interface RaspberryCluster()
@property (nonatomic, strong, readwrite) NSString *clusterId;
@property (nonatomic, strong, readwrite) NSString *title;
@property (nonatomic, strong, readwrite) NSMutableArray *raspberries;
@property (nonatomic, strong, readwrite) NSArray *relatedPackages;
@end

@implementation RaspberryCluster

- (instancetype)initWithTitle:(NSString *)aTitle {
    self = [super init];
    if(self){
        self.title = aTitle;
        self.clusterId = [DeviceSerialNumber UUIDString];
        self.raspberries = [NSMutableArray arrayWithCapacity:0];
        self.swapHasMade = NO;
    }
    return self;
}

- (instancetype)initWithCoder:(NSCoder *)aDecoder {
    
    self = [super init];
    if (self){
        
        NSString *cid = [aDecoder decodeObjectForKey:kRaspberryClusterId];
        NSString *ttl = [aDecoder decodeObjectForKey:kRaspberryClusterTitle];
        NSArray *rpis = [aDecoder decodeObjectForKey:kRaspberryClusterArray];
        BOOL swapon   = [[aDecoder decodeObjectForKey:kRaspberryClusterSwapOn] boolValue];

        self.clusterId = cid;
        self.raspberries = [rpis mutableCopy];
        if(!ISNULL_STRING(ttl)){
            self.title = ttl;
        }
        self.swapHasMade = swapon;
    }
    return self;
}

- (void)encodeWithCoder:(NSCoder *)anEncoder {
    [anEncoder encodeObject:_clusterId              forKey:kRaspberryClusterId];
    if(!ISNULL_STRING(_title)){
        [anEncoder encodeObject:_title              forKey:kRaspberryClusterTitle];
    }
    [anEncoder encodeObject:_raspberries            forKey:kRaspberryClusterArray];
    [anEncoder encodeObject:@(_swapHasMade)         forKey:kRaspberryClusterSwapOn];
}

- (id)copyWithZone:(NSZone*)zone {
    RaspberryCluster *rpiclsuter = [[[self class] allocWithZone:zone] init];
    
    if(rpiclsuter) {
        rpiclsuter.clusterId        = self.clusterId;
        rpiclsuter.title            = self.title;
        rpiclsuter.raspberries      = self.raspberries;
        rpiclsuter.swapHasMade      = self.swapHasMade;
    }
    return rpiclsuter;
}

- (void)updateHeartBeats:(NSString *)aMasterId withSlaveMAC:(NSString *)aSlaveMac forTS:(struct timeval)heartbeat {
    // check heartbeat
    @synchronized(_raspberries) {
        [_raspberries enumerateObjectsUsingBlock:^(Raspberry*  _Nonnull rpi, NSUInteger idx, BOOL * _Nonnull stop) {
            if([rpi.masterBoundAgent isEqualToString:aMasterId] && [rpi.slaveNodeMacAddr isEqualToString:aSlaveMac]){
                rpi.heartbeat = heartbeat;
            }
        }];
    }
}

- (Raspberry*)addRaspberry:(Raspberry*)aRaspberry {
    Raspberry *existing = [self getRaspberryWithName:aRaspberry.slaveNodeName];
    
    if(existing) {
        return existing;
    }
    
    @synchronized(_raspberries) {
        [_raspberries addObject:aRaspberry];
    }
    
    return aRaspberry;
}

- (NSUInteger)liveRaspberryCount {
    NSArray *filtered = [_raspberries filteredArrayUsingPredicate:[NSPredicate predicateWithFormat:@"(SELF.isAlive == YES)"]];
    return [filtered count];
}

- (NSUInteger)raspberryCount {
    return [self.raspberries count];
}

- (NSMutableArray*)getRaspberries {
    NSMutableArray *rpicluster;
    @synchronized(_raspberries) {
        rpicluster = [NSMutableArray arrayWithArray:_raspberries];
    }
    return rpicluster;
}

- (void)removeRaspberryWithName:(NSString*)aName {
    Raspberry *bookmark = [self getRaspberryWithName:aName];
    if(bookmark) {
        @synchronized(_raspberries) {
            [_raspberries removeObject:bookmark];
        }
    }
}

- (Raspberry*)getRaspberryWithName:(NSString*)aName {
    @synchronized(_raspberries) {
        for(Raspberry *rpi in _raspberries) {
            if([rpi.slaveNodeName isEqualToString:aName]) {
                return rpi;
            }
        }
    }
    
    return nil;
}

- (int)getIndexOfRaspberryWithName:(NSString*)aName {
    for(int i=0; i<_raspberries.count; ++i) {
        Raspberry *rpi = [_raspberries objectAtIndex:i];
        if([rpi.slaveNodeName isEqualToString:aName]) {
            return i;
        }
    }
    
    return -1;
}

- (void)checkRelatedPackage {
    self.relatedPackages =
        [[[PCPackageManager sharedManager] installedPackage]
         filteredArrayUsingPredicate:
            [NSPredicate predicateWithFormat:@"(SELF.clusterRelation == %@)",self.clusterId]];
}

- (void)resetNodeHeartbeat {
    @synchronized(_raspberries) {
        [_raspberries makeObjectsPerformSelector:@selector(resetHeartbeat)];
    }
}

-(void)debugOutput {
    Log(@"TITLE %@ CID %@ ", _title, _clusterId);
    for (Raspberry *rpi in _raspberries){
        Log(@"\t%@", [rpi description]);
    }
}

@end