//
//  AppDelegate+Execution.h
//  manager
//
//  Created by Almighty Kim on 11/4/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AppDelegate.h"

@interface AppDelegate(Execution)
- (void) startUpPackageWithID:(NSString *)aPackageID;
- (void) killPackageWithID:(NSString *)aPackageID;
@end
