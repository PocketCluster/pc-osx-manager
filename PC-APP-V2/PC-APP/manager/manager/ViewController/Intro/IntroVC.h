//
//  IntroVC.h
//  manager
//
//  Created by Almighty Kim on 10/19/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "BaseSetupVC.h"

@interface IntroVC : BaseSetupVC
@property (weak) IBOutlet NSTextField *versionLabel;
@property (weak) IBOutlet NSTextField *progressLabel;
@property (weak) IBOutlet NSProgressIndicator *circularProgress;
@end
