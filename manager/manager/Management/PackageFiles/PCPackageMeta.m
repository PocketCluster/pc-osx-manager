//
//  PCMeta.m
//  AFNETPROTO
//
//  Created by Almighty Kim on 11/4/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCFormulaClient.h"
#import "PCPackageMeta.h"
#import "AFHTTPRequestOperationManager.h"
#import "AFURLSessionManager.h"
#import "PCConstants.h"

NSString * const kPCPackageRelatedCluster                   = @"package-cluster-relation";

NSString * const kPCPackageMetaVersion                      = @"pkg-ver";
NSString * const kDescription                               = @"description";
NSString * const kPCPackageName                             = @"name";
NSString * const kPCPackageFamily                           = @"family";
NSString * const kPCPackageVersions                         = @"versions";
NSString * const kPCPackageVersionId                        = @"pkg-id";
NSString * const kPCPackageVersionNumber                    = @"ver-num";
NSString * const kPCPackageVersionModes                     = @"modes";
NSString * const kPCPackageVersionModesType                 = @"mode-type";
NSString * const kPCPackageVersionPort                      = @"ports";
NSString * const kPCPackageVersionLibraryDep                = @"dep-lib";
NSString * const kPCPackageVersionBigpkgDep                 = @"dep-bigpkg";

NSString * const kPCPackageVersionMasterInstallPath         = @"path-install-master";
NSString * const kPCPackageVersionSecondaryInstallPath      = @"path-install-secondary";
NSString * const kPCPackageVersionNodeInstallPath           = @"path-install-nodes";

NSString * const kPCPackageVersionMasterCompletePath        = @"path-complete-master";
NSString * const kPCPackageVersionSecondaryCompletePath     = @"path-complete-secondary";
NSString * const kPCPackageVersionNodeCompletePath          = @"path-complete-nodes";

NSString * const kPCPackageVersionMasterDownload            = @"path-download-master";
NSString * const kPCPackageVersionSecondaryDownload         = @"path-download-secondary";
NSString * const kPCPackageVersionNodeDownload              = @"path-download-nodes";

NSString * const kPCPackageVersionMasterResetCmd            = @"cmd-reset-master";
NSString * const kPCPackageVersionSecondaryResetCmd         = @"cmd-reset-secondary";
NSString * const kPCPackageVersionNodesResetCmd             = @"cmd-reset-nodes";

NSString * const kPCPackageVersionMasterUninstallCmd        = @"cmd-uninstall-master";
NSString * const kPCPackageVersionSecondaryUninstallCmd     = @"cmd-uninstall-secondary";
NSString * const kPCPackageVersionNodesUninstallCmd         = @"cmd-uninstall-nodes";

NSString * const kPCPackageVersionStartScript               = @"script-start";
NSString * const kPCPackageVersionStopScript                = @"script-stop";
NSString * const kPCPackageVersionCmdScript                 = @"script-cmd";

NSString * const kPCPackageVersionProcessCheck              = @"check-process";

static NSString * const kGithubRawFileLinkURL               = @"download_url";

@interface PCPackageMeta()

@property (nonatomic, strong, readwrite) NSString *metaVersion;
@property (nonatomic, strong, readwrite) NSString *packageName;
@property (nonatomic, strong, readwrite) NSArray<NSString*> *family;
@property (nonatomic, strong, readwrite) NSString *packageDescription;
@property (nonatomic, strong, readwrite) NSString *packageId;
@property (nonatomic, strong, readwrite) NSString *version;
@property (nonatomic, strong, readwrite) NSString *modeType;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *ports;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *libraryDependencies;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *bigPkgDependencies;

@property (nonatomic, strong, readwrite) NSArray<NSString *> *masterInstallPath;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *secondaryInstallPath;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *nodeInstallPath;

@property (nonatomic, strong, readwrite) NSArray<NSString *> *masterCompletePath;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *secondaryCompletePath;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *nodeCompletePath;

