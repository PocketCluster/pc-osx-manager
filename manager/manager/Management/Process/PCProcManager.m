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

@interface PCProcManager()<PCTaskDelegate>
@property (nonatomic, strong) PCTask *saltMinion;
@property (nonatomic, strong) PCTask *saltMaster;
@property (nonatomic, strong) PCTask *saltClear;
@end


@implementation PCProcManager
SYNTHESIZE_SINGLETON_FOR_CLASS_WITH_ACCESSOR(PCProcManager, sharedManager);


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

- (void)freshStart {
    PCTask *t = [PCTask new];
    t.taskCommand = @"ps -efw | grep salt | grep -v grep | awk '{print $2}' | xargs kill";
    t.delegate = self;
    self.saltClear = t;
    [t launchTask];
}

#pragma mark - PCTaskDelegate
-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    
    if (aPCTask == self.saltClear){
        Log(@"all salt processes killed");
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


@end
