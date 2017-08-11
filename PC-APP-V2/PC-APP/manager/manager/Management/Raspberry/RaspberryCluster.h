//
//  RaspberryCluster.h
//  manager
//
//  Created by Almighty Kim on 11/2/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "Raspberry.h"

@interface RaspberryCluster : NSObject <NSCoding, NSCopying>
@property (nonatomic, strong, readonly) NSString *clusterId;
@property (nonatomic, strong, readonly) NSString *title;
@property (nonatomic, strong, readonly) NSArray *relatedPackages;
@property (nonatomic, readwrite) BOOL swapHasMade;

- (instancetype)initWithTitle:(NSString *)aTitle;
- (void)updateHeartBeats:(NSString *)aMasterId withSlaveMAC:(NSString *)aSlaveMac forTS:(struct timeval)heatbeat;
- (Raspberry*)addRaspberry:(Raspberry*)aRaspberry;
- (NSUInteger)liveRaspberryCount;
- (NSUInteger)raspberryCount;
- (NSMutableArray*)getRaspberries;
- (void)removeRaspberryWithName:(NSString*)aName;
- (Raspberry*)getRaspberryWithName:(NSString*)aName;
- (int)getIndexOfRaspberryWithName:(NSString*)aName;
- (void)checkRelatedPackage;
- (void)resetNodeHeartbeat;

- (void)debugOutput;
@end