@property (nonatomic, strong, readwrite) NSArray<NSString *> *masterDownloadPath;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *secondaryDownloadPath;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *nodeDownloadPath;

@property (nonatomic, strong, readwrite) NSArray<NSString *> *masterResetCmd;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *secondaryResetCmd;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *nodeResetCmd;

@property (nonatomic, strong, readwrite) NSArray<NSString *> *masterUninstallCmd;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *secondaryUninstallCmd;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *nodeUninstallCmd;

@property (nonatomic, strong, readwrite) NSArray<NSString *> *startScript;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *stopScript;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *cmdScript;

@property (nonatomic, strong, readwrite) NSArray<NSString *> *processCheck;

@end

#pragma mark - PCPACKAGEMETA NSCODING
@implementation PCPackageMeta (NSCoding)
- (void)encodeWithCoder:(NSCoder *)aCoder {
    
    [aCoder encodeObject:self.clusterRelation       forKey:kPCPackageRelatedCluster];
    
    [aCoder encodeObject:self.metaVersion           forKey:kPCPackageMetaVersion];
    [aCoder encodeObject:self.packageName           forKey:kPCPackageName];
    [aCoder encodeObject:self.family                forKey:kPCPackageFamily];
    [aCoder encodeObject:self.packageDescription    forKey:kDescription];

    [aCoder encodeObject:self.packageId             forKey:kPCPackageVersionId];
    [aCoder encodeObject:self.version               forKey:kPCPackageVersionNumber];
    [aCoder encodeObject:self.modeType              forKey:kPCPackageVersionModesType];
    [aCoder encodeObject:self.ports                 forKey:kPCPackageVersionPort];
    [aCoder encodeObject:self.libraryDependencies   forKey:kPCPackageVersionLibraryDep];
    [aCoder encodeObject:self.bigPkgDependencies    forKey:kPCPackageVersionBigpkgDep];
    
    [aCoder encodeObject:self.masterInstallPath     forKey:kPCPackageVersionMasterInstallPath];
    [aCoder encodeObject:self.secondaryInstallPath  forKey:kPCPackageVersionSecondaryInstallPath];
    [aCoder encodeObject:self.nodeInstallPath       forKey:kPCPackageVersionNodeInstallPath];
    
    [aCoder encodeObject:self.masterCompletePath    forKey:kPCPackageVersionMasterCompletePath];
    [aCoder encodeObject:self.secondaryCompletePath forKey:kPCPackageVersionSecondaryCompletePath];
    [aCoder encodeObject:self.nodeCompletePath      forKey:kPCPackageVersionNodeCompletePath];

    [aCoder encodeObject:self.masterDownloadPath    forKey:kPCPackageVersionMasterDownload];
    [aCoder encodeObject:self.secondaryDownloadPath forKey:kPCPackageVersionSecondaryDownload];
    [aCoder encodeObject:self.nodeDownloadPath      forKey:kPCPackageVersionNodeDownload];
    
    [aCoder encodeObject:self.masterResetCmd        forKey:kPCPackageVersionMasterResetCmd];
    [aCoder encodeObject:self.secondaryResetCmd     forKey:kPCPackageVersionSecondaryResetCmd];
    [aCoder encodeObject:self.nodeResetCmd          forKey:kPCPackageVersionNodesResetCmd];
    
    [aCoder encodeObject:self.masterUninstallCmd    forKey:kPCPackageVersionMasterUninstallCmd];
    [aCoder encodeObject:self.secondaryUninstallCmd forKey:kPCPackageVersionSecondaryUninstallCmd];
    [aCoder encodeObject:self.nodeUninstallCmd      forKey:kPCPackageVersionNodesUninstallCmd];
    
    [aCoder encodeObject:self.startScript           forKey:kPCPackageVersionStartScript];
    [aCoder encodeObject:self.stopScript            forKey:kPCPackageVersionStopScript];
    [aCoder encodeObject:self.cmdScript             forKey:kPCPackageVersionCmdScript];

    [aCoder encodeObject:self.processCheck          forKey:kPCPackageVersionProcessCheck];
}

