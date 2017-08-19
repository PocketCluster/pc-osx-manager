//
//  AgreementVC.m
//  manager
//
//  Created by Almighty Kim on 8/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AgreementVC.h"

@implementation AgreementVC

- (void)finishConstruction {
    [super finishConstruction];
    [self setTitle:@"PocketCluster User Agreement"];
}

- (void)viewDidLoad {
    [super viewDidLoad];

    [[((BaseBrandView *)self.view) contentBox] removeFromSuperview];
    ((BaseBrandView *)self.view).contentBox = nil;

    [self.agreement setVerticallyResizable:YES];
    [self.agreement setHorizontallyResizable:NO];
    [self.agreement setEditable:NO];
    
    [self.agreement.textStorage
     setAttributedString:[[NSAttributedString alloc]
                          initWithPath:[[NSBundle mainBundle] pathForResource:@"EULA" ofType:@"rtf"]
                          documentAttributes:nil]];
}

-(IBAction)agreed:(id)sender {
    [self.stageControl shouldControlProgressFrom:self withParam:nil];
}

-(IBAction)declined:(id)sender {
    [self.stageControl shouldControlRevertFrom:self withParam:nil];
}

#pragma mark - StageStep
-(void)didControl:(NSObject<StepControl> *)aControl progressFrom:(NSObject<StageStep> *)aStep withResult:(NSDictionary *)aResult {
    if (aStep != self) {
        Log(@"this is not identical step!");
        return;
    }
}

-(void)didControl:(NSObject<StepControl> *)aControl revertFrom:(NSObject<StageStep> *)aStep withResult:(NSDictionary *)aResult {
    if (aStep != self) {
        Log(@"this is not identical step!");
        return;
    }
}
@end
