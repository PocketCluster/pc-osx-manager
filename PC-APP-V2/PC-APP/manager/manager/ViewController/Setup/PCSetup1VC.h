//
//  PCSetup1VC.h
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "BaseSetupVC.h"

@interface PCSetup1VC : BaseSetupVC
@property (nonatomic, strong) IBOutlet NSView *pannel;
- (IBAction)setup:(id)sender;
- (IBAction)cancel:(id)sender;
@end
