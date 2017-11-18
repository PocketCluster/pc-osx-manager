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
#import "StatusCache.h"

static NSString * const kPkgColTag  = @"pkgCol";
static NSString * const kSizeColTag = @"sizeCol";

@interface PCSetup3VC()<PCRouteRequest>
- (void)_enableControls;
- (void)_disableControls;
@end

@implementation PCSetup3VC {
    NSInteger _selectedIndex;
    BOOL      _isInstalling;
}

- (void) finishConstruction {
    [super finishConstruction];
    [self setTitle:@"Install Package"];

    _isInstalling = NO;
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

    // TODO : checking user authed

    [self _disableControls];
    [PCRouter routeRequestGet:RPATH_PACKAGE_LIST_AVAILABLE];
}

#pragma mark - NSWindowDelegate
- (BOOL)windowShouldClose:(NSWindow *)sender {
    if (_isInstalling) {
        [ShowAlert
         showWarningAlertWithTitle:@"Please wait until the install finishes"
         message:@"The installation takes some time. We'll let you know as soon as it's done."];
        return NO;
    }

    return YES;
}

#pragma mark - NSTableViewDataSourceDelegate
- (NSInteger)numberOfRowsInTableView:(NSTableView *)tableView {
    return [[[StatusCache SharedStatusCache] packageList] count];
}

- (nullable id)tableView:(NSTableView *)tableView objectValueForTableColumn:(nullable NSTableColumn *)aTableColumn row:(NSInteger)row {
    if (aTableColumn == nil) {
        return nil;
    }
    Package *pkg = [[[StatusCache SharedStatusCache] packageList] objectAtIndex:row];
    if ([[aTableColumn identifier] isEqualToString:kPkgColTag]) {
        return pkg.packageDescription;

    } else if ([[aTableColumn identifier] isEqualToString:kSizeColTag]) {
        return pkg.totalImageSize;
    }
    return nil;
}

#pragma mark - NSTableViewDelegate
// disable table row text editing
- (BOOL)tableView:(NSTableView *)tableView shouldEditTableColumn:(NSTableColumn *)tableColumn row:(NSInteger)row {
    return NO;
}

// enable table row selection
- (BOOL)selectionShouldChangeInTableView:(NSTableView *)tableView {
    return YES;
}

- (BOOL)tableView:(NSTableView *)tableView shouldSelectRow:(NSInteger)row {
    return (![[[StatusCache SharedStatusCache] packageList] objectAtIndex:(NSUInteger)row].installed);
}

- (NSIndexSet *)tableView:(NSTableView *)aTableView
selectionIndexesForProposedSelection:(NSIndexSet *)anIndex {
    if ([anIndex count] == 0) {
        return anIndex;
    }

    NSInteger row = (NSInteger)anIndex.firstIndex;
    if (![[[StatusCache SharedStatusCache] packageList] objectAtIndex:(NSUInteger)row].installed) {
        _selectedIndex = row;
    }
    return anIndex;
}

#pragma mark - Monitoring Package
// this show all the available package from api backend
- (void) onAvailableListUpdateWith:(StatusCache *)aCache success:(BOOL)isSuccess error:(NSString *)anErrMsg {

    NSArray<Package *> *list = [aCache packageList];
    if (isSuccess && list != nil && [list count]) {
        [self.packageTable reloadData];

    } else {
        [ShowAlert
         showWarningAlertWithTitle:@"Temporarily Unavailable"
         message:anErrMsg];
    }

    [self _enableControls];
}

// this show all the installed package in the system
- (void) onInstalledListUpdateWith:(StatusCache *)aCache success:(BOOL)isSuccess error:(NSString *)anErrMsg {
    NSArray<Package *> *list = [aCache packageList];
    if (isSuccess && list != nil && [list count]) {
        [self.packageTable reloadData];

    } else {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to confirm installed package"
         message:anErrMsg];
    }

    [self _enableControls];
}

#pragma mark - Setup UI states
- (void)_enableControls {
    [self.btnInstall    setEnabled:YES];
    [self.btnCancel     setEnabled:YES];
    [self.progressLabel setStringValue:@""];
    [self.progressBar   displayIfNeeded];

    [self.circularProgress setHidden:YES];
    [self.circularProgress stopAnimation:nil];
    [self.circularProgress displayIfNeeded];
    [self.circularProgress removeFromSuperview];
    [self setCircularProgress:nil];

    _isInstalling = NO;
}

- (void)_disableControls {
    [self.btnInstall    setEnabled:NO];
    [self.btnCancel     setEnabled:NO];
    [self.progressLabel setStringValue:@""];
    [self.progressBar   displayIfNeeded];

    NSProgressIndicator *ind = [[NSProgressIndicator alloc] initWithFrame:(NSRect){{20.0, 20.0}, {16.0, 16.0}}];
    [ind setControlSize:NSSmallControlSize];
    [ind setStyle:NSProgressIndicatorSpinningStyle];
    [self.view addSubview:ind];
    [ind setHidden:NO];
    [ind setIndeterminate:YES];
    [ind startAnimation:self];
    [ind displayIfNeeded];

    [self setCircularProgress:ind];

    _isInstalling = YES;
}

#pragma mark - IBACTION
-(IBAction)install:(id)sender {
    static const double unit_gigabyte = 1073741824.0;
    static const double unit_megabyte = 1048576.0;

    NSString *rpPkgInstall = [NSString stringWithUTF8String:RPATH_PACKAGE_INSTALL];
    NSString *rpPkgInstProg = [NSString stringWithUTF8String:RPATH_PACKAGE_INSTALL_PROGRESS];
    WEAK_SELF(self);

    //[self.stageControl shouldControlProgressFrom:self withParam:nil];

    NSArray<Package *>* pkgList = [[StatusCache SharedStatusCache] packageList];
    if (_selectedIndex == -1 || (NSInteger)[pkgList count] <= _selectedIndex ) {
        return;
    }
    Package *targetPkg = [pkgList objectAtIndex:_selectedIndex];
    if ([targetPkg installed]) {
        [ShowAlert
         showWarningAlertWithTitle:@"Installed Package"
         message:[targetPkg packageDescription]];
        return;
    }

    /// TODO checking if user is authed
    
    [self _disableControls];

    [[PCRouter sharedRouter]
     addPostRequest:self
     onPath:rpPkgInstall
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         if ([[response valueForKeyPath:@"package-install.status"] boolValue]) {
             [ShowAlert
              showWarningAlertWithTitle:@"Installation Completed!"
              message:[targetPkg packageDescription]];

         } else {
             [ShowAlert
              showWarningAlertWithTitle:@"Temporarily Unavailable"
              message:[response valueForKeyPath:@"package-install.error"]];
         }

         [[PCRouter sharedRouter] delPostRequest:belf onPath:rpPkgInstall];
         [[PCRouter sharedRouter] delPostRequest:belf onPath:rpPkgInstProg];
         // ask installed package status to update this and other UI parts
         [PCRouter routeRequestGet:RPATH_PACKAGE_LIST_INSTALLED];
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
     withRequestBody:@{@"pkg-id":[[[StatusCache SharedStatusCache] packageList] objectAtIndex:(NSUInteger)_selectedIndex].packageID}];
}

-(IBAction)cancel:(id)sender {
    [self.stageControl shouldControlRevertFrom:self withParam:nil];    
}
@end
