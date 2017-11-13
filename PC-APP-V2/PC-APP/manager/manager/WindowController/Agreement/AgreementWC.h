//
//  AgreementWC.h
//  manager
//
//  Created by Almighty Kim on 8/16/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "BaseWindowController.h"
#import "StepStage.h"
#import "MonitorAppCheck.h"
#import "MonitorStatus.h"

@interface AgreementWC : BaseWindowController <MonitorAppCheck, MonitorStatus, StepControl>
@end
