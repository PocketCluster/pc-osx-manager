//
//  NodeSetupWC.m
//  manager
//
//  Created by Almighty Kim on 11/18/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "BaseSetupVC.h"
#import "PCSetup2VC.h"
#import "PCSetup3VC.h"

#import "NodeSetupWC.h"

@interface NodeSetupWC ()
@property (nonatomic, strong) NSArray<NSViewController<StageStep> *>* viewControllers;
@end

@implementation NodeSetupWC {
    NSUInteger _vcIndex;
}

- (instancetype) initWithWindowNibName:(NSString *)windowNibName {
    self = [super initWithWindowNibName:windowNibName];
    if (self != nil) {
        self.viewControllers =
        @[[[PCSetup2VC alloc] initWithStageControl:self nibName:@"PCSetup2VC" bundle:[NSBundle mainBundle]],
          [[PCSetup3VC alloc] initWithStageControl:self nibName:@"PCSetup3VC" bundle:[NSBundle mainBundle]]];

        // current index
        _vcIndex = 0;

        // prepare the first viewcontroller
        [[self.viewControllers objectAtIndex:_vcIndex] control:self askedProgressWithParam:nil];
    }
    return self;
}

- (void) dealloc {
    [self.viewControllers makeObjectsPerformSelector:@selector(prepareDestruction)];
    self.viewControllers = nil;
}

- (void)windowDidLoad {
    [super windowDidLoad];

    [self.window setTitle:[[self.viewControllers objectAtIndex:0] title]];
    [[self.window contentView] addSubview:[[self.viewControllers objectAtIndex:0] view]];
}

#pragma mark - NSWindowDelegate
- (BOOL)windowShouldClose:(NSWindow *)sender {
    return [(BaseSetupVC *)[self.viewControllers objectAtIndex:_vcIndex] windowShouldClose:sender];
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

    // prepare next stage
    [[self.viewControllers objectAtIndex:nextIndex] control:self askedProgressWithParam:nil];

    // make progress
    [[[self.viewControllers objectAtIndex:prevIndex] view] removeFromSuperview];
    [self.window setTitle:[[self.viewControllers objectAtIndex:nextIndex] title]];
    [[self.window contentView] addSubview:[[self.viewControllers objectAtIndex:nextIndex] view]];

    // notify prev viewController
    [[self.viewControllers objectAtIndex:prevIndex] didControl:self progressedFrom:aStep withResult:nil];
}

-(void)shouldControlRevertFrom:(NSObject<StageStep> *)aStep withParam:(NSDictionary *)aParam {
    
    NSViewController<StageStep> *prevStep = (NSViewController<StageStep> *)aStep;
    NSUInteger prevIndex = [self.viewControllers indexOfObject:prevStep];
    NSUInteger nextIndex = 0;
    
    if (1 <= prevIndex) {
        nextIndex = prevIndex - 1;
    } else {
        Log(@"end of control");
        return;
    }

    // save index
    _vcIndex = nextIndex;

    // prepare next stage
    [[self.viewControllers objectAtIndex:nextIndex] control:self askedRevertWithParam:nil];

    // make progress. this can safe current view states including cursor. but, that's not necessary.
    //[[[self.viewControllers objectAtIndex:prevIndex] view] removeFromSuperviewWithoutNeedingDisplay];
    [[[self.viewControllers objectAtIndex:prevIndex] view] removeFromSuperview];
    [self.window setTitle:[[self.viewControllers objectAtIndex:nextIndex] title]];
    [[self.window contentView] addSubview:[[self.viewControllers objectAtIndex:nextIndex] view]];

    // notify prev viewController
    [[self.viewControllers objectAtIndex:prevIndex] didControl:self revertedFrom:aStep withResult:nil];
}

#pragma mark - Monitoring Package
// this show all the available package from api backend
- (void) onAvailableListUpdateWith:(StatusCache *)aCache success:(BOOL)isSuccess error:(NSString *)anErrMsg {
    for (id<MonitorPackage> vc in self.viewControllers) {
        if ([vc conformsToProtocol:@protocol(MonitorPackage)]) {
            [vc onAvailableListUpdateWith:aCache success:isSuccess error:anErrMsg];
        }
    }
}

// this show all the installed package in the system
- (void) onInstalledListUpdateWith:(StatusCache *)aCache success:(BOOL)isSuccess error:(NSString *)anErrMsg {
    for (id<MonitorPackage> vc in self.viewControllers) {
        if ([vc conformsToProtocol:@protocol(MonitorPackage)]) {
            [vc onInstalledListUpdateWith:aCache success:isSuccess error:anErrMsg];
        }
    }
}

@end
