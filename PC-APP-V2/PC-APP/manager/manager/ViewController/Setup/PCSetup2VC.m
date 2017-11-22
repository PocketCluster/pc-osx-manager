//
//  PCSetup2VC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "ShowAlert.h"
#import "PCRouter.h"
#import "PCSetup2VC.h"
#import "PCConstants.h"

@interface PCSetup2VC ()<PCRouteRequest>
@property (nonatomic, strong) NSArray *nodeList;
@end

@implementation PCSetup2VC

- (void) finishConstruction {
    [super finishConstruction];
    [self setTitle:@"Build Cluster"];
    
    WEAK_SELF(self);

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_NODE_REG_START)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         Log(@"path %@ response %@", path, response);
         if (![[response valueForKeyPath:@"node-reg-start.status"] boolValue]) {
             [ShowAlert
              showTerminationAlertWithTitle:@"Unable to add new node"
              message:[response valueForKeyPath:@"node-reg-start.error"]];
         }
     }];

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_NODE_UNREG_LIST)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         Log(@"path %@ response %@", path, response);

         NSArray<NSDictionary *>* list = [response valueForKeyPath:@""];
         if (belf != nil) {
             [belf setNodeList:list];
             [[belf nodeTable] reloadData];
         }
     }];

    [[PCRouter sharedRouter]
     addPostRequest:self
     onPath:@(RPATH_NODE_REG_CANDIDATE)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         Log(@"path %@ response %@", path, response);
         if (belf != nil) {
         }
     }];

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_NODE_REG_CONFIRM)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         Log(@"path %@ response %@", path, response);
         if (belf != nil) {
         }
     }];

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_NODE_REG_STOP)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         Log(@"path %@ response %@", path, response);
         if (belf != nil) {
         }
     }];
}

- (void) prepareDestruction {
    [super prepareDestruction];

    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_NODE_REG_START)];
    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_NODE_UNREG_LIST)];
    [[PCRouter sharedRouter] delPostRequest:self onPath:@(RPATH_NODE_REG_CANDIDATE)];
    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_NODE_REG_CONFIRM)];
    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_NODE_REG_STOP)];
}

- (void) viewDidLoad {
    [super viewDidLoad];
    
    [[((BaseBrandView *)self.view) contentBox] addSubview:self.pannel];
    self.pannel = nil;
}

#pragma mark - NSTableViewDataSourceDelegate
- (NSInteger)numberOfRowsInTableView:(NSTableView *)tableView {
    return [self.nodeList count];
}

- (nullable id)tableView:(NSTableView *)tableView objectValueForTableColumn:(nullable NSTableColumn *)tableColumn row:(NSInteger)row {
    return [self.nodeList objectAtIndex:row];
}

#pragma mark - NSTableViewDelegate
-(NSView *)tableView:(NSTableView *)aTableView viewForTableColumn:(NSTableColumn *)aTableColumn row:(NSInteger)row{
    NSDictionary *nd = [self.nodeList objectAtIndex:row];
    NSTableCellView *nv = [aTableView makeViewWithIdentifier:@"nodeview" owner:self];
    if([aTableColumn.identifier isEqualToString:@"nodename"]){
        [nv.textField setStringValue:[nd valueForKey:@"name"]];
    }else{
        [nv.textField setStringValue:[nd valueForKey:@"addr"]];
    }
    return nv;
}

- (BOOL)selectionShouldChangeInTableView:(NSTableView *)tableView {
    return NO;
}

- (BOOL)tableView:(NSTableView *)tableView shouldSelectRow:(NSInteger)row {
    return NO;
}

#pragma mark - IBACTION
-(IBAction)build:(id)sender {
    [PCRouter routeRequestGet:RPATH_NODE_REG_START];
    return;
    
    [self.stageControl shouldControlProgressFrom:self withParam:nil];

    // return if there is no node
    if ([self.nodeList count] == 0){
        // NSAlert
        return;
    }
}

-(IBAction)cancel:(id)sender {
    [self.stageControl shouldControlRevertFrom:self withParam:nil];
}

@end
