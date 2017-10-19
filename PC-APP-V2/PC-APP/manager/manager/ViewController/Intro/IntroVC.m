//
//  IntroVC.m
//  manager
//
//  Created by Almighty Kim on 10/19/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "IntroVC.h"

@implementation IntroVC

- (void)finishConstruction {
    [super finishConstruction];
    [self setTitle:@"PocketCluster - v0.1.4 Early Evaulation Version"];
}

- (void)viewDidLoad {
    [super viewDidLoad];
    
    [[((BaseBrandView *)self.view) contentBox] removeFromSuperview];
    ((BaseBrandView *)self.view).contentBox = nil;
}

@end
