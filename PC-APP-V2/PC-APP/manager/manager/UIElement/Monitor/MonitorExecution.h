//
//  MonitorExecution.h
//  manager
//
//  Created by Almighty Kim on 11/4/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "StatusCache.h"
#import "Package.h"

@protocol MonitorExecution <NSObject>
@required
- (void) onExecutionStartup:(Package *)aPackage;

- (void) didExecutionStartup:(Package *)aPackage success:(BOOL)isSuccess error:(NSString *)anErrMsg;

- (void) onExecutionKill:(Package *)aPackage;

- (void) didExecutionKill:(Package *)aPackage success:(BOOL)isSuccess error:(NSString *)anErrMsg;

- (void) onExecutionProcess:(Package *)aPackage success:(BOOL)isSuccess error:(NSString *)anErrMsg;
@end