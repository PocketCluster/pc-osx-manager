//
//  AppDelegate+Monitor.h
//  manager
//
//  Created by Almighty Kim on 10/22/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AppDelegate.h"
#import "MonitorStatus.h"
#import "MonitorPackage.h"

@interface AppDelegate(Monitor) <MonitorStatus, MonitorPackage>
@end
