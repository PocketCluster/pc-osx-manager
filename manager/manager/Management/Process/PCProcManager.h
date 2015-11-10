//
//  PCProcManager.h
//  manager
//
//  Created by Almighty Kim on 11/2/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCPkgProc.h"

@interface PCProcManager : NSObject
+(instancetype)sharedManager;

- (void)startSalt;
- (void)stopSalt;
- (void)freshSaltStart;

- (void)startWebServer;
- (void)stopWebServer;

- (void)addPackageProcess:(PCPkgProc *)aPackageProcess;
- (void)removePackageProcess:(PCPkgProc *)aPackageProcess;
- (void)refreshPackageProcessesStatus;
- (void)haltRefreshTimer;
- (void)refreshUpdateTimer;
@end
