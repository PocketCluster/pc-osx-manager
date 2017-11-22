//
//  AgreementWC.m
//  manager
//
//  Created by Almighty Kim on 8/16/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "PCConstants.h"
#import "AgreementVC.h"
#import "UserCheckVC.h"
#import "PCSetup1VC.h"
#import "PCSetup2VC.h"
#import "PCSetup3VC.h"

#import "AgreementWC.h"

enum {
    SETUPVC_AGREEMENT = 0,
    SETUPVC_USER_CHECK,
    SETUPVC_INTRO,
    SETUPVC_BUILD_CLUSTER,
    SETUPVC_INSTALL_PKG,
    SETUPVC_COUNT,
};

@interface AgreementWC ()
@property (nonatomic, strong) NSArray<NSViewController<StageStep> *>* viewControllers;
@end

@implementation AgreementWC {
    NSUInteger _vcIndex;
}

- (instancetype) initWithWindowNibName:(NSString *)windowNibName {
    self = [super initWithWindowNibName:windowNibName];
    if (self != nil) {
        self.viewControllers =
            @[[[AgreementVC alloc] initWithStageControl:self nibName:@"AgreementVC" bundle:[NSBundle mainBundle]],
              [[UserCheckVC alloc] initWithStageControl:self nibName:@"UserCheckVC" bundle:[NSBundle mainBundle]],
              [[PCSetup1VC alloc] initWithStageControl:self nibName:@"PCSetup1VC" bundle:[NSBundle mainBundle]],
              [[PCSetup2VC alloc] initWithStageControl:self nibName:@"PCSetup2VC" bundle:[NSBundle mainBundle]],
              [[PCSetup3VC alloc] initWithStageControl:self nibName:@"PCSetup3VC" bundle:[NSBundle mainBundle]]];

        // current index
        _vcIndex = 0;
    }
    return self;
}

- (void) dealloc {
    [self.viewControllers makeObjectsPerformSelector:@selector(prepareDestruction)];
    self.viewControllers = nil;
}

- (void)windowDidLoad {
    [super windowDidLoad];
    
    [self.window setTitle:[[self.viewControllers objectAtIndex:_vcIndex] title]];
    [[self.window contentView] addSubview:[[self.viewControllers objectAtIndex:_vcIndex] view]];
}

#pragma mark - Stage Control
-(void)shouldControlProgressFrom:(NSObject<StageStep> *)aStep withParam:(NSDictionary *)aParam {
    
    NSViewController<StageStep> *prevStep = (NSViewController<StageStep> *)aStep;
    NSUInteger prevIndex = [self.viewControllers indexOfObject:prevStep];
    NSUInteger nextIndex = 0;

    if (prevIndex < ([self.viewControllers count] - 1)) {
        nextIndex = prevIndex + 1;
    } else {
        Log(@"end of control");
        return;
    }

    // save index
    _vcIndex = nextIndex;

    [[[self.viewControllers objectAtIndex:prevIndex] view] removeFromSuperview];
    [self.window setTitle:[[self.viewControllers objectAtIndex:nextIndex] title]];
    [[self.window contentView] addSubview:[[self.viewControllers objectAtIndex:nextIndex] view]];

    [[self.viewControllers objectAtIndex:prevIndex] didControl:self progressFrom:aStep withResult:nil];
}

-(void)shouldControlRevertFrom:(NSObject<StageStep> *)aStep withParam:(NSDictionary *)aParam {
    NSViewController<StageStep> *prevStep = (NSViewController<StageStep> *)aStep;
    NSUInteger prevIndex = [self.viewControllers indexOfObject:prevStep];

    if (prevIndex <= SETUPVC_USER_CHECK) {
        [[NSApplication sharedApplication] terminate:nil];
    } else {
        [self close];
    }

#ifdef _REWIND_PAGES_
    NSViewController<StageStep> *prevStep = (NSViewController<StageStep> *)aStep;
    NSUInteger prevIndex = [self.viewControllers indexOfObject:prevStep];
    NSUInteger nextIndex = 0;

    if (1 <= prevIndex) {
        nextIndex = prevIndex - 1;
    } else {
        Log(@"end of control");
        return;
    }

    // this can safe current view states including cursor. but, that's not necessary.
    //[[[self.viewControllers objectAtIndex:prevIndex] view] removeFromSuperviewWithoutNeedingDisplay];
    [[[self.viewControllers objectAtIndex:prevIndex] view] removeFromSuperview];
    [self.window setTitle:[[self.viewControllers objectAtIndex:nextIndex] title]];
    [[self.window contentView] addSubview:[[self.viewControllers objectAtIndex:nextIndex] view]];

    [[self.viewControllers objectAtIndex:prevIndex] didControl:self progressFrom:aStep withResult:nil];
#endif
}


#pragma mark - MonitorAppCheck

// check system readiness
- (void) didAppCheckSystemReadiness:(BOOL)isReady {
}

// check app has been expried
- (void) didAppCheckAppExpiration:(BOOL)isExpired {
}

// check if first time run
- (void) didAppCheckIsFirstRun:(BOOL)isFirstRun {
}

// check if user is authed
- (void) didAppCheckUserAuthed:(BOOL)isUserAuthed {
    if (isUserAuthed) {
        [self shouldControlProgressFrom:[self.viewControllers objectAtIndex:1] withParam:nil];
    } else {
        [(UserCheckVC *)[self.viewControllers objectAtIndex:1] enableControls];
    }
}

#pragma mark - MonitorStatus
// show initial message
- (void) setupWithInitialCheckMessage {
    [[(PCSetup1VC *)[self.viewControllers objectAtIndex:2] progressLabel] setStringValue:kAppCheckMessageInit];
}

// show "service starting" message.
- (void) setupWithStartServicesMessage {
    [[(PCSetup1VC *)[self.viewControllers objectAtIndex:2] progressLabel] setStringValue:kAppCheckMessageSrvcStart];
}

// services online timeup. Display service status. This is paired method that
// needs to be initiated by previous call to `setupWithStartServicesMessage`
- (void) onNotifiedWith:(StatusCache *)aCache serviceOnlineTimeup:(BOOL)isSuccess {
}

// show "checking nodes" message
- (void) setupWithCheckingNodesMessage {
    [[(PCSetup1VC *)[self.viewControllers objectAtIndex:2] progressLabel] setStringValue:kAppCheckMessageCheckingNode];
}

// nodes online timeup. Display node state no matter what. This is paired method that
// needs to be initiated by previous call to `setupWithCheckingNodesMessage`
- (void) onNotifiedWith:(StatusCache *)aCache nodeOnlineTimeup:(BOOL)isSuccess {
}

// update services
- (void) updateServiceStatusWith:(StatusCache *)aCache {
}

// update nodes
- (void) updateNodeStatusWith:(StatusCache *)aCache {
    if (![aCache activateMenuBeforeNodeTimeup]) {
        return;
    }

    [[(PCSetup1VC *)[self.viewControllers objectAtIndex:2] progressLabel] setStringValue:@"Ready to Setup a Raspberry Pi cluster"];
}

@end
