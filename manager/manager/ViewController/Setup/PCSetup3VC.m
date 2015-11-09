//
//  PCSetup3VC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCSetup3VC.h"
#import "PCPackageMeta.h"
#import "PCConstants.h"
#import "Util.h"
#import "PCTask.h"
#import "RaspberryManager.h"

@interface PCSetup3VC()<PCTaskDelegate>
@property (nonatomic, strong) NSMutableArray<PCPackageMeta *> *packageList;
@property (nonatomic, strong) NSMutableArray<NSString *> *downloadFileList;

@property (nonatomic, strong) PCTask *javaTask;
@property (nonatomic, strong) PCTask *saltBaseTask;
@property (nonatomic, strong) PCTask *saltMasterTask;
@property (nonatomic, strong) PCTask *saltMinionTask;
@property (nonatomic, strong) PCTask *saltMasterCompleteTask;
@property (readwrite, nonatomic) BOOL canContinue;
@property (readwrite, nonatomic) BOOL canGoBack;

- (void)setUIToProceedState;
- (void)resetUIForFailure;
- (void)setToNextStage;
- (void)setProgMessage:(NSString *)aMessage value:(double)aValue;

- (void)startPackageInstall;
- (void)failedPackageInstall;
- (void)finalizeInstallProcess;
- (void)downloadMetaFiles;
@end

@implementation PCSetup3VC

@synthesize canContinue;
@synthesize canGoBack;

-(instancetype)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil {
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if(self){
        self.packageList = [NSMutableArray arrayWithCapacity:0];
        self.downloadFileList = [NSMutableArray arrayWithCapacity:0];
        
        [self resetToInitialState];
        
        WEAK_SELF(self);
        [PCPackageMeta metaPackageListWithBlock:^(NSArray<PCPackageMeta *> *packages, NSError *error) {
            if(belf != nil){
                [belf.packageList addObjectsFromArray:packages];
                [belf.packageTable reloadData];
            }
        }];
    }
    return self;
}

#pragma mark - NSTableViewDataSourceDelegate
- (NSInteger)numberOfRowsInTableView:(NSTableView *)tableView {
    return [self.packageList count];
}

- (nullable id)tableView:(NSTableView *)tableView objectValueForTableColumn:(nullable NSTableColumn *)tableColumn row:(NSInteger)row {
    return [self.packageList objectAtIndex:row];
}

#pragma mark - NSTableViewDelegate
-(NSView *)tableView:(NSTableView *)aTableView viewForTableColumn:(NSTableColumn *)aTableColumn row:(NSInteger)row{
    PCPackageMeta *meta = [self.packageList objectAtIndex:row];
    NSTableCellView *nv = [aTableView makeViewWithIdentifier:@"packageview" owner:self];
    [nv.textField setStringValue:[meta packageDescription]];
    return nv;
}

- (BOOL)selectionShouldChangeInTableView:(NSTableView *)tableView {
    return NO;
}

- (BOOL)tableView:(NSTableView *)tableView shouldSelectRow:(NSInteger)row {
    return NO;
}

