//
//  PCProcManager.m
//  manager
//
//  Created by Almighty Kim on 11/2/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "SynthesizeSingleton.h"
#import "PCProcManager.h"
#import "PCTask.h"
#import "PCConstants.h"
#import <GCDWebServers/GCDWebServers.h>

//#define THREADED_SERVER 1

@interface PCProcManager()<PCTaskDelegate, GCDWebServerDelegate>
@property (nonatomic, strong) GCDWebServer *webServer;
@property (nonatomic, strong) PCTask *saltMinion;
@property (nonatomic, strong) PCTask *saltMaster;
@property (nonatomic, strong) PCTask *saltClear;


@property (nonatomic, strong) NSMutableArray<PCPkgProc *> *packageProcesses;
@property (nonatomic, strong) NSTimer *refreshTimer;

-(void)_webServerStart;
-(void)_webServerStop;

- (void)addPackageProcess:(PCPkgProc *)aPackageProcess;
- (void)removePackageProcess:(PCPkgProc *)aPackageProcess;
@end


@implementation PCProcManager{
    volatile bool _isWebServerRunning;
}
SYNTHESIZE_SINGLETON_FOR_CLASS_WITH_ACCESSOR(PCProcManager, sharedManager);
- (instancetype)init {
    
    self = [super init];
    if (self){
        self.packageProcesses = [NSMutableArray<PCPkgProc *> arrayWithCapacity:0];

        GCDWebServer *ws = [[GCDWebServer alloc] init];
        [ws addGETHandlerForBasePath:@"/" directoryPath:WEB_SERVER_ROOT_PATH indexFilename:nil cacheAge:0 allowRangeRequests:YES];
        [ws setDelegate:self];
        self.webServer = ws;
        _isWebServerRunning = false;
    }
    return self;
}

#pragma mark - SALT MANAGEMENT
- (void)startSalt {
    if(!self.saltMinion){
        PCTask *minion = [[PCTask alloc] init];
        minion.taskCommand = @"salt-minion";
        
        //TODO: enabling delegate take 130% of CPU due to crazy # of invokation to NSNotificationCenter - ?
        //minion.delegate = self;
        self.saltMinion = minion;
        [minion launchTask];
    }
    
    if(!self.saltMaster){
        PCTask *master = [[PCTask alloc] init];
        master.taskCommand = @"salt-master";

        //TODO: enabling delegate take 130% of CPU due to crazy # of invokation to NSNotificationCenter - ?
        //master.delegate = self;
        self.saltMaster = master;
        [master launchTask];
    }
}

- (void)stopSalt {
    if(self.saltMinion){
        [self.saltMinion cancelTask];
        self.saltMinion = nil;
    }

    if (self.saltMaster){
        [self.saltMaster cancelTask];
        self.saltMaster = nil;
    }
}

- (void)freshSaltStart {

#if 0
    if(self.saltMinion != nil && self.saltMaster != nil){
        return;
    }
#else
    [self stopSalt];
#endif
    
    PCTask *t = [PCTask new];
    t.taskCommand = @"ps -efw | grep salt | grep -v grep | awk '{print $2}' | xargs kill";
    t.delegate = self;
    self.saltClear = t;
    [t launchTask];
}

#pragma mark - PCTaskDelegate
-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    if (aPCTask == self.saltClear){
        [self.saltClear cancelTask];
        self.saltClear = nil;
        [self startSalt];
    }
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {
    
}

-(BOOL)task:(PCTask *)aPCTask isOutputClosed:(id<PCTaskDelegate>)aDelegate {
    return NO;
}


#pragma mark - WebServer Control

-(void)_webServerStart {
    _isWebServerRunning = true;
    @autoreleasepool {
        NSDictionary *options =
        @{GCDWebServerOption_Port:@(WEB_SERVER_PORT)
          ,GCDWebServerOption_RequestNATPortMapping:@(NO)
          ,GCDWebServerOption_BindToLocalhost:@(NO)};
        NSError *error;
        [self.webServer startWithOptions:options error:&error];
    }
}

-(void)_webServerStop {
    @autoreleasepool {
        [self.webServer stop];
    }
    _isWebServerRunning = false;
}

-(void)startWebServer {
    
    if(_isWebServerRunning){
        return;
    }

#ifdef THREADED_SERVER
    [self performSelectorInBackground:@selector(_webServerStart) withObject:nil];
#else
    [self _webServerStart];
#endif
}

-(void)stopWebServer {

    if (!_isWebServerRunning){
        return;
    }
#ifdef THREADED_SERVER
    [self performSelectorInBackground:@selector(_webServerStop) withObject:nil];
#else
    [self _webServerStop];
#endif
}

#pragma mark - WebServer Delegates
- (void)webServerDidStart:(GCDWebServer*)server {
}

- (void)webServerDidCompleteBonjourRegistration:(GCDWebServer*)server {
}

- (void)webServerDidUpdateNATPortMapping:(GCDWebServer*)server {
}

- (void)webServerDidConnect:(GCDWebServer*)server {
}

- (void)webServerDidDisconnect:(GCDWebServer*)server {
}

- (void)webServerDidStop:(GCDWebServer*)server {
}


#pragma mark - Package Process Manage
- (PCPkgProc *)startPackageProcess:(PCPackageMeta *)aMetaPacakge {
    
    PCPkgProc *proc = [self findPackageProcess:aMetaPacakge];
    if(proc != nil){
        return proc;
    }

    proc = [[PCPkgProc alloc] initWithPackageMeta:aMetaPacakge];
    [self addPackageProcess:proc];
    [proc startPackageProcess];
    return proc;
}

- (void)stopPackageProcess:(PCPackageMeta *)aMetaPacakge {

    PCPkgProc *proc = [self findPackageProcess:aMetaPacakge];
    if(proc != nil){
        [proc stopPackageProcess];
        [self removePackageProcess:proc];
    }
}

- (void)addPackageProcess:(PCPkgProc *)aPackageProcess {
    
    if(aPackageProcess == nil){
        return;
    }
    
    @synchronized(self.packageProcesses) {
        [self.packageProcesses addObject:aPackageProcess];
    }
}

- (void)removePackageProcess:(PCPkgProc *)aPackageProcess {
    
    if(aPackageProcess == nil){
        return;
    }
    
    @synchronized(self.packageProcesses) {
        [self.packageProcesses removeObject:aPackageProcess];
    }
}

- (PCPkgProc *)findPackageProcess:(PCPackageMeta *)aMetaPackage {
    
    if(aMetaPackage == nil){
        return nil;
    }

    PCPkgProc *proc = nil;
    @synchronized(self.packageProcesses) {
        for(PCPkgProc *pkg in self.packageProcesses) {
            
            // there is only one meta package object at any given time
            if ( pkg.package == aMetaPackage ) {
                proc = pkg;
                break;
            }
        }
    }
    
    return proc;
}

- (void)refreshPackageProcessesStatus {
    @synchronized(self.packageProcesses) {
        [self.packageProcesses makeObjectsPerformSelector:@selector(refreshProcessStatus)];
    }
}

- (void)haltPackageProcessRefreshTimer {
    if (self.refreshTimer) {
        [self.refreshTimer invalidate];
        self.refreshTimer = nil;
    }
}

- (void)startPackageProcessUpdateTimer {
    
    [self haltPackageProcessRefreshTimer];

    self.refreshTimer =
        [NSTimer
         scheduledTimerWithTimeInterval:PROCESS_REFRESH_TIME_INTERVAL
         target:self
         selector:@selector(refreshPackageProcessesStatus)
         userInfo:nil
         repeats:YES];
}

@end
