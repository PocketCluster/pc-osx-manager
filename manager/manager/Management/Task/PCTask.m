//
//  PCTask.m
//  manager
//
//  Created by Almighty Kim on 10/23/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCTask.h"
#import "Util.h"
#import "AppDelegate.h"


@implementation PCTask
-(instancetype)init {
    self = [super init];
    

    if(self){
        CFUUIDRef uuid = CFUUIDCreate(NULL);
        self.taskUUID = (__bridge_transfer NSString *)CFUUIDCreateString(NULL, uuid);
        CFRelease(uuid);
    }

    return self;
}

-(void)launchTask {
    
    NSPipe *taskOutputPipe = [NSPipe pipe];
    [self.task setStandardInput:[NSFileHandle fileHandleWithNullDevice]];
    [self.task setStandardOutput:taskOutputPipe];
    [self.task setStandardError:taskOutputPipe];
    
    //set up Askpass handler for sudo
    NSString *askPassPath = [NSBundle pathForResource:@"Askpass" ofType:@"" inDirectory:[[NSBundle mainBundle] bundlePath]];
    NSMutableDictionary *env = [[[NSProcessInfo processInfo] environment] mutableCopy];
    [env setObject:@"NONE" forKey:@"DISPLAY"];
    [env setObject:askPassPath forKey:@"SUDO_ASKPASS"];
    [self.task setEnvironment:env];
    
    NSFileHandle *fh = [taskOutputPipe fileHandleForReading];
    [fh waitForDataInBackgroundAndNotify];
    
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(receivedOutput:) name:NSFileHandleDataAvailableNotification object:fh];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(taskCompletion:) name: NSTaskDidTerminateNotification object:self.task];
    
    [[self task] launch];
}

- (void)taskCompletion:(NSNotification*)notif {
    NSTask *task = [notif object];
    
    
    NSString *notificationText;
    
    if(task.terminationStatus != 0) {
        notificationText = @"Task completed with errors";
        
    } else {
        notificationText = @"Task completed successfully";
    }
    
    /*
     [[Util getApp] showUserNotificationWithTitle:notificationText informativeText:[NSString stringWithFormat:@"%@ %@", name, self.taskAction] taskWindowUUID:self.windowUUID];
     */
    
    //notify app task is complete
    [[NSNotificationCenter defaultCenter] postNotificationName:@"vagrant-manager.task-completed" object:nil userInfo:@{@"target": self.target}];

/*
    if([[NSUserDefaults standardUserDefaults] boolForKey:@"autoCloseTaskWindows"] && task.terminationStatus == 0) {
        dispatch_async(dispatch_get_global_queue(0,0), ^{
            [self close];
        });
    }
 */
}

- (void)receivedOutput:(NSNotification*)notif {
    NSFileHandle *fh = [notif object];
    //NSData *data = [fh availableData];
    //NSString *str = [[NSString alloc] initWithData:data encoding:NSASCIIStringEncoding];
    
    @synchronized(self) {
  
        if(self.task.isRunning) {
            [fh waitForDataInBackgroundAndNotify];
        }

    }
}

- (void)cancelTask {
    [self.task interrupt],[self.task terminate],self.task = nil;
}

@end
