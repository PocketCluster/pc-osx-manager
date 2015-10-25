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

@interface PCTask()
- (NSTask *)defaultTask;
- (NSString *)sudoCommand;
@end

@implementation PCTask
-(instancetype)init {
    self = [super init];
    
    if(self){
        CFUUIDRef uuid = CFUUIDCreate(NULL);
        self.taskUUID = (__bridge_transfer NSString *)CFUUIDCreateString(NULL, uuid);
        CFRelease(uuid);
        
        _sudoCommand = NO;
    }

    return self;
}

-(void)dealloc
{
    [[NSNotificationCenter defaultCenter] removeObserver:self];
}

- (NSTask *)defaultTask
{
    NSTask *task = [[NSTask alloc] init];
    [task setLaunchPath:@"/bin/bash"];
    
    if(_sudoCommand){
        [task setArguments:@[@"-l", @"-c", [self sudoCommand]]];
    }else{
        [task setArguments:@[@"-l", @"-c", self.taskCommand]];
    }
    
    return task;
}

- (NSString *)sudoCommand
{
    return [NSString stringWithFormat:@"sudo -A %@", self.taskCommand];
}


-(void)launchTask {
    
    if (!self.task){
        self.task = [self defaultTask];
    }
    
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
    
    @synchronized(self) {
        [self.delegate task:self taskCompletion:task];
    }
}

- (void)receivedOutput:(NSNotification*)notif {
    NSFileHandle *fh = [notif object];
    
    @synchronized(self) {
        if(self.delegate && ![self.delegate task:self isOutputClosed:self.delegate]){
            
            [self.delegate task:self recievedOutput:fh];

            if(self.task.isRunning){
                [fh waitForDataInBackgroundAndNotify];
            }
        }
    }
}

- (void)cancelTask {
    [self.task interrupt],[self.task terminate],self.task = nil;
}

@end
