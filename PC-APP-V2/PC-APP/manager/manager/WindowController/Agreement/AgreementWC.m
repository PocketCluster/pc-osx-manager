//
//  AgreementWC.m
//  manager
//
//  Created by Almighty Kim on 8/16/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AgreementWC.h"
#import "AgreementVC.h"
#import "UserCheckVC.h"

@interface AgreementWC ()
@property (nonatomic, strong) NSArray<NSViewController *>* viewControllers;
@end

@implementation AgreementWC

- (instancetype) initWithWindowNibName:(NSString *)windowNibName {
    self = [super initWithWindowNibName:windowNibName];
    if (self != nil) {
        self.viewControllers =
            @[[[AgreementVC alloc] initWithNibName:@"AgreementVC" bundle:[NSBundle mainBundle]],
              [[UserCheckVC alloc] initWithNibName:@"UserCheckVC" bundle:[NSBundle mainBundle]]];
    }
    return self;
}

- (void) dealloc {
    self.viewControllers = nil;
}

- (void)windowDidLoad {
    [super windowDidLoad];
    
    [self.window setTitle:[[self.viewControllers objectAtIndex:1] title]];
    [[self.window contentView] addSubview:[[self.viewControllers objectAtIndex:1] view]];
}

@end
