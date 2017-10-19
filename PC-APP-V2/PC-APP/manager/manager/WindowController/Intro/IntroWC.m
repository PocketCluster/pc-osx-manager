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


#pragma mark - update message
// show initial message
- (void) setupInitialCheckMessage {
    [self.viewController.progressLabel setStringValue:@"Initializing..."];
}

// show "service starting..." message
- (void) setupStartServices {
    [self.viewController.progressLabel setStringValue:@"Starting Services..."];
}

// services online timeup
- (void) onNotifiedWith:(StatusCache *)aCache forServiceOnline:(BOOL)isSuccess {
    [self.viewController.progressLabel setStringValue:@"Checking Nodes..."];
}

// nodes online timeup
- (void) onNotifiedWith:(StatusCache *)aCache forNodeOnline:(BOOL)isSuccess {
    WEAK_SELF(self);
    
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         if(belf){
             [belf close];
         }
     }];
}

// update services
- (void) updateServiceStatusWith:(StatusCache *)aCache {
    
}

// update nodes
- (void) updateNodeStatusWith:(StatusCache *)aCache {
    
    // quickly filter out the worst case scenarios when 'node online timeup' noti has not fired
    if (![[StatusCache SharedStatusCache] showOnlineNode]) {
        if (![[StatusCache SharedStatusCache] isNodeListValid] || \
            ![[StatusCache SharedStatusCache] isAllRegisteredNodesReady]) {
            return;
        }
    }
    
    // -- as 'node online timeup' noti should have been kicked, check strict manner --
    // node list should be valid at this point
    if (![[StatusCache SharedStatusCache] isNodeListValid]) {
        return;
    }
    
    WEAK_SELF(self);

    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         if(belf){
             [belf close];
         }
     }];
}

@end
