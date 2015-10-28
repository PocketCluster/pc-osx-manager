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


@interface PCSetup2RPVC ()<PCTaskDelegate, GCDAsyncUdpSocketDelegate>
@property (strong, nonatomic) PCTask *sudoTask;
@property (strong, nonatomic) PCTask *userTask;
@property (nonatomic, strong) GCDAsyncUdpSocket *udpSocket;
@property (nonatomic, strong) NSMutableArray<LinkInterface *> *localInterfaces;
@property (atomic, strong) NSMutableArray *nodeList;
@end

@implementation PCSetup2RPVC

-(instancetype)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil {

    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    
    if(self){
        [[Util getApp] addMultDelegateToQueue:self];
        self.localInterfaces = [NSMutableArray arrayWithCapacity:0];
        [self.localInterfaces addObjectsFromArray: [PCInterfaceList all]];
        
        self.nodeList = [NSMutableArray arrayWithCapacity:0];
    }
    
    return self;
}

-(void)dealloc {
    [[Util getApp] removeMultDelegateFromQueue:self];
}


-(IBAction)startBulding:(id)sender
{
    
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
//            Log(@"same mac address found! %@ %@",[node valueForKey:@"pc_sl_ma"], [m valueForKey:@"pc_sl_ma"]);
            doesNodeExist = true;
            break;
        }
    }

    if (!doesNodeExist){
        
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

    

#if 0
    NSString *sn = [[DeviceSerialNumber deviceSerialNumber] lowercaseString];
    
    NSString *hn = [[[NSHost currentHost] localizedName] lowercaseString];
    
    NSString *ha = [[NSHost currentHost] address];
    
    NSMutableDictionary* n = [NSMutableDictionary dictionaryWithDictionary:m];
    [n setValuesForKeysWithDictionary:
     @{@"pc_ma_ct":@"ct_fix_bound",
       @"pc_ma_hn":hn,
       @"pc_ma_ba":sn,
       
#warning fix!
       //@"pc_ma_i4":ha,
       @"pc_ma_i4":@"192.168.1.152",
       @"pc_ma_i6":@""}];
    
    [n setValue:@"rpi-node1" forKey:@"pc_sl_nm"];
    
    if (n)
    {
        NSLog(@"%@",[m description]);
    }
    else
    {
        NSLog(@"Error converting received data into UTF-8 String");
    }
    
    [self.udpSocket sendData:[n BSONRepresentation] toHost:@"239.193.127.127" port:10061 withTimeout:-1 tag:0];
#endif

}

#pragma mark - PCTaskDelegate
-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    
    if(self.sudoTask){
        self.sudoTask = nil;
        
        NSString *basePath  = [[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"];
        NSString *userSetup = [NSString stringWithFormat:@"%@/setup/vagrant_user_setup.sh",basePath];
        
        PCTask *userTask = [PCTask new];
        userTask.taskCommand = [NSString stringWithFormat:@"sh %@ %@", userSetup, basePath];
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

@end