#pragma mark - PCTaskDelegate
-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    
    if(aTask.terminationStatus != 0) {
        [self resetUIForFailure];
        [self.progressLabel setStringValue:@"Installation Error. Please try again."];
        
        self.javaTask = nil;
        self.saltMasterTask = nil;
        self.saltMinionTask = nil;
        return;
    }
    
    [self setUIToProceedState];

    if(self.javaTask == aPCTask){
        WEAK_SELF(self);
        [[NSOperationQueue mainQueue] addOperationWithBlock:^{
            [belf downloadMetaFiles];
        }];
        
        self.javaTask = nil;
    }

    if(self.saltMasterTask == aPCTask ){
        [self setProgMessage:@"Setting up slave nodes..." value:80.0];

        PCTask *smt = [PCTask new];
        smt.taskCommand = [NSString stringWithFormat:@"salt \'pc-node*\' state.sls hadoop/2-4-0/datanode/cluster/install"];
        smt.delegate = self;
        self.saltMinionTask = smt;
        [smt launchTask];

        self.saltMasterTask = nil;
    }
    
    if(self.saltMinionTask == aPCTask){
        
        NSUInteger nc = 3;
        PCClusterType t = [[Util getApp] loadClusterType];
        switch (t) {
            case PC_CLUTER_VAGRANT:{
                nc = 3;
                break;
            }
            case PC_CLUSTER_RASPBERRY: {
                nc = [[[[RaspberryManager sharedManager] clusters] objectAtIndex:0] raspberryCount];
                break;
            }
            case PC_CLUSTER_NONE:
            default:
                break;
        }
        
        PCTask *smc = [PCTask new];
        smc.taskCommand = [NSString stringWithFormat:@"salt \'pc-master\' state.sls hadoop/2-4-0/namenode/cluster/complete pillar=\'{numnodes: %ld}\'",nc];
        smc.delegate = self;
        self.saltMasterCompleteTask = smc;
        [smc launchTask];

        self.saltMinionTask = nil;
    }
    
    
    if (self.saltMasterCompleteTask == aPCTask) {
        
        [self finalizeInstallProcess];

        self.saltMasterCompleteTask = nil;
    }
        
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {
    NSData *data = [aFileHandler availableData];
    NSString *str = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
    
    Log(@"STR %@",str);
}

-(BOOL)task:(PCTask *)aPCTask isOutputClosed:(id<PCTaskDelegate>)aDelegate {
    return NO;
}

#pragma mark - DPSetupWindowDelegate
-(void)resetToInitialState {
    self.canContinue = NO;
    self.canGoBack = NO;
}

- (NSString *)continueButtonTitle {
    return @"Finish";
}

#pragma mark - Setup UI status
- (void)setUIToProceedState {
    [self resetToInitialState];
    
    [self.installBtn setEnabled:NO];

    [self.circularProgress startAnimation:nil];
}

-(void)resetUIForFailure {
    [self resetToInitialState];
    
    [self.installBtn setEnabled:YES];
    
    [self.circularProgress stopAnimation:nil];
    [self.progressBar setDoubleValue:0.0];
    [self.progressBar displayIfNeeded];
}

-(void)setToNextStage {
    self.canContinue = YES;
    self.canGoBack = NO;
    [self.installBtn setEnabled:NO];
    [self.circularProgress stopAnimation:nil];
}

-(void)setProgMessage:(NSString *)aMessage value:(double)aValue {
    [self.progressLabel setStringValue:aMessage];
    [self.progressBar setDoubleValue:aValue];
    [self.progressBar displayIfNeeded];
}

#pragma mark - Install Start
-(void)startPackageInstall {
    [self setToNextStage];
    [self setProgMessage:@"Setting up master node..." value:60.0];
    
    PCTask *smt = [PCTask new];
    smt.taskCommand = [NSString stringWithFormat:@"salt \'pc-master\' state.sls hadoop/2-4-0/namenode/cluster/install"];
    smt.delegate = self;
    self.saltMasterTask = smt;
    [smt launchTask];
}

-(void)failedPackageInstall {
    [self resetUIForFailure];
}

-(void)finalizeInstallProcess {

    PCClusterType t = [[Util getApp] loadClusterType];
    switch (t) {
        case PC_CLUTER_VAGRANT:{
            [[Util getApp] startVagrantSetupService];
            break;
        }
        case PC_CLUSTER_RASPBERRY: {
            [[Util getApp] startRaspberrySetupService];
            break;
        }
        case PC_CLUSTER_NONE:
        default:
            break;
    }

    [self setProgMessage:@"Installation completed!" value:100.0];
    [self setToNextStage];
}

