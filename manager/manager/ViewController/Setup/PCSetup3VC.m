//
//  PCSetup3VC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCSetup3VC.h"

#import "PCConstants.h"
#import "VagrantManager.h"
#import "RaspberryManager.h"
#import "PCPackageManager.h"
#import "PCTask.h"
#import "Util.h"

typedef enum INSTALL_PROGRESS {
    
    IP_SALT_INIT_JOB_CHECKER = 0,
    IP_SALT_MASTER_INSTALL,
    IP_SALT_NODE_INSTALL,
    IP_SALT_MASTER_COMPLETE,
    IP_FINALIZE_INSTALL

} INSTALL_PROGRESS;


@interface PCSetup3VC()<PCTaskDelegate>
@property (nonatomic, strong) NSMutableArray<PCPackageMeta *> *packageList;
@property (nonatomic, strong) NSMutableArray<NSString *> *downloadFileList;

@property (nonatomic, strong) PCTask *saltMasterInstallTask;
@property (nonatomic, strong) PCTask *saltMinionInstallTask;
@property (nonatomic, strong) PCTask *saltMasterCompleteTask;

@property (nonatomic, strong) PCTask *saltJobTask;

@property (readwrite, nonatomic) BOOL canContinue;
@property (readwrite, nonatomic) BOOL canGoBack;

- (void)setUIToProceedState;
- (void)resetUIForFailure;
- (void)setToNextStage;
- (void)setProgMessage:(NSString *)aMessage value:(double)aValue;

- (void)checkLiveSaltJob;
- (void)startInstallProcessWithMasterNode;
- (void)nodeInstallProcess;
- (void)masterCompleteProcess;
- (void)finalizeInstallProcess;

- (void)downloadMetaFiles;
@end

@implementation PCSetup3VC {
    INSTALL_PROGRESS _install_marker;
    BOOL _isJobStillRunning;
}
@synthesize canContinue;
@synthesize canGoBack;

