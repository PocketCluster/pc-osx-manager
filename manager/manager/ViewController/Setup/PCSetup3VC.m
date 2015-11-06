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

@interface PCSetup3VC()<PCTaskDelegate>
@property (strong, nonatomic) NSDictionary *progDict;
@property (nonatomic, strong) NSMutableArray<PCPackageMeta *> *packageList;
@property (nonatomic, strong) NSMutableArray<NSString *> *downloadFileList;
@property (nonatomic, strong) PCTask *saltMasterTask;
@property (nonatomic, strong) PCTask *saltMinionTask;
@property (readwrite, nonatomic) BOOL canContinue;
@property (readwrite, nonatomic) BOOL canGoBack;

-(void)setUIToProceedState;
-(void)resetUIForFailure;
-(void)setToNextStage;

-(void)startPackageInstall;
-(void)failedPackageInstall;
-(void)finalizeInstallProcess;
@end

@implementation PCSetup3VC

@synthesize canContinue;
@synthesize canGoBack;

-(instancetype)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil {
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if(self){
        self.packageList = [NSMutableArray arrayWithCapacity:0];
        self.downloadFileList = [NSMutableArray arrayWithCapacity:0];
        
        self.progDict = @{@"SUDO_SETUP_STEP_0":@[@"Setting up base configuration.",@20.0]
                          ,@"SUDO_SETUP_DONE":@[@"Finishing configuration.",@30.0]
                          ,@"USER_SETUP_STEP_0":@[@"Starting Vagrant.",@40.0]
                          ,@"USER_SETUP_STEP_1":@[@"Setting up connection.",@50.0]
                          ,@"USER_SETUP_STEP_2":@[@"Finalizing...",@90.0]
                          ,@"USER_SETUP_DONE":@[@"Done!",@100.0]};

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
        
        self.saltMasterTask = nil;
        self.saltMinionTask = nil;
        return;
    }
    
    [self setUIToProceedState];


    if(self.saltMasterTask == aPCTask ){
        
        PCTask *smt = [PCTask new];
        smt.taskCommand = [NSString stringWithFormat:@"salt 'pc-node*' state.sls hadoop/2-4-0/datanode/cluster/init"];
        smt.delegate = self;
        self.saltMinionTask = smt;
        [smt launchTask];

        self.saltMasterTask = nil;
    }
    
    if(self.saltMinionTask == aPCTask){
        
        [self finalizeInstallProcess];
        
        self.saltMinionTask = nil;
    }
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {
    NSData *data = [aFileHandler availableData];
    NSString *str = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
    
    Log(@"STR %@",str);
    
    NSArray *p = nil;
    for (NSString *key in self.progDict) {
        if ([str containsString:key]){
            p = [self.progDict valueForKey:key];
            break;
        }
    }
    
    if(p != nil){
        [self.progressLabel setStringValue:[p objectAtIndex:0]];
        [self.progressBar setDoubleValue:[[p objectAtIndex:1] doubleValue]];
        [self.progressBar displayIfNeeded];
    }
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
    [self.progressBar setDoubleValue:100.0];
    [self.progressBar displayIfNeeded];
}

#pragma mark - Install Start
-(void)startPackageInstall {
    [self setToNextStage];
    
    PCTask *smt = [PCTask new];
    smt.taskCommand = [NSString stringWithFormat:@"salt 'pc-master' state.sls hadoop/2-4-0/namenode/cluster/init"];
    smt.delegate = self;
    self.saltMasterTask = smt;
    [smt launchTask];
}

-(void)failedPackageInstall {
    [self resetUIForFailure];
}

-(void)finalizeInstallProcess {

    PCClusterType t = [[Util getApp] loadClusterType];
    t = PC_CLUTER_VAGRANT;
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

    [self setToNextStage];
}

#pragma mark - IBACTION
-(IBAction)install:(id)sender {
    
    WEAK_SELF(self);
    
    // if there is no package to install, just don't do it.
    if(![self.packageList count]){
        return;
    }

    [self setUIToProceedState];
    
    [self.progressBar setDoubleValue:10.0];
    [self.progressBar displayIfNeeded];

    
    for(PCPackageMeta *meta in belf.packageList){

        if(!belf){
            return;
        }
        
        NSString *mpath = [meta.masterFilePath objectAtIndex:0];
        NSString *npath = [meta.nodeFilePath objectAtIndex:0];
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
@end




