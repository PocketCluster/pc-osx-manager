//
//  PCPkgInstallWC.m
//  manager
//
//  Created by Almighty Kim on 11/24/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCPkgInstallWC.h"
#import "PCConstants.h"
//#import "PCPackageManager.h"
//#import "PCTask.h"
#import "Util.h"

#import "Package.h"
#import "Router.h"

#define BASE_PROGRESS_PERCENTAGE  (30.0)
#define TOTAL_PROGRESS_DURATION  (60.0)

@interface PCPkgInstallWC ()
//<PCTaskDelegate>
@property (nonatomic, strong) NSMutableArray<Package *> *packageList;
@property (nonatomic, strong) NSMutableArray<NSString *> *downloadFileList;

@property (nonatomic, strong) id saltMasterInstallTask;
@property (nonatomic, strong) id saltSecondInstallTask;
@property (nonatomic, strong) id saltNodeInstallTask;

@property (nonatomic, strong) id saltMasterCompleteTask;
@property (nonatomic, strong) id saltSecondCompleteTask;
@property (nonatomic, strong) id saltNodeCompleteTask;

@property (nonatomic, strong) id saltJobTask;

- (NSUInteger)getNodeCount;
- (double)getDeltaProgress:(double)aProgressMark;

- (void)resetToInitialState;
- (void)setUIToProceedState;
- (void)resetUIForFailure;
- (void)setToNextStage;
- (void)setProgMessage:(NSString *)aMessage value:(double)aValue;

- (void)checkLiveSaltJob;
- (void)proceedTargetStage;
- (void)startInstallProcessForMaster;
- (void)startInstallProcessForSecondary;
- (void)startInstallProcessForNode:(NSUInteger)aStartNode;

- (void)startCompletionForMaster;
- (void)startCompletionForSecondary;
- (void)startCompletionForNode:(NSUInteger)aStartNode;
- (void)finalizeInstallProcess;

- (void)downloadMetaFiles;
@end

@implementation PCPkgInstallWC {
    PKG_INSTALL_PROGRESS _install_marker;
    BOOL _isJobStillRunning;
    NSUInteger _target_package_index;
    NSUInteger _total_target_count;
    NSUInteger _target_installed;
}

- (void)windowDidLoad {
    [super windowDidLoad];
    // Implement this method to handle any initialization after your window controller's window has been loaded from its nib file.
}

-(instancetype)initWithWindowNibName:(NSString *)windowNibName {

    self = [super initWithWindowNibName:windowNibName];
    if(self){

        self.packageList = [NSMutableArray arrayWithCapacity:0];
        self.downloadFileList = [NSMutableArray arrayWithCapacity:0];
        
        _install_marker = PI_INIT_JOB_CHECKER;
        _isJobStillRunning = NO;
        _target_package_index = 0;
        [self resetToInitialState];

#if 0
        // TODO move this process to package manager or somewhere to make it more formalized
        WEAK_SELF(self);
        [PCPackageMeta metaPackageListWithBlock:^(NSArray<PCPackageMeta *> *packages, NSError *error) {
            if(belf != nil){

                for (PCPackageMeta *meta in packages){
                    for (PCPackageMeta *installed in [[PCPackageManager sharedManager] installedPackage]) {
                        if ([meta.packageId isEqualToString:installed.packageId ]) {
                            [meta setInstalled:YES];
                        }
                    }
                }
                
                [belf.packageList addObjectsFromArray:packages];
                [belf.packageTable reloadData];
            }
        }];
#endif
    }
    return self;
}

#pragma mark - NSTableViewDataSourceDelegate
- (NSInteger)numberOfRowsInTableView:(NSTableView *)tableView {
    return [self.packageList count];
}

- (nullable id)tableView:(NSTableView *)tableView objectValueForTableColumn:(nullable NSTableColumn *)tableColumn row:(NSInteger)row {
    return [[self.packageList objectAtIndex:row] description];
}