- (instancetype)initWithCoder:(NSCoder *)aDecoder {
    self = [super init];
    if (self) {

        self.clusterRelation       = [aDecoder decodeObjectForKey:kPCPackageRelatedCluster];

        self.metaVersion           = [aDecoder decodeObjectForKey:kPCPackageMetaVersion];
        self.packageName           = [aDecoder decodeObjectForKey:kPCPackageName];
        self.family                = [aDecoder decodeObjectForKey:kPCPackageFamily];
        self.packageDescription    = [aDecoder decodeObjectForKey:kDescription];

        self.packageId             = [aDecoder decodeObjectForKey:kPCPackageVersionId];
        self.version               = [aDecoder decodeObjectForKey:kPCPackageVersionNumber];
        self.modeType              = [aDecoder decodeObjectForKey:kPCPackageVersionModesType];
        self.ports                 = [aDecoder decodeObjectForKey:kPCPackageVersionPort];
        self.libraryDependencies   = [aDecoder decodeObjectForKey:kPCPackageVersionLibraryDep];
        self.bigPkgDependencies    = [aDecoder decodeObjectForKey:kPCPackageVersionBigpkgDep];
        
        self.masterInstallPath     = [aDecoder decodeObjectForKey:kPCPackageVersionMasterInstallPath];
        self.secondaryInstallPath  = [aDecoder decodeObjectForKey:kPCPackageVersionSecondaryInstallPath];
        self.nodeInstallPath       = [aDecoder decodeObjectForKey:kPCPackageVersionNodeInstallPath];
        
        self.masterCompletePath    = [aDecoder decodeObjectForKey:kPCPackageVersionMasterCompletePath];
        self.secondaryCompletePath = [aDecoder decodeObjectForKey:kPCPackageVersionSecondaryCompletePath];
        self.nodeCompletePath      = [aDecoder decodeObjectForKey:kPCPackageVersionNodeCompletePath];
        
        self.masterDownloadPath    = [aDecoder decodeObjectForKey:kPCPackageVersionMasterDownload];
        self.secondaryDownloadPath = [aDecoder decodeObjectForKey:kPCPackageVersionSecondaryDownload];
        self.nodeDownloadPath      = [aDecoder decodeObjectForKey:kPCPackageVersionNodeDownload];
        
        self.masterResetCmd        = [aDecoder decodeObjectForKey:kPCPackageVersionMasterResetCmd];
        self.secondaryResetCmd     = [aDecoder decodeObjectForKey:kPCPackageVersionSecondaryResetCmd];
        self.nodeResetCmd          = [aDecoder decodeObjectForKey:kPCPackageVersionNodesResetCmd];
        
        self.masterUninstallCmd    = [aDecoder decodeObjectForKey:kPCPackageVersionMasterUninstallCmd];
        self.secondaryUninstallCmd = [aDecoder decodeObjectForKey:kPCPackageVersionSecondaryUninstallCmd];
        self.nodeUninstallCmd      = [aDecoder decodeObjectForKey:kPCPackageVersionNodesUninstallCmd];
        
        self.startScript           = [aDecoder decodeObjectForKey:kPCPackageVersionStartScript];
        self.stopScript            = [aDecoder decodeObjectForKey:kPCPackageVersionStopScript];
        self.cmdScript             = [aDecoder decodeObjectForKey:kPCPackageVersionCmdScript];

        self.processCheck          = [aDecoder decodeObjectForKey:kPCPackageVersionProcessCheck];
        
    }
    return self;
}
@end

@implementation PCPackageMeta

