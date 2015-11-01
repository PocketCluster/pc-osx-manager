//
//  PCSetup2RPVC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCSetup2RPVC.h"

#import "RaspberryManager.h"
#import "PCInterfaceList.h"
#import "DeviceSerialNumber.h"
#import "BSONSerialization.h"
#import "Util.h"
#import "PCTask.h"
#import "NullStringChecker.h"
#import "RaspberryManager.h"
#import <SystemConfiguration/SCNetworkConfiguration.h>

#define MAX_SUPPORTED_NODE (6)

@interface PCSetup2RPVC ()<PCTaskDelegate, GCDAsyncUdpSocketDelegate>
-(void)refreshInterface;

@property (strong, nonatomic) PCTask *sudoTask;
@property (strong, nonatomic) PCTask *userTask;
@property (nonatomic, strong) GCDAsyncUdpSocket *udpSocket;
//@property (nonatomic, strong) NSMutableArray<LinkInterface *> *localInterfaces;
@property (atomic, strong) NSMutableArray *nodeList;
@property (strong, nonatomic) NSDictionary *progDict;


@property (nonatomic, strong) NSString *deviceSerial;
@property (nonatomic, strong) NSString *hostName;
@property (nonatomic, strong) LinkInterface *interface;

@end

@implementation PCSetup2RPVC

-(instancetype)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil {

    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    
    if(self){
        [[RaspberryManager sharedManager] addMultDelegateToQueue:self];
        self.nodeList = [NSMutableArray arrayWithCapacity:0];
        
        self.progDict = @{@"SUDO_SETUP_STEP_0":@[@"Base config done...",@10.0]
                          ,@"SUDO_SETUP_DONE":@[@"Start setting up Vagrant",@20.0]
                          ,@"USER_SETUP_STEP_0":@[@"USER_SETUP_STEP_0",@30.0]
                          ,@"USER_SETUP_STEP_1":@[@"USER_SETUP_STEP_1",@50.0]
                          ,@"USER_SETUP_STEP_2":@[@"USER_SETUP_STEP_2",@90.0]
                          ,@"USER_SETUP_DONE":@[@"USER_SETUP_DONE",@100.0]};
        
        
        self.deviceSerial = [[DeviceSerialNumber deviceSerialNumber] lowercaseString];
        self.hostName = [[[NSHost currentHost] localizedName] lowercaseString];

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
    [[RaspberryManager sharedManager] removeMultDelegateFromQueue:self];
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
        if([[node valueForKey:SLAVE_NODE_MACADDR] isEqualToString:[m valueForKey:SLAVE_NODE_MACADDR]]){
            doesNodeExist = true;
            break;
        }
    }

    if (!doesNodeExist && self.nodeList.count <= MAX_SUPPORTED_NODE){
        
        NSString *sn = self.deviceSerial;
        NSString *hn = self.hostName;
        //NSString *ha = [[NSHost currentHost] address];
        NSString *ha = [self.interface ip4Address];
        
        NSMutableDictionary* n = [NSMutableDictionary dictionaryWithDictionary:m];
        [n setValuesForKeysWithDictionary:
         @{MASTER_COMMAND_TYPE:@"ct_fix_bound",
           MASTER_HOSTNAME:hn,
           MASTER_BOUND_AGENT:sn,
           MASTER_IP4_ADDRESS:ha,
           MASTER_IP6_ADDRESS:@""}];
        
        [self.nodeList addObject:n];
        [self.nodeList sortUsingComparator:^NSComparisonResult(NSDictionary*  _Nonnull node1, NSDictionary* _Nonnull node2) {
            return [[node1 valueForKey:ADDRESS] compare:[node2 valueForKey:ADDRESS] options:NSNumericSearch];
        }];
        
        for (int i = 0; i < [self.nodeList count]; i++){
            NSMutableDictionary *nd = [self.nodeList objectAtIndex:i];
            [nd setValue:[NSString stringWithFormat:@"pc-node%d",(i + 1)] forKey:SLAVE_NODE_NAME];
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
        /*
         [[Util getApp] startSalt];
         sleep(4);
         */
        
        NSString *basePath  = [[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"];
        NSString *userSetup = [NSString stringWithFormat:@"%@/setup/raspberry_user_setup.sh",basePath];
        
        PCTask *userTask = [PCTask new];
        userTask.taskCommand = [NSString stringWithFormat:@"sh %@ %@ %ld", userSetup, basePath, [self.nodeList count]];
        userTask.delegate = self;
        
        self.userTask = userTask;
        [userTask launchTask];
        
        self.sudoTask = nil;
    }else{
        self.userTask = nil;
        [self.progressBar stopAnimation:self];
    }
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {
    
    NSData *data = [aFileHandler availableData];
    __block NSString *str = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
    
    Log(@"STR %@",str);
    
    NSArray *p = nil;
    for (NSString *key in self.progDict) {
        if ([str containsString:key]){
            p = [self.progDict valueForKey:key];
            break;
        }
    }

    if(p != nil){
        [self.progressLabel setStringValue:[p objectAtIndex:0]];
        [self.progressBar setDoubleValue:[[p objectAtIndex:1] doubleValue]];
        [self.progressBar displayIfNeeded];
    }
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

    // save to local configuration
    for (NSUInteger i = 0; i < nodeCount; ++i){
        NSDictionary *node = [self.nodeList objectAtIndex:i];
        [[RaspberryManager sharedManager] addRaspberry:[[Raspberry alloc] initWithDictionary:node]];
        //[[Util getApp] multicastData:[node BSONRepresentation]];
        sleep(1);
    }

    [[RaspberryManager sharedManager] saveRaspberries];
    
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
    
    [self.progressBar startAnimation:self];
    [self.buildBtn setEnabled:NO];
}

@end
