//
//  IntroWC.m
//  manager
//
//  Created by Almighty Kim on 10/19/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "StepStage.h"
#import "IntroVC.h"
#import "IntroWC.h"

@interface IntroWC()<StepControl>
@property (nonatomic, strong) IntroVC* viewController;
@end

@implementation IntroWC
- (instancetype) initWithWindowNibName:(NSString *)windowNibName {
    self = [super initWithWindowNibName:windowNibName];
    if (self != nil) {
        self.viewController = [[IntroVC alloc] initWithStageControl:self nibName:@"IntroVC" bundle:[NSBundle mainBundle]];
    }
    return self;
}

- (void) dealloc {
    self.viewController = nil;
}

- (void)windowDidLoad {
    [super windowDidLoad];
    [self.window setBackgroundColor:[NSColor whiteColor]];
    [self.window setTitleVisibility:NSWindowTitleHidden];
    [self.window setTitlebarAppearsTransparent:YES];
    [self.window setShowsResizeIndicator:NO];
    [self.window setShowsToolbarButton:NO];
    self.window.styleMask |= NSFullSizeContentViewWindowMask;

    [[self.window contentView] addSubview:[self.viewController view]];
}

#pragma mark - Stage Control
-(void)shouldControlProgressFrom:(NSObject<StageStep> *)aStep withParam:(NSDictionary *)aParam {
}

-(void)shouldControlRevertFrom:(NSObject<StageStep> *)aStep withParam:(NSDictionary *)aParam {
}
@end
