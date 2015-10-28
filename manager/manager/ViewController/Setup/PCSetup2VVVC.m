//
//  PCSetup2VVVC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCSetup2VVVC.h"

#import "Util.h"
#import "PCTask.h"

//#import "PCTaskOutputWindow.h"

@interface PCSetup2VVVC ()<PCTaskDelegate>
@property (strong, nonatomic) PCTask *sudoTask;
@property (strong, nonatomic) PCTask *userTask;
@end

@implementation PCSetup2VVVC


#pragma mark - PCTaskDelegate

-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    
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

#pragma mark - IBACTION
-(IBAction)build:(id)sender
{
    NSString *basePath    = [[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"];
    NSString *sudoSetup = [NSString stringWithFormat:@"%@/setup/vagrant_sudo_setup.sh",basePath];
    
    PCTask *sudoTask = [PCTask new];
    sudoTask.taskCommand = [NSString stringWithFormat:@"sh %@ %@", sudoSetup, basePath];
    sudoTask.sudoCommand = YES;
    sudoTask.delegate = self;
    
    self.sudoTask = sudoTask;
    
    [sudoTask launchTask];
}

@end
