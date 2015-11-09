//
//  PCMeta.h
//  AFNETPROTO
//
//  Created by Almighty Kim on 11/4/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

extern NSString * const kPCPackageMetaVersion;
extern NSString * const kDescription;
extern NSString * const kPCPackageName;
extern NSString * const kPCPackageFamily;
extern NSString * const kPCPackageVersions;
extern NSString * const kPCPackageVersionNumber;
extern NSString * const kPCPackageVersionModes;
extern NSString * const kPCPackageVersionModesType;

extern NSString * const kPCPackageVersionMasterInstallPath;
extern NSString * const kPCPackageVersionSecondaryInstallPath;
extern NSString * const kPCPackageVersionNodeInstallPath;

extern NSString * const kPCPackageVersionMasterCompletePath;
extern NSString * const kPCPackageVersionSecondaryCompletePath;
extern NSString * const kPCPackageVersionNodeCompletePath;

extern NSString * const kPCPackageVersionMasterDownload;
extern NSString * const kPCPackageVersionSecondaryDownload;
extern NSString * const kPCPackageVersionNodeDownload;

extern NSString * const kPCPackageVersionStartScript;
extern NSString * const kPCPackageVersionStopScript;
extern NSString * const kPCPackageVersionCmdScript;

extern NSString * const kPCPackageVersionProcessCheck;
extern NSString * const kPCPackageVersionResetCmd;
extern NSString * const kPCPackageVersionUninstallCmd;

@interface PCPackageMeta : NSObject

@property (nonatomic, readonly) NSString *metaVersion;
@property (nonatomic, readonly) NSString *packageName;
@property (nonatomic, readonly) NSArray<NSString*> *family;
@property (nonatomic, readonly) NSString *packageDescription;
@property (nonatomic, readonly) NSString *version;
@property (nonatomic, readonly) NSString *modeType;
@property (nonatomic, readonly) NSArray<NSString *> *ports;
@property (nonatomic, readonly) NSArray<NSString *> *libraryDependencies;
@property (nonatomic, readonly) NSArray<NSString *> *bigPkgDependencies;

@property (nonatomic, readonly) NSArray<NSString *> *masterInstallPath;
@property (nonatomic, readonly) NSArray<NSString *> *secondaryInstallPath;
@property (nonatomic, readonly) NSArray<NSString *> *nodeInstallPath;

@property (nonatomic, readonly) NSArray<NSString *> *masterCompletePath;
@property (nonatomic, readonly) NSArray<NSString *> *secondaryCompletePath;
@property (nonatomic, readonly) NSArray<NSString *> *nodeCompletePath;

@property (nonatomic, readonly) NSArray<NSString *> *masterDownloadPath;
@property (nonatomic, readonly) NSArray<NSString *> *secondaryDownloadPath;
@property (nonatomic, readonly) NSArray<NSString *> *nodeDownloadPath;

@property (nonatomic, readonly) NSArray<NSString *> *startScript;
@property (nonatomic, readonly) NSArray<NSString *> *stopScript;
@property (nonatomic, readonly) NSArray<NSString *> *cmdScript;

@property (nonatomic, readonly) NSArray<NSString *> *processCheck;
@property (nonatomic, readonly) NSArray<NSString *> *resetCommands;
@property (nonatomic, readonly) NSArray<NSString *> *uninstallCommands;

+ (NSURLSessionDataTask *)metaPackageListWithBlock:(void (^)(NSArray<PCPackageMeta *> *packages, NSError *error))block;
+ (NSURLSessionDataTask *)packageFileListOn:(NSString *)aPath WithBlock:(void (^)(NSArray<NSString *> *fileList, NSError *error))block;
+ (BOOL)makeIntermediateDirectories:(NSString *)aPath;
+ (void)downloadFileFromURL:(NSString *)URL
                   basePath:(NSString *)aBasePath
                 completion:(void (^)(NSString *URL, NSURL *filePath))completionBlock
                    onError:(void (^)(NSString *URL, NSError *error))errorBlock;
@end

@interface PCPackageMeta (NSCoding) <NSCoding>
@end