#pragma mark - NSTableViewDelegate
-(NSView *)tableView:(NSTableView *)aTableView viewForTableColumn:(NSTableColumn *)aTableColumn row:(NSInteger)row{
#if 0
    Package *meta = [self.packageList objectAtIndex:row];
    NSTableCellView *nv = [aTableView makeViewWithIdentifier:@"packageview" owner:self];
    [nv.textField setStringValue:[meta packageDescription]];

    if([meta isInstalled]){
        [nv.textField setTextColor:[NSColor lightGrayColor]];
    }else{
        [nv.textField setTextColor:[NSColor blackColor]];
    }

    return nv;
#endif

    return [aTableView makeViewWithIdentifier:@"packageview" owner:self];
}

- (BOOL)selectionShouldChangeInTableView:(NSTableView *)tableView {
    return NO;
}

- (BOOL)tableView:(NSTableView *)tableView shouldSelectRow:(NSInteger)row {
    return NO;
}

#if 0
#pragma mark - PCTaskDelegate
-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    
    //FIXME: should always check active job status
    //int term_status = aTask.terminationStatus;
    int term_status = -1;

    if(self.saltJobTask == aPCTask) {
        
        // need to invalidate job task first
        self.saltJobTask = nil;
        
        // check if a salt job is still running
        if(_isJobStillRunning){
            
            [self checkLiveSaltJob];

        // no job is running. let's proceed
        }else{
            
            switch (_install_marker) {
                case PI_INIT_JOB_CHECKER:{
                    
                    //FIXME: this could be a root of serious bug!
                    _install_marker = PI_MASTER_INSTALL;
                    [self downloadMetaFiles];
                    break;
                }

                case PI_MASTER_INSTALL:{
                    
                    // if secondary install script exists
                    if ([[self.packageList objectAtIndex:_target_package_index].secondaryInstallPath count]){
                        
                        _install_marker = PI_SECONDARY_INSTALL;
                        [self startInstallProcessForSecondary];
                        
                    // if secondary install script DNE
                    }else{

                        // if node install script exists
                        if([[self.packageList objectAtIndex:_target_package_index].nodeInstallPath count]){
                            _install_marker = PI_NODE_INSTALL;
                            [self startInstallProcessForNode:1];
                            
                        // if node install script DNE
                        }else{
                            _install_marker = PI_MASTER_COMPLETE;
                            [self startCompletionForMaster];
                        }

                    }

                    break;
                }
                    
                    
                case PI_SECONDARY_INSTALL:{
                    _install_marker = PI_NODE_INSTALL;
                    [self startInstallProcessForNode:2];
                    break;
                }

                case PI_NODE_INSTALL: {
                    _install_marker = PI_MASTER_COMPLETE;
                    [self startCompletionForMaster];
                    break;
                }
                    
                case PI_MASTER_COMPLETE: {
                    
                    // if secondary complete script exists
                    if ([[self.packageList objectAtIndex:_target_package_index].secondaryCompletePath count]){
                        
                        _install_marker = PI_SECONDARY_COMPLETE;
                        [self startCompletionForSecondary];
                        
                    // if secondary complete script DNE
                    }else{

                        // if node complete script exists
                        if([[self.packageList objectAtIndex:_target_package_index].nodeCompletePath count]){

                            _install_marker = PI_NODE_COMPLETE;
                            [self startCompletionForNode:1];
                            
                        // if node complete script DNE
                        }else{

                            _install_marker = PI_FINALIZE_INSTALL;
                            [self finalizeInstallProcess];
                        }
                    }
                    break;
                }
                    
                case PI_SECONDARY_COMPLETE: {
                    _install_marker = PI_NODE_COMPLETE;
                    [self startCompletionForNode:2];
                    break;
                }
                    
                case PI_NODE_COMPLETE: {
                    _install_marker = PI_FINALIZE_INSTALL;
                    [self finalizeInstallProcess];
                    break;
                }
                
                case PI_FINALIZE_INSTALL:
                default:
                    break;
            }
        }
    }
    

    if(self.saltMasterInstallTask == aPCTask ){
        
        if(term_status == 0){
            
            // if secondary install script exists
            if ([[self.packageList objectAtIndex:_target_package_index].secondaryInstallPath count]){
                
                _install_marker = PI_SECONDARY_INSTALL;
                [self startInstallProcessForSecondary];
                
            // if secondary install script DNE
            }else{
                
                // if node install script exists
                if([[self.packageList objectAtIndex:_target_package_index].nodeInstallPath count]){
                    _install_marker = PI_NODE_INSTALL;
                    [self startInstallProcessForNode:1];
                    
                    // if node install script DNE
                }else{
                    _install_marker = PI_MASTER_COMPLETE;
                    [self startCompletionForMaster];
                }
                
            }

        } else {
            Log(@"There is an while exec %@", [[self.packageList objectAtIndex:_target_package_index].masterInstallPath objectAtIndex:0]);
            [self checkLiveSaltJob];
        }
        
        self.saltMasterInstallTask = nil;
    }
    
    
    if(self.saltSecondInstallTask == aPCTask){

        if(term_status == 0){
            _install_marker = PI_NODE_INSTALL;
            [self startInstallProcessForNode:2];
        }else{
            Log(@"There is an while exec %@", [[self.packageList objectAtIndex:_target_package_index].secondaryInstallPath objectAtIndex:0]);
            [self checkLiveSaltJob];
        }

        self.saltSecondInstallTask = nil;
    }
    
    if(self.saltNodeInstallTask == aPCTask){
        
        if (term_status == 0){
            _install_marker = PI_MASTER_COMPLETE;
            [self startCompletionForMaster];
        } else {
            Log(@"There is an while exec %@", [[self.packageList objectAtIndex:_target_package_index].nodeInstallPath objectAtIndex:0]);
            [self checkLiveSaltJob];
        }
        
        self.saltNodeInstallTask = nil;
    }
    
    
    if (self.saltMasterCompleteTask == aPCTask) {
        
        if (term_status == 0){
            
            // if secondary complete script exists
            if ([[self.packageList objectAtIndex:_target_package_index].secondaryCompletePath count]){
                
                _install_marker = PI_SECONDARY_COMPLETE;
                [self startCompletionForSecondary];

            // if secondary complete script DNE
            }else{

                // if node complete script exists
                if([[self.packageList objectAtIndex:_target_package_index].nodeCompletePath count]){

                    _install_marker = PI_NODE_COMPLETE;
                    [self startCompletionForNode:1];

                // if node complete script DNE
                }else{
                    _install_marker = PI_FINALIZE_INSTALL;
                    [self finalizeInstallProcess];
                }
            }
            
        } else {
            Log(@"There is an while exec %@", [[self.packageList objectAtIndex:_target_package_index].masterCompletePath objectAtIndex:0]);
            [self checkLiveSaltJob];
        }
        
        self.saltMasterCompleteTask = nil;
    }
    
    
    if(self.saltSecondCompleteTask == aPCTask) {
        
        if(term_status == 0){
            _install_marker = PI_NODE_COMPLETE;
            [self startCompletionForNode:2];
        }else{
            Log(@"There is an while exec %@", [[self.packageList objectAtIndex:_target_package_index].secondaryCompletePath objectAtIndex:0]);
            [self checkLiveSaltJob];
        }

        self.saltSecondCompleteTask = nil;
    }
    
    if(self.saltNodeCompleteTask == aPCTask) {
        
        if(term_status == 0){
            _install_marker = PI_FINALIZE_INSTALL;
            [self finalizeInstallProcess];
        }else {
            Log(@"There is an while exec %@", [[self.packageList objectAtIndex:_target_package_index].nodeCompletePath objectAtIndex:0]);
            [self checkLiveSaltJob];
        }

        self.saltNodeCompleteTask = nil;
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
#endif

#pragma mark - Utils

-(NSUInteger)getNodeCount {
    NSUInteger nc = 3;

#if 0
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
            nc = 0;
            break;
    }
#endif

    return nc;
}

