//
//  PCPkgInstallWC.h
//  manager
//
//  Created by Almighty Kim on 11/24/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "BaseWindowController.h"
#import "StepStage.h"
#import "MonitorPackage.h"

@interface PCPkgInstallWC : BaseWindowController <NSWindowDelegate, StepControl, MonitorPackage>
@end
