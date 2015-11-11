//
//  PCPkgProc.m
//  manager
//
//  Created by Almighty Kim on 11/10/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCPkgProc.h"
#import "PCTask.h"
#import "PCConstants.h"

@interface PCPkgProc()<PCTaskDelegate>
@property (nonatomic, weak, readwrite) PCPackageMeta *package;
@property (strong, nonatomic) NSMutableDictionary *procCheckDict;
@property (strong, nonatomic) PCTask *procCheckTask;
@property (strong, nonatomic) PCTask *procStartTask;
@property (strong, nonatomic) PCTask *procStopTask;
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
        
        self.procCheckTask = nil;
    }
    
    if(self.procStartTask == aPCTask){
        self.procStartTask = nil;
    }

    if(self.procStopTask == aPCTask){
        self.procStopTask = nil;
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
            BOOL ppcheck = [str containsString:pkgproc] & [str containsString:pName] & [str containsString:pVer] & [str containsString:pMode];
            if(ppcheck){
                [pcl setValue:@(YES) forKey:pkgproc];
            }
        }
        self.procCheckDict = pcl;
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
    PCTask *spp = [PCTask new];
    spp.taskCommand = [NSString stringWithFormat:@"bash %@",[self.package.startScript objectAtIndex:0]];
    spp.delegate = self;
    self.procStartTask = spp;
    [spp launchTask];
}

- (void)stopPackageProcess {
    PCTask *spp = [PCTask new];
    spp.taskCommand = [NSString stringWithFormat:@"bash %@",[self.package.stopScript objectAtIndex:0]];
    spp.delegate = self;
    self.procStartTask = spp;
    [spp launchTask];
}

@end
