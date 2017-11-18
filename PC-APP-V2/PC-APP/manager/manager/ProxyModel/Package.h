//
//  Package
//  PocketCluster
//
//  Created by Almighty Kim on 11/4/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

/*
 * package process status is self-explanatory. -ing state means it's in transition
 * between ui command -> package is starting or stopping
 */
typedef enum execState {
    ExecIdle    = 0,
    ExecStarting,
    ExecStarted,
    ExecRun,
    ExecStopping
    // ExecStopped = ExecIdle
} ExecState;

@interface Package : NSObject

@property (nonatomic, readonly) BOOL installed;
@property (nonatomic, readonly) NSString *packageDescription;
@property (nonatomic, readonly) NSString *packageID;
@property (nonatomic, readonly) NSString *menuName;
@property (nonatomic, readonly) ExecState execState;

// all three below are in GB
@property (nonatomic, readonly) NSString *totalImageSize;

+ (NSArray<Package *> *)packagesFromList:(NSArray<NSDictionary *> *)aList;

// update all package status except `package id` with a package w/ same id.
- (void) updateWithPackage:(Package *)newSamePackage;
- (void) updateExecState:(ExecState)state;
@end
