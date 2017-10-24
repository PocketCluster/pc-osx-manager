//
//  PCPkgInstallWC.m
//  manager
//
//  Created by Almighty Kim on 11/24/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCSetup3VC.h"
#import "PCPkgInstallWC.h"

@interface PCPkgInstallWC ()
@property (nonatomic, strong) PCSetup3VC* viewController;
@end

@implementation PCPkgInstallWC

- (instancetype) initWithWindowNibName:(NSString *)windowNibName {
    self = [super initWithWindowNibName:windowNibName];
    if (self != nil) {
        self.viewController = [[PCSetup3VC alloc] initWithStageControl:self nibName:@"PCSetup3VC" bundle:[NSBundle mainBundle]];
    }
    return self;
}

- (void) dealloc {
    self.viewController = nil;
}

- (void)windowDidLoad {
    [super windowDidLoad];
    
    [self.window setTitle:[self.viewController title]];
    [[self.window contentView] addSubview:[self.viewController view]];
}

#pragma mark - NSWindowDelegate
- (BOOL)windowShouldClose:(NSWindow *)sender {
    return [self.viewController windowShouldClose:sender];
}

#pragma mark - Stage Control
-(void)shouldControlProgressFrom:(NSObject<StageStep> *)aStep withParam:(NSDictionary *)aParam {
}

-(void)shouldControlRevertFrom:(NSObject<StageStep> *)aStep withParam:(NSDictionary *)aParam {
}

@end
