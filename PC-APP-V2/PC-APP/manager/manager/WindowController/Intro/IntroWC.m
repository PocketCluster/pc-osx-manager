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
    [self.window setTitleVisibility:NSWindowTitleVisible];
    [self.window setTitlebarAppearsTransparent:YES];
    [self.window setShowsResizeIndicator:NO];
    [self.window setShowsToolbarButton:NO];
    self.window.styleMask |= NSFullSizeContentViewWindowMask;

/*
    // (2017/10/24) this isn't strictly necessary as of now
    [self.window setTitle:
     [@"PocketCluster {VERSION} - Early Evaluation"
      stringByReplacingOccurrencesOfString:@"{VERSION}"
      withString:[[[NSBundle mainBundle] infoDictionary] valueForKey:@"CFBundleShortVersionString"]]];
*/

    [[self.window contentView] addSubview:[self.viewController view]];

    // begin with initialization method
    [self.viewController.progressLabel setStringValue:@"Initializing..."];
}

#pragma mark - Stage Control
-(void)shouldControlProgressFrom:(NSObject<StageStep> *)aStep withParam:(NSDictionary *)aParam {
}

-(void)shouldControlRevertFrom:(NSObject<StageStep> *)aStep withParam:(NSDictionary *)aParam {
}


#pragma mark - update message
// show initial message
- (void) setupWithInitialCheckMessage {
    [self.viewController.progressLabel setStringValue:@"Initializing..."];
}

// show "service starting..." message
- (void) setupWithStartServicesMessage {
    [self.viewController.progressLabel setStringValue:@"Starting Services..."];
}

- (void) setupWithCheckingNodesMessage {
    [self.viewController.progressLabel setStringValue:@"Checking Nodes..."];
}

// services online timeup
- (void) onNotifiedWith:(StatusCache *)aCache serviceOnlineTimeup:(BOOL)isSuccess {

    if ([aCache serviceError] != nil) {
        WEAK_SELF(self);
        [[NSOperationQueue mainQueue]
         addOperationWithBlock:^{
             if(belf){
                 [belf close];
             }
         }];
    }
}

// nodes online timeup
- (void) onNotifiedWith:(StatusCache *)aCache nodeOnlineTimeup:(BOOL)isSuccess {
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
    if (![[StatusCache SharedStatusCache] timeUpNodeOnline]) {
        if (![[StatusCache SharedStatusCache] isNodeListValid] || \
            ![[StatusCache SharedStatusCache] isRegisteredNodesAllOnline]) {
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
