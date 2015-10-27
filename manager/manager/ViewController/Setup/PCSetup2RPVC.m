//
//  PCSetup2RPVC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCSetup2RPVC.h"
#import "Util.h"
#import "PCTask.h"
#import "DeviceSerialNumber.h"

#import "GCDAsyncUdpSocket.h"

@interface PCSetup2RPVC ()<PCTaskDelegate>
@property (strong, nonatomic) PCTask *sudoTask;
@property (strong, nonatomic) PCTask *userTask;
@property (nonatomic, strong) GCDAsyncUdpSocket *udpSocket;
@end

@interface PCSetup2RPVC ()

@end

@implementation PCSetup2RPVC

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do view setup here.
}

-(IBAction)startBulding:(id)sender
{
    
}

-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    
    Log(@"%s",__PRETTY_FUNCTION__);
    
    if(self.sudoTask){
        self.sudoTask = nil;
        
        NSString *basePath    = [[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"];
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
    return false;
}





@end
