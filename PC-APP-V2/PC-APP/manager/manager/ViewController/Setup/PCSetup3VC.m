//
//  PCSetup3VC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "Package.h"
#import "PCSetup3VC.h"
#import "PCRouter.h"
#import "ShowAlert.h"

@interface PCSetup3VC()<PCRouteRequest>
@property (nonatomic, strong) NSMutableArray<Package *> *packageList;
@end

@implementation PCSetup3VC {
    NSInteger _selectedIndex;
}

- (void) finishConstruction {
    [super finishConstruction];
    [self setTitle:@"Install Package"];
    self.packageList = [NSMutableArray<Package *> arrayWithCapacity:0];
}

- (void) viewDidLoad {
    [super viewDidLoad];
    
    [[((BaseBrandView *)self.view) contentBox] addSubview:self.pannel];
    self.pannel = nil;
}

- (void) viewDidAppear {
    [super viewDidAppear];
    
    // reset selected index
    _selectedIndex = -1;
    
    /*** checking user authed ***/
    WEAK_SELF(self);
    NSString *rpPkgList = [NSString stringWithUTF8String:RPATH_PACKAGE_LIST];
    
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:rpPkgList
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         
         if ([[response valueForKeyPath:@"package-list.status"] boolValue]) {
             NSArray<Package *> *list = [Package packagesFromList:[response valueForKeyPath:@"package-list.list"]];
             if (list != nil) {
                 [self.packageList addObjectsFromArray:list];
                 [self.packageTable reloadData];
             } else {
                 [ShowAlert
                  showWarningAlertWithTitle:@"Temporarily Unavailable"
                  message:@"Unable to retrieve available packages. Please try a bit later."];
             }
             
         } else {
             [ShowAlert
              showWarningAlertWithTitle:@"Temporarily Unavailable"
              message:[response valueForKeyPath:@"package-list.error"]];
         }

         [self _enableControls];
         [[PCRouter sharedRouter] delGetRequest:belf onPath:rpPkgList];
     }];
    
    [self _disableControls];
    [PCRouter routeRequestGet:RPATH_PACKAGE_LIST];
}

#pragma mark - NSTableViewDataSourceDelegate
- (NSInteger)numberOfRowsInTableView:(NSTableView *)tableView {
    return [self.packageList count];
}

- (nullable id)tableView:(NSTableView *)tableView objectValueForTableColumn:(nullable NSTableColumn *)tableColumn row:(NSInteger)row {
    return [[self.packageList objectAtIndex:row] description];
}

#pragma mark - NSTableViewDelegate
-(NSView *)tableView:(NSTableView *)aTableView viewForTableColumn:(NSTableColumn *)aTableColumn row:(NSInteger)row {
    Package *meta = [self.packageList objectAtIndex:row];
    NSTableCellView *nv = [aTableView makeViewWithIdentifier:@"packageview" owner:self];
    [nv.textField setStringValue:[meta packageDescription]];
    return nv;
}

// disable table row text editing
- (BOOL)tableView:(NSTableView *)tableView
shouldEditTableColumn:(NSTableColumn *)tableColumn
              row:(NSInteger)row {
    return NO;
}

// enable table row selection
- (BOOL)selectionShouldChangeInTableView:(NSTableView *)tableView {
    return YES;
}

- (BOOL)tableView:(NSTableView *)tableView shouldSelectRow:(NSInteger)row {
    return (![self.packageList objectAtIndex:(NSUInteger)row].installed);
}

- (NSIndexSet *)tableView:(NSTableView *)aTableView
selectionIndexesForProposedSelection:(NSIndexSet *)anIndex {
    NSInteger row = (NSInteger)anIndex.firstIndex;
    if (![self.packageList objectAtIndex:(NSUInteger)row].installed) {
        _selectedIndex = row;
    }
    return anIndex;
}

#pragma mark - Setup UI states
-(void)_enableControls {
    [self.btnInstall setEnabled:YES];
    [self.btnCancel setEnabled:YES];
    [self.progressLabel setStringValue:@""];
    [self.circularProgress setHidden:YES];
    [self.circularProgress stopAnimation:nil];
    [self.progressBar displayIfNeeded];
}

-(void)_disableControls {
    [self.btnInstall setEnabled:NO];
    [self.btnCancel setEnabled:NO];
    [self.progressLabel setStringValue:@""];
    [self.circularProgress setHidden:NO];
    [self.circularProgress startAnimation:nil];
    [self.progressBar setDoubleValue:0.0];
    [self.progressBar displayIfNeeded];
}

- (void)_setUIToProceedState {
    [self.btnInstall setEnabled:NO];
    [self.circularProgress startAnimation:nil];
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
    //[self.stageControl shouldControlProgressFrom:self withParam:nil];

    if (_selectedIndex == -1 || (NSInteger)[self.packageList count] <= _selectedIndex ) {
        return;
    }
    
    /*** checking user authed ***/
    WEAK_SELF(self);
    NSString *rpPkgInst = [NSString stringWithUTF8String:RPATH_PACKAGE_INSTALL];
    
    [[PCRouter sharedRouter]
     addPostRequest:self
     onPath:rpPkgInst
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         
         if ([[response valueForKeyPath:@"package-install.status"] boolValue]) {
         } else {
             [ShowAlert
              showWarningAlertWithTitle:@"Temporarily Unavailable"
              message:[response valueForKeyPath:@"package-install.error"]];
         }
         
         [self _enableControls];
         [[PCRouter sharedRouter] delGetRequest:belf onPath:rpPkgInst];
     }];
    
    [self _disableControls];
    [PCRouter
     routeRequestPost:RPATH_PACKAGE_INSTALL
     withRequestBody:@{@"package-id":[self.packageList objectAtIndex:(NSUInteger)_selectedIndex].packageID}];
}

-(IBAction)cancel:(id)sender {
    [self.stageControl shouldControlRevertFrom:self withParam:nil];    
}
@end
