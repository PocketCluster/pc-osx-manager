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

NSString * const kDescription                       = @"description";
NSString * const kPCPackageName                     = @"name";
NSString * const kPCPackageFamily                   = @"family";
NSString * const kPCPackageVersions                 = @"versions";
NSString * const kPCPackageVersionNumber            = @"ver-num";
NSString * const kPCPackageVersionModes             = @"modes";
NSString * const kPCPackageVersionModesType         = @"type";
NSString * const kPCPackageVersionPort              = @"ports";
NSString * const kPCPackageVersionLibraryDep        = @"dep-lib";
NSString * const kPCPackageVersionBigpkgDep         = @"dep-bigpkg";
NSString * const kPCPackageVersionMasterPath        = @"masters-path";
NSString * const kPCPackageVersionNodesPath         = @"nodes-path";

NSString * const kGithubRawFileLinkURL              = @"download_url";

@interface PCPackageMeta()
@property (nonatomic, strong, readwrite) NSString *packageName;
@property (nonatomic, strong, readwrite) NSArray<NSString*> *family;
@property (nonatomic, strong, readwrite) NSString *packageDescription;
@property (nonatomic, strong, readwrite) NSString *version;
@property (nonatomic, strong, readwrite) NSString *modeType;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *ports;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *libraryDependencies;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *bigPkgDependencies;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *masterFilePath;
@property (nonatomic, strong, readwrite) NSArray<NSString *> *nodeFilePath;
@end


#pragma mark - PCPACKAGEMETA NSCODING
@implementation PCPackageMeta (NSCoding)
- (void)encodeWithCoder:(NSCoder *)aCoder {
    [aCoder encodeObject:self.packageName           forKey:kPCPackageName];
    [aCoder encodeObject:self.family                forKey:kPCPackageFamily];
    [aCoder encodeObject:self.packageDescription    forKey:kDescription];
    [aCoder encodeObject:self.version               forKey:kPCPackageVersionNumber];
    [aCoder encodeObject:self.modeType              forKey:kPCPackageVersionModesType];
    [aCoder encodeObject:self.ports                 forKey:kPCPackageVersionPort];
    [aCoder encodeObject:self.libraryDependencies   forKey:kPCPackageVersionLibraryDep];
    [aCoder encodeObject:self.bigPkgDependencies    forKey:kPCPackageVersionBigpkgDep];
    [aCoder encodeObject:self.masterFilePath        forKey:kPCPackageVersionMasterPath];
    [aCoder encodeObject:self.nodeFilePath          forKey:kPCPackageVersionNodesPath];
}

- (instancetype)initWithCoder:(NSCoder *)aDecoder {
    self = [super init];
    if (self) {

        self.packageName           = [aDecoder decodeObjectForKey:kPCPackageName];
        self.family                = [aDecoder decodeObjectForKey:kPCPackageFamily];
        self.packageDescription    = [aDecoder decodeObjectForKey:kDescription];
        self.version               = [aDecoder decodeObjectForKey:kPCPackageVersionNumber];
        self.modeType              = [aDecoder decodeObjectForKey:kPCPackageVersionModesType];
        self.ports                 = [aDecoder decodeObjectForKey:kPCPackageVersionPort];
        self.libraryDependencies   = [aDecoder decodeObjectForKey:kPCPackageVersionLibraryDep];
        self.bigPkgDependencies    = [aDecoder decodeObjectForKey:kPCPackageVersionBigpkgDep];
        self.masterFilePath        = [aDecoder decodeObjectForKey:kPCPackageVersionMasterPath];
        self.nodeFilePath          = [aDecoder decodeObjectForKey:kPCPackageVersionNodesPath];
        
    }
    return self;
}
@end

@implementation PCPackageMeta
#pragma mark - Class Methods
+ (NSURLSessionDataTask *)metaPackageListWithBlock:(void (^)(NSArray<PCPackageMeta *> *packages, NSError *error))block {
    return [[PCFormulaClient sharedGithubRawFileClient]
            GET:@"meta/packages.json"
            parameters:nil
            success:^(NSURLSessionDataTask * __unused task, id JSON){

                NSMutableArray<PCPackageMeta *> *metaList = [NSMutableArray arrayWithCapacity:0];
                NSArray *parray = (NSArray *)JSON;
                for (NSDictionary *pi in parray){
                    
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
                            meta.packageName        = pName;
                            meta.family             = pFami;
                            meta.packageDescription = [NSString stringWithFormat:@"%@ %@ %@", pDesc, vDesc, [mode objectForKey:kDescription]];
                            meta.version            = vNum;
                            meta.modeType           = [mode objectForKey:kPCPackageVersionModesType];
                            meta.ports              = [mode objectForKey:kPCPackageVersionPort];
                            meta.libraryDependencies= [mode objectForKey:kPCPackageVersionLibraryDep];
                            meta.bigPkgDependencies = [mode objectForKey:kPCPackageVersionBigpkgDep];
                            meta.masterFilePath     = [mode objectForKey:kPCPackageVersionMasterPath];
                            meta.nodeFilePath       = [mode objectForKey:kPCPackageVersionNodesPath];

                            [metaList addObject:meta];
                        }
                    }
                }
                
                if (block) {
                    block(metaList, nil);
                }
            } failure:^(NSURLSessionDataTask *__unused task, NSError *error) {
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
    NSString *targetPath = [NSString stringWithFormat:@"%@/%@",kPOCKET_CLUSTER_SALT_STATE_PATH, aPath];
    NSError *error = nil;
    BOOL isDirectory;
    if([[NSFileManager defaultManager] fileExistsAtPath:targetPath isDirectory:&isDirectory]){
        return YES;
    }
    
    BOOL result = [[NSFileManager defaultManager] createDirectoryAtPath:targetPath withIntermediateDirectories:YES attributes:nil error:&error];
    if(!result || error){
        Log(@"Error: Create folder failed %@ %@", targetPath, [error debugDescription]);
        return NO;
    }
    return YES;
}

+ (void) downloadFileFromURL:(NSString *)URL
                    basePath:(NSString *)aBasePath
                  completion:(void (^)(NSURL *filePath))completionBlock
                     onError:(void (^)(NSError *error))errorBlock {
    
    //Configuring the session manager
    __block AFURLSessionManager *manager = [PCFormulaClient sharedDownloadManager];
    __block NSString *fileName = [[URL componentsSeparatedByString:@"/"] lastObject];
    
    //Start the download
    [[manager
      downloadTaskWithRequest:[NSURLRequest requestWithURL:[NSURL URLWithString:URL]]
      progress:nil
      destination:^NSURL *(NSURL *targetPath, NSURLResponse *response) {
          return [NSURL URLWithString:[NSString stringWithFormat:@"file://%@/%@", aBasePath, fileName]];
      } completionHandler:^(NSURLResponse *response, NSURL *filePath, NSError *error) {
          if (!error) {
              //If there's no error, return the completion block
              completionBlock(filePath);
          } else {
              //Otherwise return the error block
              errorBlock(error);
          }
      }] resume];    
}

@end
