//
//  Package
//  PocketCluster
//
//  Created by Almighty Kim on 11/4/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

typedef enum execState {
    ExecIdle = 0,
    ExecStarting,
    ExecRun,
    ExecStopping
} ExecState;

@interface Package : NSObject

@property (nonatomic, readonly) BOOL installed;
@property (nonatomic, readonly) NSString *packageDescription;
@property (nonatomic, readonly) NSString *packageID;
@property (nonatomic, readonly) NSString *menuName;
@property (nonatomic, readonly) ExecState state;

// all three below are in GB
@property (nonatomic, readonly) NSString *totalImageSize;

+ (NSArray<Package *> *)packagesFromList:(NSArray<NSDictionary *> *)aList;

// update all package status except `package id` with a package w/ same id.
- (void) updateWithPackage:(Package *)newSamePackage;
- (void) updateExecState:(ExecState)state;
@end