- (void)downloadMetaFiles {
    WEAK_SELF(self);
    [belf setProgMessage:@"Downloading a meta package..." value:40.0];

    for(PCPackageMeta *meta in belf.packageList){
        
        if(!belf){
            return;
        }

        NSString *mpath = [meta.masterDownloadPath objectAtIndex:0];
        NSString *npath = [meta.nodeDownloadPath objectAtIndex:0];
        NSString *mBasePath = [NSString stringWithFormat:@"%@/%@",kPOCKET_CLUSTER_SALT_STATE_PATH ,mpath];
        NSString *nBasePath = [NSString stringWithFormat:@"%@/%@",kPOCKET_CLUSTER_SALT_STATE_PATH ,npath];
        
        [PCPackageMeta makeIntermediateDirectories:mBasePath];
        [PCPackageMeta makeIntermediateDirectories:nBasePath];
        
        [PCPackageMeta
         packageFileListOn:mpath
         WithBlock:^(NSArray<NSString *> *mFileList, NSError *mError) {
             
             if(belf && !mError){
                 @synchronized(belf.downloadFileList) {
                     [belf.downloadFileList addObjectsFromArray:mFileList];
                 }
                 
                 [PCPackageMeta
                  packageFileListOn:npath
                  WithBlock:^(NSArray<NSString *> *nFileList, NSError *nError) {
                      
                      if(belf && !nError){
                          @synchronized(belf.downloadFileList) {
                              [belf.downloadFileList addObjectsFromArray:nFileList];
                          }
                          
                          for(NSString *mFile in mFileList){
                              [PCPackageMeta
                               downloadFileFromURL:mFile
                               basePath:mBasePath
                               completion:^(NSString *URL, NSURL *filePath) {
                                   
                                   if(belf){
                                       @synchronized(belf.downloadFileList) {
                                           [belf.downloadFileList removeObject:URL];
                                           
                                           if(![belf.downloadFileList count]){
                                               [[NSOperationQueue mainQueue] addOperationWithBlock:^{
                                                   if(belf){
                                                       [belf startPackageInstall];
                                                   }
                                               }];
                                           }
                                           
                                       }
                                       Log(@"%@ %ld",filePath, [belf.downloadFileList count]);
                                   }
                                   
                               }
                               onError:^(NSString *URL, NSError *error) {
                                   Log(@"Master - %@",[error description]);
                                   [self failedPackageInstall];
                               }];
                          }
                          
                          for(NSString *nFile in nFileList){
                              [PCPackageMeta
                               downloadFileFromURL:nFile
                               basePath:nBasePath
                               completion:^(NSString *URL, NSURL *filePath) {
                                   
                                   if(belf){
                                       @synchronized(belf.downloadFileList) {
                                           [belf.downloadFileList removeObject:URL];
                                           
                                           if(![belf.downloadFileList count]){
                                               [[NSOperationQueue mainQueue] addOperationWithBlock:^{
                                                   if(belf){
                                                       [belf startPackageInstall];
                                                   }
                                               }];
                                           }
                                           
                                       }
                                       Log(@"%@ %ld",filePath, [belf.downloadFileList count]);
                                   }
                                   
                               }
                               onError:^(NSString *URL, NSError *error) {
                                   Log(@"Node - %@",[error description]);
                                   [self failedPackageInstall];
                               }];
                          }
                      }
                  }];
             }
         }];
    }
}


#pragma mark - IBACTION
-(IBAction)install:(id)sender {
    
    // if there is no package to install, just don't do it.
    if(![self.packageList count]){
        return;
    }

    [self setUIToProceedState];
    [self setProgMessage:@"Installing Java to slave nodes..." value:10.0];

    // start java installation task
    PCTask *jt = [PCTask new];
    PCClusterType t = [[Util getApp] loadClusterType];
    switch (t) {
        case PC_CLUTER_VAGRANT:{
            jt.taskCommand = @"salt 'pc-node*' state.sls 'base/oracle-java8'";
            break;
        }
        case PC_CLUSTER_RASPBERRY: {
            jt.taskCommand = @"salt 'pc-node*' state.sls 'base/openjdk-7'";
            break;
        }
        case PC_CLUSTER_NONE:
        default:{
            [self resetUIForFailure];
            return;
        }
    }

    jt.delegate = self;
    self.javaTask = jt;
    [jt launchTask];
}
@end
