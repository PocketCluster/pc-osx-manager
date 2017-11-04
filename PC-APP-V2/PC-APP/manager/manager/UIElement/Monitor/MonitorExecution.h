//
//  MonitorExecution.h
//  manager
//
//  Created by Almighty Kim on 11/4/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "StatusCache.h"

@protocol MonitorExecution <NSObject>
@required
- (void) onExecutionStartup:(StatusCache *)aCache package:(NSString *)aPackageID;

- (void) didExecutionStartup:(StatusCache *)aCache package:(NSString *)aPackageID success:(BOOL)isSuccess error:(NSString *)anErrMsg;

- (void) onExecutionKill:(StatusCache *)aCache package:(NSString *)aPackageID;

- (void) didExecutionKill:(StatusCache *)aCache package:(NSString *)aPackageID success:(BOOL)isSuccess error:(NSString *)anErrMsg;

- (void) onExecutionProcess:(StatusCache *)aCache package:(NSString *)aPackageID success:(BOOL)isSuccess error:(NSString *)anErrMsg;
@end