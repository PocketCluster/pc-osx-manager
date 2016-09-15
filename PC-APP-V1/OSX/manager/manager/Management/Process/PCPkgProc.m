//
//  PCPkgProc.m
//  manager
//
//  Created by Almighty Kim on 11/10/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCPkgProc.h"
#import "TaskOutputWindow.h"
#import "PCConstants.h"
#import "PCProcManager.h"
#import "PCTask.h"
#import "Util.h"

#define USE_TASK_WINDOW

@interface PCPkgProc()<PCTaskDelegate>
@property (nonatomic, weak, readwrite) PCPackageMeta *package;
@property (strong, nonatomic) NSMutableDictionary *procCheckDict;
@property (strong, nonatomic) PCTask *procCheckTask;

@property (strong, nonatomic) PCTask *procStartTask;
@property (weak, nonatomic) TaskOutputWindow *startWindow;

@property (strong, nonatomic) PCTask *procStopTask;
@property (weak, nonatomic) TaskOutputWindow *stopWindow;

- (PCTask *)executionTask:(NSString *)anAction;
@end

@implementation PCPkgProc{
    BOOL _isAlive;
}
@dynamic isAlive;

- (instancetype)initWithPackageMeta:(PCPackageMeta *)aPackageMeta {
    self = [super init];
    if(self){
        _isAlive = NO;
        self.procCheckDict = nil;
        self.package = aPackageMeta;
    }
    return self;
}