-(double)getDeltaProgress:(double)aProgressMark {
    double delta_duration = (TOTAL_PROGRESS_DURATION)/(double)_total_target_count;
    double delta_progress = (BASE_PROGRESS_PERCENTAGE) + (delta_duration * _target_installed) + (delta_duration * aProgressMark);
    return delta_progress;
}

#pragma mark - UI status
-(void)resetToInitialState {
    [self.installBtn setEnabled:YES];
    [self.closeBtn setEnabled:YES];
}

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
#if 0
    PCClusterType t = [[Util getApp] loadClusterType];
    switch (t) {
        case PC_CLUTER_VAGRANT:{
            [[VagrantManager sharedManager] refreshVagrantMachines];
            break;
        }
        case PC_CLUSTER_RASPBERRY: {
            [[RaspberryManager sharedManager] refreshRaspberryClusters];
            break;
        }
        case PC_CLUSTER_NONE:
        default:
            break;
    }
#endif
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
#if 0
    PCTask *clsjt = [PCTask new];
    clsjt.taskCommand = @"salt-run jobs.active";
    clsjt.delegate = self;
    self.saltJobTask = clsjt;
    
    [clsjt performSelector:@selector(launchTask) withObject:nil afterDelay:5.0];
#endif
}

//FIXME: introducing an intermediate function stage could cause hineous bug.
// Be very careful and monitor closely
-(void)proceedTargetStage {
    Log(@"%s",__PRETTY_FUNCTION__);
#if 0
    BOOL foundTarget = NO;
    
    for (NSUInteger i = 0; i < [self.packageList count]; ++i){

        PCPackageMeta *meta = [self.packageList objectAtIndex:i];
        if(meta.isInstalled){
            continue;
        }
        
        foundTarget = YES;
        _target_package_index = i;
        Log(@"next target %ld %@", i, [meta debugDescription]);
        break;
    }

    // we have no packages to install. close installation
    if(foundTarget){
        //FIXME: this could be a root of serious bug!
        _install_marker = PI_MASTER_INSTALL;
        [self startInstallProcessForMaster];
    }else{
        [self setToNextStage];
    }
#endif
}

