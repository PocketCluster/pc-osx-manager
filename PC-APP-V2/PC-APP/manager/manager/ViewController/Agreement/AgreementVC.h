//
//  AgreementVC.h
//  manager
//
//  Created by Almighty Kim on 8/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import <Cocoa/Cocoa.h>

@interface AgreementVC : NSViewController
@property (nonatomic, assign) IBOutlet NSTextView *agreement;

-(IBAction)agreed:(id)sender;
-(IBAction)declined:(id)sender;
@end
