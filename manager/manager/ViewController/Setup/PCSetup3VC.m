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
@property (nonatomic, strong) NSMutableArray<PCPackageMeta *> *packageList;
@property (nonatomic, strong) NSMutableArray<NSString *> *downloadFileList;
@end

@implementation PCSetup3VC

-(instancetype)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil {
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if(self){
        self.packageList = [NSMutableArray arrayWithCapacity:0];
        self.downloadFileList = [NSMutableArray arrayWithCapacity:0];

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
    
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {
    
    NSData *data = [aFileHandler availableData];
    __block NSString *str = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
    
    Log(@"STR %@",str);
}

-(BOOL)task:(PCTask *)aPCTask isOutputClosed:(id<PCTaskDelegate>)aDelegate {
    return NO;
}

#pragma mark - IBACTION
-(IBAction)install:(id)sender {

    WEAK_SELF(self);

    // if there is no package to install, just don't do it.
    if(![self.packageList count]){
        return;
    }
    
    for(PCPackageMeta *meta in belf.packageList){
        
        if(!belf){
            return;
        }

        __block NSString *mpath = [meta.masterFilePath objectAtIndex:0];
        [PCPackageMeta makeIntermediateDirectories:mpath];
        
        [PCPackageMeta
         packageFileListOn:mpath
         WithBlock:^(NSArray<NSString *> *fileList, NSError *error) {
             
             if(belf){
                 @synchronized(belf.downloadFileList) {
                     [belf.downloadFileList addObjectsFromArray:fileList];
                 }
             }
             
             for(NSString *file in fileList){
                 [PCPackageMeta
                  downloadFileFromURL:file
                  basePath:[NSString stringWithFormat:@"%@/%@",kPOCKET_CLUSTER_SALT_STATE_PATH ,mpath]
                  completion:^(NSString *URL, NSURL *filePath) {
                      
                      if(belf){
                          @synchronized(belf.downloadFileList) {
                              [belf.downloadFileList removeObject:URL];
                          }
                          Log(@"%@ %ld",filePath, [belf.downloadFileList count]);
                      }
                  }
                  onError:^(NSString *URL, NSError *error) {
                      Log(@"%@",[error description]);
                  }];
             }
         }];
        
        __block NSString *npath = [meta.nodeFilePath objectAtIndex:0];
        [PCPackageMeta makeIntermediateDirectories:npath];
        [PCPackageMeta
         packageFileListOn:npath
         WithBlock:^(NSArray<NSString *> *fileList, NSError *error) {
             
             if(belf){
                 @synchronized(belf.downloadFileList) {
                     [belf.downloadFileList addObjectsFromArray:fileList];
                 }
             }
             
             for(NSString *file in fileList){
                 [PCPackageMeta
                  downloadFileFromURL:file
                  basePath:[NSString stringWithFormat:@"%@/%@",kPOCKET_CLUSTER_SALT_STATE_PATH ,npath]
                  completion:^(NSString *URL, NSURL *filePath) {
                      
                      if(belf){
                          @synchronized(belf.downloadFileList) {
                              [belf.downloadFileList removeObject:URL];
                          }
                          Log(@"%@ %ld",filePath, [belf.downloadFileList count]);
                      }
                  }
                  onError:^(NSString *URL, NSError *error) {
                      Log(@"%@",[error description]);
                  }];
             }
         }];
    }
    
    
}


@end
