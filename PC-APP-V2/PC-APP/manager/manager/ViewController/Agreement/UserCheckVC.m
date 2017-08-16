//
//  UserCheckVC.m
//  manager
//
//  Created by Almighty Kim on 8/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "UserCheckVC.h"
#import "BaseBrandView.h"

@interface UserCheckVC ()
@end

@implementation UserCheckVC

- (instancetype) initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil {
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if (self != nil) {
        [self setTitle:@"Inviatation Check"];
    }
    return self;
}

- (void)viewDidLoad {
    [super viewDidLoad];

    [[((BaseBrandView *)self.view) contentBox] removeFromSuperview];
    ((BaseBrandView *)self.view).contentBox = nil;
}

-(IBAction)check:(id)sender {
    
}

-(IBAction)cancel:(id)sender {
    
}

@end
