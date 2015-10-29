//
//  PCSetup2RPVC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCSetup2RPVC.h"

#import "PCInterfaceList.h"
#import "DeviceSerialNumber.h"
#import "GCDAsyncUdpSocket.h"
#import "BSONSerialization.h"
#import "Util.h"
#import "PCTask.h"
#import "NullStringChecker.h"
#import <SystemConfiguration/SCNetworkConfiguration.h>


#define MAX_SUPPORTED_NODE (6)

@interface PCSetup2RPVC ()<PCTaskDelegate, GCDAsyncUdpSocketDelegate>
-(void)refreshInterface;

@property (strong, nonatomic) PCTask *sudoTask;
@property (strong, nonatomic) PCTask *userTask;
@property (nonatomic, strong) GCDAsyncUdpSocket *udpSocket;
//@property (nonatomic, strong) NSMutableArray<LinkInterface *> *localInterfaces;
@property (atomic, strong) NSMutableArray *nodeList;
@property (nonatomic, strong) LinkInterface *interface;
@end

@implementation PCSetup2RPVC

-(instancetype)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil {

    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    
    if(self){
        [[Util getApp] addMultDelegateToQueue:self];
        self.nodeList = [NSMutableArray arrayWithCapacity:0];
        
        [self refreshInterface];
    }
    
    return self;
}

-(void)viewDidAppear {
    if (self.interface){
        [self.warningLabel setHidden:YES];
    }else{
        [self.warningLabel setHidden:NO];
    }
}


-(void)dealloc {
    [[Util getApp] removeMultDelegateFromQueue:self];
}


-(void)refreshInterface {
#if 0
    self.localInterfaces = [NSMutableArray arrayWithCapacity:0];
    for (LinkInterface *iface in [PCInterfaceList all]){
        if ([iface.kind isEqualToString:(__bridge NSString *)kSCNetworkInterfaceTypeEthernet]){
            Log(@"%@-%@ %@", [iface BSDName],[iface ip4Address],[iface kind]);
        }
    }
#endif
    
    for (LinkInterface *iface in [PCInterfaceList all]){
        if (!ISNULL_STRING(iface.ip4Address) && [iface.kind isEqualToString:(__bridge NSString *)kSCNetworkInterfaceTypeEthernet]){
            self.interface = iface;
            return;
        }
    }
}

#pragma mark - GCDAsyncUdpSocketDelegate
- (void)udpSocket:(GCDAsyncUdpSocket *)sock
   didReceiveData:(NSData *)data
      fromAddress:(NSData *)address
withFilterContext:(id)filterContext
{
    NSDictionary *m =[NSDictionary dictionaryWithBSON:data];
    
    BOOL doesNodeExist = false;
    for (NSDictionary *node in self.nodeList){
        if([[node valueForKey:@"pc_sl_ma"] isEqualToString:[m valueForKey:@"pc_sl_ma"]]){
            doesNodeExist = true;
            break;
        }
    }

    if (!doesNodeExist && self.nodeList.count <= 6){
        
        NSString *sn = [[DeviceSerialNumber deviceSerialNumber] lowercaseString];
        NSString *hn = [[[NSHost currentHost] localizedName] lowercaseString];
        NSString *ha = [[NSHost currentHost] address];

        NSMutableDictionary* n = [NSMutableDictionary dictionaryWithDictionary:m];
        [n setValuesForKeysWithDictionary:
         @{@"pc_ma_ct":@"ct_fix_bound",
           @"pc_ma_hn":hn,
           @"pc_ma_ba":sn,
           //@"pc_ma_i4":ha,
           @"pc_ma_i4":@"192.168.1.152",
           @"pc_ma_i6":@""}];

        [self.nodeList addObject:n];
        [self.nodeList sortUsingComparator:^NSComparisonResult(NSDictionary*  _Nonnull node1, NSDictionary* _Nonnull node2) {
            return [[node1 valueForKey:@"address"] compare:[node2 valueForKey:@"address"] options:NSNumericSearch];
        }];
        
        for (int i = 0; i < [self.nodeList count]; i++){
            NSMutableDictionary *nd = [self.nodeList objectAtIndex:i];
            [nd setValue:[NSString stringWithFormat:@"pc-node%d",(i + 1)] forKey:@"pc_sl_nm"];
        }
        
        [self.nodeTable reloadData];
    }
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
        [nv.textField setStringValue:[nd valueForKey:@"pc_sl_nm"]];
    }else{
        [nv.textField setStringValue:[nd valueForKey:@"address"]];
    }
    
    return nv;
}

- (BOOL)selectionShouldChangeInTableView:(NSTableView *)tableView {
    return NO;
}

- (BOOL)tableView:(NSTableView *)tableView shouldSelectRow:(NSInteger)row {
    return NO;
}

#pragma mark - PCTaskDelegate
-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    
    if(self.sudoTask){
        self.sudoTask = nil;
        
        NSString *basePath  = [[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"];
        NSString *userSetup = [NSString stringWithFormat:@"%@/setup/raspberry_user_setup.sh",basePath];
        
        PCTask *userTask = [PCTask new];
        userTask.taskCommand = [NSString stringWithFormat:@"sh %@ %@ %ld", userSetup, basePath, [self.nodeList count]];
        userTask.delegate = self;
        
        self.userTask = userTask;
        [userTask launchTask];
    }else{
        self.userTask = nil;
    }
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {
    
    NSData *data = [aFileHandler availableData];
    NSString *str = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
    
    Log(@"%@",str);
}

-(BOOL)task:(PCTask *)aPCTask isOutputClosed:(id<PCTaskDelegate>)aDelegate {
    return NO;
}


#pragma mark - IBACTION
-(IBAction)build:(id)sender
{
    // update interface status
    [self refreshInterface];
    
    // if there is no Ethernet, do not proceed.
    if (self.interface){
        [self.warningLabel setHidden:YES];
    }else{
        [self.warningLabel setHidden:NO];
        return;
    }

    // return if there is no node
    NSUInteger nodeCount = MIN([self.nodeList count], MAX_SUPPORTED_NODE);
    if (nodeCount <= 0){
        // NSAlert
        return;
    }

    // setup only six nodes
#if 0
    // save to local configuration
    for (NSUInteger i = 0; i < nodeCount; ++i){
        NSDictionary *node = [self.nodeList objectAtIndex:i];
        [[Util getApp] multicastData:[node BSONRepresentation]];
        sleep(1);
    }
#endif

    NSMutableString *nodeip = [NSMutableString new];
    for (NSUInteger i = 0; i < nodeCount; ++i){
        NSDictionary *node = [self.nodeList objectAtIndex:i];
        [nodeip appendString:[NSString stringWithFormat:@"%@ ", [node valueForKey:@"address"]]];
    }
    
    NSString *basePath = [[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"];
    NSString *sudoSetup = [NSString stringWithFormat:@"%@/setup/raspberry_sudo_setup.sh %@ %@ %@",basePath, basePath, self.interface.ip4Address, nodeip];
    
    PCTask *sudoTask = [PCTask new];
    sudoTask.taskCommand = [NSString stringWithFormat:@"sh %@",sudoSetup];
    sudoTask.sudoCommand = YES;
    sudoTask.delegate = self;
    
    self.sudoTask = sudoTask;
    
    [sudoTask launchTask];
}

@end
