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


-(IBAction)vagrantUp:(id)sender
{

    NSString *basePath    = [[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"];
    NSString *sudoSetup = [NSString stringWithFormat:@"%@/setup/vagrant_sudo_setup.sh",basePath];
    NSString *userSetup = [NSString stringWithFormat:@"%@/setup/vagrant_user_setup.sh",basePath];

#if 0
    __block PCTaskOutputWindow *tow = [[PCTaskOutputWindow alloc] initWithNibName:nil bundle:nil];
    tow.taskCommand = [NSString stringWithFormat:@"sh %@ %@", setupScript, basePath];
    tow.sudoCommand = YES;
    
    WEAK_SELF(self);
    [belf.view.window
     beginSheet:tow
     completionHandler:^(NSModalResponse returnCode) {
         [belf.view.window endSheet:tow];
     }];
    
    [tow launchTask];
    
    //[[Util getApp] addOpenWindow:tow];
    
#elif 1
    
    PCTask *sudoTask = [PCTask new];
    sudoTask.taskCommand = [NSString stringWithFormat:@"sh %@ %@", sudoSetup, basePath];
    sudoTask.sudoCommand = YES;
    sudoTask.delegate = self;
    
    self.sudoTask = sudoTask;
    
    [sudoTask launchTask];
#else
    
    
    PCTask *userTask = [PCTask new];
    userTask.taskCommand = [NSString stringWithFormat:@"sh %@ %@", userSetup, basePath];
    userTask.delegate = self;
    
    [self.taskQueue addObject:userTask];
    
    [userTask launchTask];

#endif
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
