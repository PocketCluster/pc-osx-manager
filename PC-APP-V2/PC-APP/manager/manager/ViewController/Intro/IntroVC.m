//
//  IntroVC.m
//  manager
//
//  Created by Almighty Kim on 10/19/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "IntroVC.h"

@implementation IntroVC
@synthesize versionLabel;
@synthesize progressLabel;
@synthesize circularProgress;

- (void)viewDidLoad {
    [super viewDidLoad];

    NSString *vStr = [NSString stringWithFormat:@"Version %@ - Early Evaluation",[[[NSBundle mainBundle] infoDictionary] valueForKey:@"CFBundleShortVersionString"]];
    [self.versionLabel setStringValue:vStr];
    [self.circularProgress startAnimation:self];
    [self.circularProgress displayIfNeeded];
}

@end
