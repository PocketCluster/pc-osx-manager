//
//  PCSetup1VC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "BFPageControl.h"
#import "PCSetup1VC.h"

#define NUM_INTRO_PAGES  3

@interface PCSetup1VC (Private)<BFPageControlDelegate>
@end

@implementation PCSetup1VC {
    __strong NSArray<NSString *>* _introImages;
}

- (void) finishConstruction {
    [super finishConstruction];
    [self setTitle:@"Setup a New Cluster"];
    _introImages = @[@"intro_01", @"intro_02", @"intro_03"];
}

- (void) viewDidLoad {
    [super viewDidLoad];

    // Setup page control
    NSRect frame = self.pannel.frame;
    BFPageControl *control = [[BFPageControl alloc] init];
    [control setDelegate: self];
    [control setNumberOfPages: NUM_INTRO_PAGES];
    [control setIndicatorDiameterSize: 15];
    [control setIndicatorMargin:5];
    [control setCurrentPage: 0];
    [control setDrawingBlock: ^(NSRect frame, NSView *aView, BOOL isSelected, BOOL isHighlighted){

        frame = CGRectInset(frame, 2.0, 2.0);
        NSBezierPath *path = [NSBezierPath bezierPathWithOvalInRect: CGRectMake(frame.origin.x, frame.origin.y + 1.5, frame.size.width, frame.size.height)];
        [[NSColor whiteColor] set];
        [path fill];

        path = [NSBezierPath bezierPathWithOvalInRect: frame];
        NSColor *color = isSelected ? [NSColor colorWithCalibratedRed: (115.0 / 255.0) green: (115.0 / 255.0) blue: (115.0 / 255.0) alpha: 1.0] :
        [NSColor colorWithCalibratedRed: (217.0 / 255.0) green: (217.0 / 255.0) blue: (217.0 / 255.0) alpha: 1.0];

        if(isHighlighted) {
            color = [NSColor colorWithCalibratedRed: (150.0 / 255.0) green: (150.0 / 255.0) blue: (150.0 / 255.0) alpha: 1.0];
        }

        [color set];
        [path fill];

        frame = CGRectInset(frame, 0.5, 0.5);
        [[NSColor colorWithCalibratedRed: (25.0 / 255.0) green: (25.0 / 255.0) blue: (25.0 / 255.0) alpha: 0.15] set];
        [NSBezierPath setDefaultLineWidth: 1.0];
        [[NSBezierPath bezierPathWithOvalInRect: frame] stroke];
    }];
    [self.pannel addSubview: control];
    CGSize size = [control intrinsicContentSize];
    [control setFrame: CGRectMake((frame.size.width - size.width)/2, 50, size.width, size.height)];

    [[((BaseBrandView *)self.view) contentBox] addSubview:self.pannel];
    self.pannel = nil;

    [self _disableControls];
}

#pragma mark - 
-(void)pageControl: (BFPageControl *)pageControl didSelectPageAtIndex: (NSInteger)index {
    [self.introImage setImage:[NSImage imageNamed:[_introImages objectAtIndex:index]]];
}

- (IBAction)setup:(id)sender {
    [self.stageControl shouldControlProgressFrom:self withParam:nil];
}

- (IBAction)cancel:(id)sender {
    [self.stageControl shouldControlRevertFrom:self withParam:nil];
}

-(void) enableControls {
    [self.circularProgress setHidden:YES];
    [self.circularProgress stopAnimation:nil];
    [self.circularProgress displayIfNeeded];
    [self.circularProgress removeFromSuperview];
    [self setCircularProgress:nil];

    [self.btnCancel setEnabled:YES];
    [self.btnSetup setEnabled:YES];
}

-(void) _disableControls {
    [self.btnCancel setEnabled:NO];
    [self.btnSetup setEnabled:NO];

    NSProgressIndicator *ind = [[NSProgressIndicator alloc] initWithFrame:(NSRect){{20.0, 20.0}, {16.0, 16.0}}];
    [ind setControlSize:NSSmallControlSize];
    [ind setStyle:NSProgressIndicatorSpinningStyle];
    [self.view addSubview:ind];
    [ind setHidden:NO];
    [ind setIndeterminate:YES];
    [ind startAnimation:self];
    [ind displayIfNeeded];

    [self setCircularProgress:ind];
}

@end