#pragma mark - PCTaskDelegate
- (void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {

    if(self.procCheckTask == aPCTask){
        BOOL isAlive = YES;
        for(NSString *pn in self.procCheckDict){
            NSNumber *pls = (NSNumber *)[self.procCheckDict objectForKey:pn];
            isAlive = (isAlive & [pls boolValue]);
        }
        self.procCheckDict = nil;
        @synchronized(self) {
            _isAlive = isAlive;
        }
        
        [[NSNotificationCenter defaultCenter]
         postNotificationName:kPOCKET_CLUSTER_PACKAGE_PROCESS_STATUS
         object:nil
         userInfo:
         @{kPOCKET_CLUSTER_PACKAGE_PROCESS_ISALIVE: @(isAlive),
           kPOCKET_CLUSTER_PACKAGE_IDENTIFICATION:self.package.packageId}];

//Log(@"%s, %@ is %@",__PRETTY_FUNCTION__, self.package.packageName, isAlive?@"ALIVE":@"DEAD");
        
        self.procCheckTask = nil;
    }
    
    if(self.procStartTask == aPCTask){
        [self.startWindow task:aPCTask taskCompletion:aTask], self.startWindow = nil;
        self.procStartTask = nil;
    }

    if(self.procStopTask == aPCTask){
        
        __weak NSString *pid = self.package.packageId;

        // should fire dead notification.
        [[NSNotificationCenter defaultCenter]
         postNotificationName:kPOCKET_CLUSTER_PACKAGE_PROCESS_STATUS
         object:nil
         userInfo:
         @{kPOCKET_CLUSTER_PACKAGE_PROCESS_ISALIVE:@(NO),
           kPOCKET_CLUSTER_PACKAGE_IDENTIFICATION:pid}];
        
        [self.stopWindow task:aPCTask taskCompletion:aTask], self.stopWindow = nil;
        self.procStopTask = nil;

        //FIXME: this need to be included and should be one-stop function call
        [[PCProcManager sharedManager] removePackageProcess:self];

    }
}

- (void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {

    if(self.procCheckTask == aPCTask){

        NSData *data = [aFileHandler availableData];
        NSString *str = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
        
        NSString *pName = self.package.packageName;
        NSString *pVer  = self.package.version;
        NSString *pMode = self.package.modeType;

        NSMutableDictionary *pcl = [NSMutableDictionary dictionaryWithDictionary:self.procCheckDict];
        for (NSString *pkgproc in self.procCheckDict){
            
            //FIXME: Spark cannot make use of package name, version or, mode. What should I do?
            BOOL ppcheck = [str containsString:pkgproc]; //& [str containsString:pName] & [str containsString:pVer] & [str containsString:pMode];
            if(ppcheck){
                [pcl setValue:@(YES) forKey:pkgproc];
            }
        }
        self.procCheckDict = pcl;
    }
    
    if(self.procStartTask == aPCTask) {
        [self.startWindow task:aPCTask recievedOutput:aFileHandler];
    }
    
    if(self.procStopTask == aPCTask) {
        [self.stopWindow task:aPCTask recievedOutput:aFileHandler];
    }
}

- (BOOL)task:(PCTask *)aPCTask isOutputClosed:(id<PCTaskDelegate>)aDelegate {return NO;}

#pragma mark - Alive
- (BOOL)isAlive {
    BOOL alive = NO;
    
    @synchronized(self) {
        alive = _isAlive;
    }

    return alive;
}

#pragma mark - Instance Methods
- (void)refreshProcessStatus {

    NSMutableDictionary *pcd = [NSMutableDictionary dictionary];
    for(NSString *pn in self.package.processCheck){
        [pcd setObject:@(NO) forKey:pn];
    }
    self.procCheckDict = pcd;

    PCTask *pct = [PCTask new];
    pct.taskCommand = @"jps -v";
    pct.delegate = self;
    self.procCheckTask = pct;
    
    [pct launchTask];
}

- (void)startPackageProcess {
#ifdef USE_TASK_WINDOW
    Assert([NSThread isMainThread], @"startPackageProcess should run in Main Thread");
    
    PCTask *spp = [self executionTask:[NSString stringWithFormat:@"bash %@",[self.package.startScript objectAtIndex:0]]];
    TaskOutputWindow *ow = [[TaskOutputWindow alloc] initWithWindowNibName:@"TaskOutputWindow"];
    ow.taskOperator = spp;

    self.procStartTask = spp;
    self.startWindow = ow;
    
    [NSApp activateIgnoringOtherApps:YES];
    [ow showWindow:[Util getApp]];
    
    [[Util getApp] addOpenWindow:ow];
#else
    PCTask *spp = [PCTask new];
    spp.taskCommand = [NSString stringWithFormat:@"bash %@",[self.package.startScript objectAtIndex:0]];
    spp.delegate = self;
    self.procStartTask = spp;
    [spp launchTask];
#endif
}

- (void)stopPackageProcess {
#ifdef USE_TASK_WINDOW
    Assert([NSThread isMainThread], @"stopPackageProcess should run in Main Thread");
    
    PCTask *spp = [self executionTask:[NSString stringWithFormat:@"bash %@",[self.package.stopScript objectAtIndex:0]]];
    TaskOutputWindow *ow = [[TaskOutputWindow alloc] initWithWindowNibName:@"TaskOutputWindow"];
    ow.taskOperator = spp;
    
    self.procStopTask = spp;
    self.stopWindow = ow;
    
    [NSApp activateIgnoringOtherApps:YES];
    [ow showWindow:[Util getApp]];

    [[Util getApp] addOpenWindow:ow];
#else
    PCTask *spp = [PCTask new];
    spp.taskCommand = [NSString stringWithFormat:@"bash %@",[self.package.stopScript objectAtIndex:0]];
    spp.delegate = self;
    self.procStopTask = spp;
    [spp launchTask];
#endif
}

- (PCTask *)executionTask:(NSString *)anAction {

    NSArray *cmdParts = @[anAction,@"--no-color"];
    NSString *command = [cmdParts componentsJoinedByString:@" "];
    
    NSTask *task = [[NSTask alloc] init];
    [task setLaunchPath:@"/bin/bash"];
    [task setArguments:@[@"-c", @"-l", command]];
    
    PCTask *spp = [PCTask new];
    spp.task = task;
    spp.taskCommand = command;
    spp.taskAction = command;
    spp.delegate = self;

    return spp;
}
@end
