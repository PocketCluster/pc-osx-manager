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
@property (nonatomic, strong) NSArray<NSViewController<StageStep> *>* viewControllers;
@end

@implementation AgreementWC

- (instancetype) initWithWindowNibName:(NSString *)windowNibName {
    self = [super initWithWindowNibName:windowNibName];
    if (self != nil) {
        self.viewControllers =
            @[[[AgreementVC alloc] initWithStageControl:self nibName:@"AgreementVC" bundle:[NSBundle mainBundle]],
              [[UserCheckVC alloc] initWithStageControl:self nibName:@"UserCheckVC" bundle:[NSBundle mainBundle]]];
    }
    return self;
}

- (void) dealloc {
    self.viewControllers = nil;
}

- (void)windowDidLoad {
    [super windowDidLoad];
    
    [self.window setTitle:[[self.viewControllers objectAtIndex:0] title]];
    [[self.window contentView] addSubview:[[self.viewControllers objectAtIndex:0] view]];
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

    [[[self.viewControllers objectAtIndex:prevIndex] view] removeFromSuperview];
    [self.window setTitle:[[self.viewControllers objectAtIndex:nextIndex] title]];
    [[self.window contentView] addSubview:[[self.viewControllers objectAtIndex:nextIndex] view]];

    [[self.viewControllers objectAtIndex:prevIndex] didControl:self progressFrom:aStep withResult:nil];
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

    // this can safe current view states including cursor. but, that's not necessary.
    //[[[self.viewControllers objectAtIndex:prevIndex] view] removeFromSuperviewWithoutNeedingDisplay];
    [[[self.viewControllers objectAtIndex:prevIndex] view] removeFromSuperview];
    [self.window setTitle:[[self.viewControllers objectAtIndex:nextIndex] title]];
    [[self.window contentView] addSubview:[[self.viewControllers objectAtIndex:nextIndex] view]];
    
    [[self.viewControllers objectAtIndex:prevIndex] didControl:self progressFrom:aStep withResult:nil];
}

@end
