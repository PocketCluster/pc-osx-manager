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
#import "NullStringChecker.h"

@interface PCSetup3VC()<PCRouteRequest>
@property (nonatomic, strong) NSMutableArray<Package *> *packageList;
+ (void)_enableControls:(PCSetup3VC *)vc;
+ (void)_disableControls:(PCSetup3VC *)vc;
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

         [PCSetup3VC _enableControls:belf];
         [[PCRouter sharedRouter] delGetRequest:belf onPath:rpPkgList];
     }];
    
    [PCSetup3VC _disableControls:self];
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
- (BOOL)tableView:(NSTableView *)tableView shouldEditTableColumn:(NSTableColumn *)tableColumn row:(NSInteger)row {
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
    if ([anIndex count] == 0) {
        return anIndex;
    }

    NSInteger row = (NSInteger)anIndex.firstIndex;
    if (![self.packageList objectAtIndex:(NSUInteger)row].installed) {
        _selectedIndex = row;
    }
    return anIndex;
}

#pragma mark - Setup UI states
+ (void)_enableControls:(PCSetup3VC *)vc {
    [vc.btnInstall    setEnabled:YES];
    [vc.btnCancel     setEnabled:YES];
    [vc.progressLabel setStringValue:@""];
    [vc.progressBar   displayIfNeeded];

    [vc.circularProgress setHidden:YES];
    [vc.circularProgress stopAnimation:nil];
    [vc.circularProgress displayIfNeeded];
    [vc.circularProgress removeFromSuperview];
    [vc setCircularProgress:nil];
}

+ (void)_disableControls:(PCSetup3VC *)vc {
    [vc.btnInstall    setEnabled:NO];
    [vc.btnCancel     setEnabled:NO];
    [vc.progressLabel setStringValue:@""];
    [vc.progressBar   displayIfNeeded];

    NSProgressIndicator *ind = [[NSProgressIndicator alloc] initWithFrame:(NSRect){{20.0, 20.0}, {16.0, 16.0}}];
    [ind setStyle:NSProgressIndicatorSpinningStyle];
    [vc.view addSubview:ind];
    [vc setCircularProgress:ind];
    [ind setHidden:NO];
    [ind setIndeterminate:YES];
    [ind startAnimation:vc];
    [ind displayIfNeeded];
}

#pragma mark - IBACTION
-(IBAction)install:(id)sender {
    static const double unit_gigabyte = 1073741824.0;
    static const double unit_megabyte = 1048576.0;

    //[self.stageControl shouldControlProgressFrom:self withParam:nil];

    if (_selectedIndex == -1 || (NSInteger)[self.packageList count] <= _selectedIndex ) {
        return;
    }

    [PCSetup3VC _disableControls:self];
    
    /*** checking user authed ***/
    WEAK_SELF(self);
    NSString *rpPkgInstall = [NSString stringWithUTF8String:RPATH_PACKAGE_INSTALL];
    NSString *rpPkgInstProg = [NSString stringWithUTF8String:RPATH_PACKAGE_INSTALL_PROGRESS];

    [[PCRouter sharedRouter]
     addPostRequest:self
     onPath:rpPkgInstall
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         
         if ([[response valueForKeyPath:@"package-install.status"] boolValue]) {
         } else {
             [ShowAlert
              showWarningAlertWithTitle:@"Temporarily Unavailable"
              message:[response valueForKeyPath:@"package-install.error"]];
         }

         [PCSetup3VC _enableControls:belf];
         [[PCRouter sharedRouter] delPostRequest:belf onPath:rpPkgInstall];
         [[PCRouter sharedRouter] delPostRequest:belf onPath:rpPkgInstProg];
     }];

    [[PCRouter sharedRouter]
     addPostRequest:self
     onPath:rpPkgInstProg
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         NSString *message = [response valueForKeyPath:@"package-progress.message"];
         if (!ISNULL_STRING(message)) {
             [belf.progressLabel setStringValue:message];

         } else {
             double bytes_total = [[response valueForKeyPath:@"package-progress.total-size"] doubleValue];
             NSString *stringTotal;
             if (unit_gigabyte < bytes_total) {
                 bytes_total /= unit_gigabyte;
                 stringTotal = [NSString stringWithFormat:@"%.1lf GB", bytes_total];
             } else {
                 bytes_total /= unit_megabyte;
                 stringTotal = [NSString stringWithFormat:@"%.1lf MB", bytes_total];
             }

             double bytes_received = [[response valueForKeyPath:@"package-progress.received"] doubleValue];
             NSString *stringReceived;
             if (unit_gigabyte < bytes_received) {
                 bytes_received /= unit_gigabyte;
                 stringReceived = [NSString stringWithFormat:@"%.1lf GB", bytes_received];
             } else {
                 bytes_received /= unit_megabyte;
                 stringReceived = [NSString stringWithFormat:@"%.1lf MB", bytes_received];
             }

             NSString *speed = [NSString stringWithFormat:@"Total %@ Received %@ (%.1lf MB/sec)"
                                ,stringTotal, stringReceived
                                ,([[response valueForKeyPath:@"package-progress.speed"] doubleValue] / unit_megabyte)];
             [belf.progressLabel setStringValue:speed];

             [belf.progressBar setDoubleValue:[[response valueForKeyPath:@"package-progress.done-percent"] doubleValue]];
         }
     }];

    [PCRouter
     routeRequestPost:RPATH_PACKAGE_INSTALL
     withRequestBody:@{@"pkg-id":[self.packageList objectAtIndex:(NSUInteger)_selectedIndex].packageID}];
}

-(IBAction)cancel:(id)sender {
    [self.stageControl shouldControlRevertFrom:self withParam:nil];    
}
@end
