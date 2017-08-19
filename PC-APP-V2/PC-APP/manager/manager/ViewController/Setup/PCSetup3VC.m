//
//  PCSetup3VC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "Package.h"
#import "PCSetup3VC.h"

@interface PCSetup3VC()
@property (nonatomic, strong) NSMutableArray<Package *> *packageList;
@end

@implementation PCSetup3VC

- (void) finishConstruction {
    [super finishConstruction];
    [self setTitle:@"Install Package"];
    self.packageList = [NSMutableArray arrayWithCapacity:0];
}

- (void) viewDidLoad {
    [super viewDidLoad];
    
    [[((BaseBrandView *)self.view) contentBox] addSubview:self.pannel];
    self.pannel = nil;
}

#pragma mark - NSTableViewDataSourceDelegate
- (NSInteger)numberOfRowsInTableView:(NSTableView *)tableView {
    return [self.packageList count];
}

- (nullable id)tableView:(NSTableView *)tableView objectValueForTableColumn:(nullable NSTableColumn *)tableColumn row:(NSInteger)row {
    return [[self.packageList objectAtIndex:row] description];
}

#pragma mark - NSTableViewDelegate
-(NSView *)tableView:(NSTableView *)aTableView viewForTableColumn:(NSTableColumn *)aTableColumn row:(NSInteger)row{
    Package *meta = [self.packageList objectAtIndex:row];
    NSTableCellView *nv = [aTableView makeViewWithIdentifier:@"packageview" owner:self];
    [nv.textField setStringValue:[meta packageDescription]];
    return nv;
}

- (BOOL)selectionShouldChangeInTableView:(NSTableView *)tableView {
    return NO;
}

- (BOOL)tableView:(NSTableView *)tableView shouldSelectRow:(NSInteger)row {
    return NO;
}

#pragma mark - Setup UI states
- (void)_setUIToProceedState {
    [self.btnInstall setEnabled:NO];
    [self.circularProgress startAnimation:nil];
}

-(void)_enableControls {
    [self.btnInstall setEnabled:YES];
    [self.btnCancel setEnabled:YES];
    [self.progressLabel setStringValue:@"Installation Error. Please try again."];
    [self.circularProgress stopAnimation:nil];
    [self.progressBar setDoubleValue:0.0];
    [self.progressBar displayIfNeeded];
}

-(void)_setProgressMessage:(NSString *)aMessage value:(double)aValue {
    [self.circularProgress startAnimation:nil];
    [self.progressLabel setStringValue:aMessage];
    [self.progressBar setDoubleValue:aValue];
    [self.progressBar displayIfNeeded];
}

-(void)setProgMessage:(NSString *)aMessage value:(double)aValue {
    [self.circularProgress startAnimation:nil];
    [self.progressLabel setStringValue:aMessage];
    [self.progressBar setDoubleValue:aValue];
    [self.progressBar displayIfNeeded];
}

-(void)_setToNextStage {
    [self setProgMessage:@"Installation completed!" value:100.0];
    [self.btnInstall setEnabled:NO];
    [self.circularProgress stopAnimation:nil];
}

#pragma mark - IBACTION
-(IBAction)install:(id)sender {
    [self.stageControl shouldControlProgressFrom:self withParam:nil];
}

-(IBAction)cancel:(id)sender {
    [self.stageControl shouldControlRevertFrom:self withParam:nil];    
}
@end
