//
//  MonitorPackage.h
//  manager
//
//  Created by Almighty Kim on 10/22/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "StatusCache.h"

@protocol MonitorPackage <NSObject>
@required
// this show all the available package from api backend
- (void) onAvailableListUpdateWith:(StatusCache *)aCache success:(BOOL)isSuccess error:(NSString *)anErrMsg;

// this show all the installed package in the system
- (void) onInstalledListUpdateWith:(StatusCache *)aCache success:(BOOL)isSuccess error:(NSString *)anErrMsg;
@end