-(void)startInstallProcessForMaster {
#if 0
    [self setProgMessage:@"Setting up master node..." value:[self getDeltaProgress:0.167]];

    NSUInteger nc = [self getNodeCount];
    if(nc == 0){return;}
    
    PCPackageMeta *meta = [self.packageList objectAtIndex:_target_package_index];
    PCTask *smt = [PCTask new];
    smt.taskCommand = [NSString stringWithFormat:@"salt \'pc-master\' state.sls %@ pillar=\'{numnodes: %ld}\'",[meta.masterInstallPath objectAtIndex:0],nc];
    smt.delegate = self;
    self.saltMasterInstallTask = smt;
    
    [smt performSelector:@selector(launchTask) withObject:nil afterDelay:1.0];
#endif
}

-(void)startInstallProcessForSecondary {
#if 0
    [self setProgMessage:@"Setting up secondary node..." value:[self getDeltaProgress:0.334]];
    
    NSUInteger nc = [self getNodeCount];
    if(nc == 0){return;}
    
    PCPackageMeta *meta = [self.packageList objectAtIndex:_target_package_index];
    PCTask *smt = [PCTask new];
    smt.taskCommand = [NSString stringWithFormat:@"salt \'pc-node1\' state.sls %@ pillar=\'{numnodes: %ld}\'",[meta.secondaryInstallPath objectAtIndex:0],nc];
    smt.delegate = self;
    self.saltSecondInstallTask = smt;

    [smt performSelector:@selector(launchTask) withObject:nil afterDelay:1.0];
#endif
}

- (void)startInstallProcessForNode:(NSUInteger)aStartNode {
#if 0
    [self setProgMessage:@"Setting up slave nodes..." value:[self getDeltaProgress:0.5]];
    
    NSUInteger nc = [self getNodeCount];
    if(nc == 0){return;}

    PCPackageMeta *meta = [self.packageList objectAtIndex:_target_package_index];
    PCTask *snt = [PCTask new];
    snt.taskCommand = [NSString stringWithFormat:@"salt \'pc-node[%ld-%ld]\' state.sls %@ pillar=\'{numnodes: %ld}\'",aStartNode,nc,[meta.nodeInstallPath objectAtIndex:0],nc];
    snt.delegate = self;
    self.saltNodeInstallTask = snt;
    [snt performSelector:@selector(launchTask) withObject:nil afterDelay:1.0];
#endif
}

#pragma mark - COMPLETION FLOW CONTROL