+ (void (^)(NSURLSessionDataTask *, id))metaPackageParserWithBlock:(void (^)(NSArray<PCPackageMeta *> *packages, NSError *error))metaListBlock {
    
#if 0
    void (^parser)(NSURLSessionDataTask*, id) = ^(NSURLSessionDataTask* __unused task , id JSON) {
        
    };
#endif
    
    return [^void(NSURLSessionDataTask* __unused task , id JSON) {
        NSMutableArray<PCPackageMeta *> *metaList = [NSMutableArray arrayWithCapacity:0];
        NSArray *parray = (NSArray *)JSON;
        for (NSDictionary *pi in parray){
            
            NSString *pMeta = [pi objectForKey:kPCPackageMetaVersion];
            NSString *pName = [pi objectForKey:kPCPackageName];
            NSString *pDesc = [pi objectForKey:kDescription];
            NSArray  *pFami = [pi objectForKey:kPCPackageFamily];
            
            NSArray  *pVers = [pi objectForKey:kPCPackageVersions];
            for (NSDictionary *version in pVers){
                
                NSString *vNum   = [version objectForKey:kPCPackageVersionNumber];
                NSString *vDesc  = [version objectForKey:kDescription];
                
                NSArray *vModes = [version objectForKey:kPCPackageVersionModes];
                for (NSDictionary *mode in vModes){
                    
                    PCPackageMeta *meta = [[PCPackageMeta alloc] init];
                    
                    meta.clusterRelation        = @"";
                    
                    meta.metaVersion            = pMeta;
                    meta.packageName            = pName;
                    meta.family                 = pFami;
                    meta.packageDescription     = [NSString stringWithFormat:@"%@ %@ %@", pDesc, vDesc, [mode objectForKey:kDescription]];
                    
                    meta.packageId              = [mode objectForKey:kPCPackageVersionId];
                    meta.version                = vNum;
                    meta.modeType               = [mode objectForKey:kPCPackageVersionModesType];
                    meta.ports                  = [mode objectForKey:kPCPackageVersionPort];
                    meta.libraryDependencies    = [mode objectForKey:kPCPackageVersionLibraryDep];
                    meta.bigPkgDependencies     = [mode objectForKey:kPCPackageVersionBigpkgDep];
                    
                    meta.masterInstallPath      = [mode objectForKey:kPCPackageVersionMasterInstallPath];
                    meta.secondaryInstallPath   = [mode objectForKey:kPCPackageVersionSecondaryInstallPath];
                    meta.nodeInstallPath        = [mode objectForKey:kPCPackageVersionNodeInstallPath];
                    
                    meta.masterCompletePath     = [mode objectForKey:kPCPackageVersionMasterCompletePath];
                    meta.secondaryCompletePath  = [mode objectForKey:kPCPackageVersionSecondaryCompletePath];
                    meta.nodeCompletePath       = [mode objectForKey:kPCPackageVersionNodeCompletePath];
                    
                    meta.masterDownloadPath     = [mode objectForKey:kPCPackageVersionMasterDownload];
                    meta.secondaryDownloadPath  = [mode objectForKey:kPCPackageVersionSecondaryDownload];
                    meta.nodeDownloadPath       = [mode objectForKey:kPCPackageVersionNodeDownload];
                    
                    meta.masterResetCmd         = [mode objectForKey:kPCPackageVersionMasterResetCmd];
                    meta.secondaryResetCmd      = [mode objectForKey:kPCPackageVersionSecondaryResetCmd];
                    meta.nodeResetCmd           = [mode objectForKey:kPCPackageVersionNodesResetCmd];
                    
                    meta.masterUninstallCmd     = [mode objectForKey:kPCPackageVersionMasterUninstallCmd];
                    meta.secondaryUninstallCmd  = [mode objectForKey:kPCPackageVersionSecondaryUninstallCmd];
                    meta.nodeUninstallCmd       = [mode objectForKey:kPCPackageVersionNodesUninstallCmd];
                    
                    meta.startScript            = [mode objectForKey:kPCPackageVersionStartScript];
                    meta.stopScript             = [mode objectForKey:kPCPackageVersionStopScript];
                    meta.cmdScript              = [mode objectForKey:kPCPackageVersionCmdScript];
                    
                    meta.processCheck           = [mode objectForKey:kPCPackageVersionProcessCheck];
                    
                    [metaList addObject:meta];
                }
            }
        }

        if (metaListBlock != nil) {
            metaListBlock(metaList, nil);
        }
    } copy];
    
}


