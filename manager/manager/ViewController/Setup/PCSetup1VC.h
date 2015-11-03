//
//  PCSetup1VC.h
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "DPSetupWindow.h"

@interface PCSetup1VC : NSViewController  <DPSetupWindowStageViewController>
- (IBAction)setupVagrantCluster:(id)sender;
- (IBAction)setupRaspberryCluster:(id)sender;
@end