- (void)startCompletionForMaster {
#if 0
    [self setProgMessage:@"Finishing master node..." value:[self getDeltaProgress:0.668]];
    
    NSUInteger nc = [self getNodeCount];
    if(nc == 0){return;}
    
    PCPackageMeta *meta = [self.packageList objectAtIndex:_target_package_index];
    PCTask *smc = [PCTask new];
    smc.taskCommand = [NSString stringWithFormat:@"salt \'pc-master\' state.sls %@ pillar=\'{numnodes: %ld}\'", [meta.masterCompletePath objectAtIndex:0], nc];
    smc.delegate = self;
    self.saltMasterCompleteTask = smc;
    [smc performSelector:@selector(launchTask) withObject:nil afterDelay:1.0];
#endif
}

- (void)startCompletionForSecondary {
#if 0
    [self setProgMessage:@"Finishing Secondary node..." value:[self getDeltaProgress:0.835]];
    
    NSUInteger nc = [self getNodeCount];
    if(nc == 0){return;}

    PCPackageMeta *meta = [self.packageList objectAtIndex:_target_package_index];
    PCTask *ssc = [PCTask new];
    ssc.taskCommand = [NSString stringWithFormat:@"salt \'pc-node1\' state.sls %@ pillar=\'{numnodes: %ld}\'", [meta.secondaryCompletePath objectAtIndex:0], nc];
    ssc.delegate = self;
    self.saltMasterCompleteTask = ssc;
    [ssc performSelector:@selector(launchTask) withObject:nil afterDelay:1.0];
#endif
}

- (void)startCompletionForNode:(NSUInteger)aStartNode {
#if 0
    [self setProgMessage:@"Finishing Rest of Node..." value:[self getDeltaProgress:1.0]];
    
    NSUInteger nc = [self getNodeCount];
    if(nc == 0){return;}

    PCPackageMeta *meta = [self.packageList objectAtIndex:_target_package_index];
    PCTask *snc = [PCTask new];
    snc.taskCommand = [NSString stringWithFormat:@"salt \'pc-node[%ld-%ld]\' state.sls %@ pillar=\'{numnodes: %ld}\'",aStartNode, nc, [meta.nodeCompletePath objectAtIndex:0], nc];
    snc.delegate = self;
    self.saltMasterCompleteTask = snc;
    [snc performSelector:@selector(launchTask) withObject:nil afterDelay:1.0];
#endif
}

-(void)finalizeInstallProcess {
#if 0
    //TODO: this needs to be fixed. the UUID or id should come from cluster itself
    PCPackageMeta *meta = [self.packageList objectAtIndex:_target_package_index];
    meta.installed = YES;

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

            break;
        }
        case PC_CLUSTER_RASPBERRY: {
            
            RaspberryCluster *cluster = [[[RaspberryManager sharedManager] clusters] objectAtIndex:0];
            meta.clusterRelation = cluster.clusterId;
            
            // installed package data should be available before registration begins
            [[PCPackageManager sharedManager] addInstalledPackage:meta];
            [[PCPackageManager sharedManager] saveInstalledPackage];

            break;
        }
        case PC_CLUSTER_NONE:
        default:
            break;
    }

    // count installed packages
    ++_target_installed;

    // go back to the head of iteration cycle
    [self proceedTargetStage];
#endif
}