#pragma mark - Class Methods
+ (NSURLSessionDataTask *)metaPackageListWithBlock:(void (^)(NSArray<PCPackageMeta *> *packages, NSError *error))block {
    return [[PCFormulaClient sharedGithubRawFileClient]
            GET:@"meta/packages.json"
            parameters:nil
            success:[PCPackageMeta metaPackageParserWithBlock:block]
            failure:^(NSURLSessionDataTask *__unused task, NSError *error) {
                if (block) {
                    block([NSArray array], error);
                }
            }];
}

+ (NSURLSessionDataTask *)WIPPackageListWithBlock:(void (^)(NSArray<PCPackageMeta *> *packages, NSError *error))block {
    return [[PCFormulaClient sharedWIPRawFileClient]
            GET:@"meta/packages.json"
            parameters:nil
            success:[PCPackageMeta metaPackageParserWithBlock:block]
            failure:^(NSURLSessionDataTask *__unused task, NSError *error) {
                if (block) {
                    block([NSArray array], error);
                }
            }];
}

/*
 * check file list with following commands
 * curl -i https://api.github.com/repos/pocketcluster/formulas/contents/hadoop/2-4-0/datanode/cluster
 */
+ (NSURLSessionDataTask *)packageFileListOn:(NSString *)aPath WithBlock:(void (^)(NSArray<NSString *> *fileList, NSError *error))block {

    return [[PCFormulaClient sharedGithubAPIClient]
            GET:aPath
            parameters:nil
            success:^(NSURLSessionDataTask * __unused task, id JSON){
                
                NSMutableArray<NSString *> *fl = [NSMutableArray arrayWithCapacity:0];
                NSArray *farray = (NSArray *)JSON;
                
                for(NSDictionary *fDict in farray) {
                    [fl addObject:[fDict objectForKey:kGithubRawFileLinkURL]];
                }
                
                if (block) {
                    block(fl, nil);
                }
            } failure:^(NSURLSessionDataTask *__unused task, NSError *error) {
                if (block) {
                    block([NSArray array], error);
                }
            }];
}


-(NSString *)description {
    return [NSString stringWithFormat:@"%@ - %@",[super description], [self packageDescription]];
}

+ (BOOL)makeIntermediateDirectories:(NSString *)aPath {
    NSError *error = nil;
    BOOL isDirectory;
    if([[NSFileManager defaultManager] fileExistsAtPath:aPath isDirectory:&isDirectory]){
        return YES;
    }

    BOOL result = [[NSFileManager defaultManager] createDirectoryAtPath:aPath withIntermediateDirectories:YES attributes:nil error:&error];
    if(!result || error){
        Log(@"Error: Create folder failed %@ %@", aPath, [error debugDescription]);
        return NO;
    }
    return YES;
}

+ (void) downloadFileFromURL:(NSString *)URL
                    basePath:(NSString *)aBasePath
                  completion:(void (^)(NSString *URL, NSURL *filePath))completionBlock
                     onError:(void (^)(NSString *URL, NSError *error))errorBlock {
    
    __block NSString *fileName = [[URL componentsSeparatedByString:@"/"] lastObject];
    
    //Start the download
    [[[PCFormulaClient sharedDownloadManager]
      downloadTaskWithRequest:[NSURLRequest requestWithURL:[NSURL URLWithString:URL]]
      progress:nil
      destination:^NSURL *(NSURL *targetPath, NSURLResponse *response) {
          return [NSURL URLWithString:[NSString stringWithFormat:@"file://%@/%@", aBasePath, fileName]];
      } completionHandler:^(NSURLResponse *response, NSURL *filePath, NSError *error) {
          if (!error) {
              //If there's no error, return the completion block
              completionBlock(URL, filePath);
          } else {
              //Otherwise return the error block
              errorBlock(URL, error);
          }
      }] resume];    
}

@end
