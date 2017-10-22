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
- (void) onUpdatedWith:(StatusCache *)aCache forPackageListAvailable:(BOOL)isSuccess;

// this show all the installed package in the system
- (void) onUpdatedWith:(StatusCache *)aCache forPackageListInstalled:(BOOL)isSuccess;
@end