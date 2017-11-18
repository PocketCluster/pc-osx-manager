//
//  AppDelegate+Monitor.h
//  manager
//
//  Created by Almighty Kim on 10/22/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "MonitorAppCheck.h"
#import "MonitorStatus.h"
#import "MonitorPackage.h"
#import "MonitorExecution.h"

#import "AppDelegate.h"

@interface AppDelegate(MonitorDispenser) <MonitorAppCheck, MonitorStatus, MonitorPackage, MonitorExecution>
@end