- (void)downloadMetaFiles {
#if 0
    WEAK_SELF(self);
    [self setProgMessage:@"Downloading a meta package..." value:20.0];
    
    NSMutableArray *mtlst = [NSMutableArray array];
    __block NSMutableArray *dllst = [NSMutableArray array];
    __block BOOL hasDownloadEverFailed = NO;

    for(PCPackageMeta *meta in self.packageList){
        
        if (meta.isInstalled){
            continue;
        }
        
        NSString *mpath = [meta.masterDownloadPath objectAtIndex:0];
        NSString *mBasePath = [NSString stringWithFormat:@"%@/%@",kPOCKET_CLUSTER_SALT_STATE_PATH ,mpath];
        [PCPackageMeta makeIntermediateDirectories:mBasePath];

        id mop = [PCPackageMeta packageFileListOperation:mpath withSucess:^(NSArray<NSString *> *fileList) {

            Log(@"meta files \n%@",fileList);
            
            for (NSString *furl in fileList){
                id dop = [PCPackageMeta
                          packageFileDownloadOperation:furl
                          detinationPath:mBasePath
                          completion:^(NSString *URL, NSURL *filePath) {
                              
                              Log(@"%@ DONE",URL);
                              
                          } onError:^(NSString *URL, NSError *error) {
                              
                              Log(@"Master - %@",[error description]);
                              hasDownloadEverFailed = YES;

                          }];
                [dllst addObject:dop];
            }
            
        } withFailure:^(NSError *error) {
            
            Log(@"Master - %@",[error description]);
            hasDownloadEverFailed = YES;

        }];
        [mtlst addObject:mop];
        
        
        
        // secondary master files
        NSString *spath = nil, *sBasePath = nil;
        if( meta.secondaryDownloadPath.count ){
            
            spath = [meta.secondaryDownloadPath objectAtIndex:0];
            sBasePath = [NSString stringWithFormat:@"%@/%@",kPOCKET_CLUSTER_SALT_STATE_PATH ,spath];
            [PCPackageMeta makeIntermediateDirectories:sBasePath];

            id sop = [PCPackageMeta packageFileListOperation:spath withSucess:^(NSArray<NSString *> *fileList) {
                
                Log(@"meta files \n%@",fileList);
                
                for (NSString *furl in fileList){
                    id dop = [PCPackageMeta
                              packageFileDownloadOperation:furl
                              detinationPath:sBasePath
                              completion:^(NSString *URL, NSURL *filePath) {
                                  
                                  Log(@"%@ DONE",URL);

                              } onError:^(NSString *URL, NSError *error) {

                                  Log(@"Secondary - %@",[error description]);
                                  hasDownloadEverFailed = YES;

                              }];
                    
                    [dllst addObject:dop];
                }
                
            } withFailure:^(NSError *error) {

                Log(@"Secondary - %@",[error description]);
                hasDownloadEverFailed = YES;

            }];

            
            [mtlst addObject:sop];
        }

        
        
        NSString *npath = nil, *nBasePath = nil;
        if( meta.nodeDownloadPath.count ){
            
            npath = [meta.nodeDownloadPath objectAtIndex:0];
            nBasePath = [NSString stringWithFormat:@"%@/%@",kPOCKET_CLUSTER_SALT_STATE_PATH ,npath];
            [PCPackageMeta makeIntermediateDirectories:nBasePath];
            
            id nop = [PCPackageMeta packageFileListOperation:npath withSucess:^(NSArray<NSString *> *fileList) {

                Log(@"meta files \n%@",fileList);
                
                for (NSString *furl in fileList){
                    id dop = [PCPackageMeta
                              packageFileDownloadOperation:furl
                              detinationPath:nBasePath
                              completion:^(NSString *URL, NSURL *filePath) {

                                  Log(@"%@ DONE",URL);

                              } onError:^(NSString *URL, NSError *error) {

                                  Log(@"Node - %@",[error description]);
                                  hasDownloadEverFailed = YES;

                              }];
                    [dllst addObject:dop];
                }

            } withFailure:^(NSError *error) {

                Log(@"Node - %@",[error description]);
                hasDownloadEverFailed = YES;

            }];
            [mtlst addObject:nop];
        }
    }

    [PCPackageMeta
     batchDownloadOperation:mtlst
     progressBlock:^(NSUInteger numberOfFinishedOperations, NSUInteger totalNumberOfOperations){}
     completionBlock:^(NSArray *operations) {

         [PCPackageMeta
          batchDownloadOperation:dllst
          progressBlock:^(NSUInteger numberOfFinishedOperations, NSUInteger totalNumberOfOperations){}
          completionBlock:^(NSArray *operations) {

              Log(@"filedownload all completed");
              if(hasDownloadEverFailed){
                  [belf resetUIForFailure];
              }else{
                  [belf performSelector:@selector(proceedTargetStage) withObject:nil afterDelay:0.0];
              }
         }];
     }];
#endif
}

#pragma mark - IBACTION
-(IBAction)install:(id)sender {

    // if there is no package to install, just don't do it.
    if(![self.packageList count]){
        return;
    }

    _total_target_count = [[self.packageList filteredArrayUsingPredicate:[NSPredicate predicateWithFormat:@"(SELF.installed == NO)"]] count];
    _target_installed = 0;
    if(_total_target_count == 0){
        return;
    }
    
    [self setUIToProceedState];
    [self setProgMessage:@"Check cluster status..." value:10.0];
    
    [self checkLiveSaltJob];
}
@end