-(instancetype)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil {
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if(self){
        self.packageList = [NSMutableArray arrayWithCapacity:0];
        self.downloadFileList = [NSMutableArray arrayWithCapacity:0];

        _install_marker = IP_SALT_INIT_JOB_CHECKER;
        _isJobStillRunning = NO;
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
    
    int term_status = aTask.terminationStatus;
    
    if(self.saltJobTask == aPCTask) {
        
        // need to invalidate job task first
        self.saltJobTask = nil;
        
        // check if a salt job is still running
        if(_isJobStillRunning){
            
            [self checkLiveSaltJob];
            
        // no job is running. let's proceed
        }else{
            
            switch (_install_marker) {
                case IP_SALT_INIT_JOB_CHECKER:{

                    _install_marker = IP_SALT_MASTER_INSTALL;
                    [self downloadMetaFiles];
                    break;
                }
                    
                case IP_SALT_MASTER_INSTALL:{

                    _install_marker = IP_SALT_NODE_INSTALL;
                    [self nodeInstallProcess];
                    break;
                }
                    
                case IP_SALT_NODE_INSTALL: {
                    
                    _install_marker = IP_SALT_MASTER_COMPLETE;
                    [self masterCompleteProcess];
                    break;
                }
                    
                case IP_SALT_MASTER_COMPLETE: {

                    _install_marker = IP_FINALIZE_INSTALL;
                    [self finalizeInstallProcess];
                    break;
                }
                case IP_FINALIZE_INSTALL:
                default:
                    break;
            }
        }
    }

    if(self.saltMasterInstallTask == aPCTask ){
        
        if(term_status == 0){
            _install_marker = IP_SALT_NODE_INSTALL;
            [self nodeInstallProcess];
        } else {
            Log(@"There is an while exec %@", [[self.packageList objectAtIndex:0].masterInstallPath objectAtIndex:0]);
            [self checkLiveSaltJob];
        }
        
        self.saltMasterInstallTask = nil;
    }
    
    if(self.saltMinionInstallTask == aPCTask){
        
        if (term_status == 0){
            _install_marker = IP_SALT_MASTER_COMPLETE;
            [self masterCompleteProcess];
        } else {
            Log(@"There is an while exec %@", [[self.packageList objectAtIndex:0].nodeInstallPath objectAtIndex:0]);
            [self checkLiveSaltJob];
        }
        
        self.saltMinionInstallTask = nil;
    }
    
    
    if (self.saltMasterCompleteTask == aPCTask) {
        
        if (term_status == 0){
            _install_marker = IP_FINALIZE_INSTALL;
            [self finalizeInstallProcess];
        } else {
            Log(@"There is an while exec %@", [[self.packageList objectAtIndex:0].masterCompletePath objectAtIndex:0]);
            [self checkLiveSaltJob];
        }
        
        self.saltMasterCompleteTask = nil;
    }
    
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {
    NSData *data = [aFileHandler availableData];
    NSString *str = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
    
    // check if job id is found
    if(self.saltJobTask == aPCTask){
        
        //TODO: this is really important piece of code. Is this right Job ID?
        //NSRange range = [str rangeOfString:@"^[0-9]{20}\\:$" options:NSRegularExpressionSearch];
        NSRange range = [str rangeOfString:@"[0-9]{20}\\:" options:NSRegularExpressionSearch];

        if (range.location != NSNotFound){
            Log(@"\tSALT JOB IS STILL RUNNING!!!");
            _isJobStillRunning = YES;
        }
    }
    
    Log(@"%@",str);
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
    [self.progressLabel setStringValue:@"Installation Error. Please try again."];
    [self.circularProgress stopAnimation:nil];
    [self.progressBar setDoubleValue:0.0];
    [self.progressBar displayIfNeeded];

}

-(void)setToNextStage {
    
    self.canContinue = YES;
    self.canGoBack = NO;
    
    [self setProgMessage:@"Installation completed!" value:100.0];
    [self.installBtn setEnabled:NO];
    [self.circularProgress stopAnimation:nil];
}

-(void)setProgMessage:(NSString *)aMessage value:(double)aValue {

    [self.circularProgress startAnimation:nil];
    [self.progressLabel setStringValue:aMessage];
    [self.progressBar setDoubleValue:aValue];
    [self.progressBar displayIfNeeded];

}

#pragma mark - INSTALL FLOW CONTROL

- (void)checkLiveSaltJob {

    _isJobStillRunning = NO;

    PCTask *clsjt = [PCTask new];
    clsjt.taskCommand = @"salt-run jobs.active";
    clsjt.delegate = self;
    self.saltJobTask = clsjt;
    
    [clsjt performSelector:@selector(launchTask) withObject:nil afterDelay:5.0];
}

-(void)startInstallProcessWithMasterNode {

    [self setProgMessage:@"Setting up master node..." value:40.0];
    
    PCPackageMeta *meta = [self.packageList objectAtIndex:0];
    PCTask *smt = [PCTask new];
    smt.taskCommand = [NSString stringWithFormat:@"salt \'pc-master\' state.sls %@",[meta.masterInstallPath objectAtIndex:0]];
    smt.delegate = self;
    self.saltMasterInstallTask = smt;

    [smt performSelector:@selector(launchTask) withObject:nil afterDelay:1.0];
}

- (void)nodeInstallProcess {

    [self setProgMessage:@"Setting up slave nodes..." value:60.0];
    
    PCPackageMeta *meta = [self.packageList objectAtIndex:0];
    PCTask *smt = [PCTask new];
    smt.taskCommand = [NSString stringWithFormat:@"salt \'pc-node*\' state.sls %@", [meta.nodeInstallPath objectAtIndex:0]];
    smt.delegate = self;
    self.saltMinionInstallTask = smt;
    [smt performSelector:@selector(launchTask) withObject:nil afterDelay:1.0];
}


- (void)masterCompleteProcess {

    [self setProgMessage:@"Finishing master node..." value:80.0];
    
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
    
    PCPackageMeta *meta = [self.packageList objectAtIndex:0];
    PCTask *smc = [PCTask new];
    smc.taskCommand = [NSString stringWithFormat:@"salt \'pc-master\' state.sls %@ pillar=\'{numnodes: %ld}\'", [meta.masterCompletePath objectAtIndex:0], nc];
    smc.delegate = self;
    self.saltMasterCompleteTask = smc;
    [smc performSelector:@selector(launchTask) withObject:nil afterDelay:1.0];
    
}


-(void)finalizeInstallProcess {
    
    //TODO: this needs to be fixed. the UUID or id should come from cluster itself
    PCPackageMeta *meta = [self.packageList objectAtIndex:0];

    PCClusterType t = [[Util getApp] loadClusterType];
    switch (t) {
        case PC_CLUTER_VAGRANT:{

            // FIXME : if vagrant instances are not refreshed, you cannot have an instance at this point. fix this.
            //VagrantInstance *instance = [[[VagrantManager sharedManager] getInstances] objectAtIndex:0];
            //NSString *cr = [NSString stringWithFormat:@"%@-%@-%@",instance.providerIdentifier, instance.path, instance.displayName];

            meta.clusterRelation = @"virtualbox-/pocket/boxes-Cluster 1";

            // installed package data should be available before registration begins
            [[PCPackageManager sharedManager] addInstalledPackage:meta];
            [[PCPackageManager sharedManager] saveInstalledPackage];
            
            [[Util getApp] startVagrantSetupService];
            break;
        }
        case PC_CLUSTER_RASPBERRY: {
            
            RaspberryCluster *cluster = [[[RaspberryManager sharedManager] clusters] objectAtIndex:0];
            meta.clusterRelation = cluster.clusterId;

            // installed package data should be available before registration begins
            [[PCPackageManager sharedManager] addInstalledPackage:meta];
            [[PCPackageManager sharedManager] saveInstalledPackage];
            
            [[Util getApp] startRaspberrySetupService];
            break;
        }
        case PC_CLUSTER_NONE:
        default:
            break;
    }
    
    [self setToNextStage];
}

- (void)downloadMetaFiles {
    WEAK_SELF(self);
    [belf setProgMessage:@"Downloading a meta package..." value:20.0];

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
                 
                 [belf.downloadFileList addObjectsFromArray:mFileList];
                 
                 [PCPackageMeta
                  packageFileListOn:npath
                  WithBlock:^(NSArray<NSString *> *nFileList, NSError *nError) {
                      
                      if(belf && !nError){

                          [belf.downloadFileList addObjectsFromArray:nFileList];
                          
                          for(NSString *mFile in mFileList){
                              [PCPackageMeta
                               downloadFileFromURL:mFile
                               basePath:mBasePath
                               completion:^(NSString *URL, NSURL *filePath) {
                                   
                                   if(belf){

                                       [belf.downloadFileList removeObject:URL];
                                       
Log(@"%@ %ld",URL, [belf.downloadFileList count]);
                                       
                                       if([belf.downloadFileList count] == 0){
                                           [belf performSelector:@selector(startInstallProcessWithMasterNode) withObject:nil afterDelay:0.0];
                                       }

                                   }
                               }
                               onError:^(NSString *URL, NSError *error) {
                                   Log(@"Master - %@",[error description]);
                                   [belf resetUIForFailure];
                               }];
                          }
                          
                          for(NSString *nFile in nFileList){
                              [PCPackageMeta
                               downloadFileFromURL:nFile
                               basePath:nBasePath
                               completion:^(NSString *URL, NSURL *filePath) {
                                   
                                   if(belf){
                                       [belf.downloadFileList removeObject:URL];

Log(@"%@ %ld",URL, [belf.downloadFileList count]);
                                       
                                       if([belf.downloadFileList count] == 0){
                                          [belf performSelector:@selector(startInstallProcessWithMasterNode) withObject:nil afterDelay:0.0];
                                       }
                                   }
                               }
                               onError:^(NSString *URL, NSError *error) {
                                   Log(@"Node - %@",[error description]);
                                   [belf resetUIForFailure];
                               }];
                          }
                          
                      } else {
                          [belf resetUIForFailure];
                      }
                  }];
                 
             }else{
                [belf resetUIForFailure];
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
    [self setProgMessage:@"Check cluster status..." value:10.0];

    [self checkLiveSaltJob];
}
@end
