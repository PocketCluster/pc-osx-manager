//
//  UserCheckVC.h
//  manager
//
//  Created by Almighty Kim on 8/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import <Cocoa/Cocoa.h>

@interface UserCheckVC : NSViewController
@property (nonatomic, weak) IBOutlet NSProgressIndicator *progress;

-(IBAction)check:(id)sender;
-(IBAction)cancel:(id)sender;
@end
