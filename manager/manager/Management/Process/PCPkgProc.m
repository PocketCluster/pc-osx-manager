//
//  PCPkgProc.m
//  manager
//
//  Created by Almighty Kim on 11/10/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCPkgProc.h"
#import "PCTask.h"

@interface PCPkgProc()<PCTaskDelegate>
@property (strong, nonatomic) PCTask *procCheckTask;
@end

@implementation PCPkgProc

#pragma mark - PCTaskDelegate

-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {
    NSData *data = [aFileHandler availableData];
    NSString *str = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
    
    Log(@"STR %@",str);
    
    NSString *pName = self.package.packageName;
    NSString *pVer  = self.package.version;
    NSString *pMode = self.package.modeType;
    
    BOOL isProcAlive = YES;


    for (NSString *proc in [str componentsSeparatedByString:@"\n"]){
        for (NSString *pkgproc in self.package.processCheck){
            BOOL pkgRunning = [proc containsString:pkgproc] & [proc containsString:pName] & [proc containsString:pVer] & [proc containsString:pMode];
            isProcAlive = (isProcAlive & pkgRunning);
        }
    }

}

-(BOOL)task:(PCTask *)aPCTask isOutputClosed:(id<PCTaskDelegate>)aDelegate {
    return NO;
}


-(void)refreshProcessStatus {

    PCTask *pct = [PCTask new];
    pct.taskCommand = @"";
    pct.delegate = self;
    self.procCheckTask = pct;
    
    [pct launchTask];
    
}
@end
