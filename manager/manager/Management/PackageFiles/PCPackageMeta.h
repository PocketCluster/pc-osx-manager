//
//  PCMeta.h
//  AFNETPROTO
//
//  Created by Almighty Kim on 11/4/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

extern NSString * const kPCPackageName;
extern NSString * const kPCPackageFamily;
extern NSString * const kDescription;
extern NSString * const kPCPackageVersions;
extern NSString * const kPCPackageVersionNumber;
extern NSString * const kPCPackageVersionModes;
extern NSString * const kPCPackageVersionModesType;
extern NSString * const kPCPackageVersionMasterPath;
extern NSString * const kPCPackageVersionNodesPath;

@interface PCPackageMeta : NSObject
+ (NSURLSessionDataTask *)metaPackageListWithBlock:(void (^)(NSArray<PCPackageMeta *> *packages, NSError *error))block;
+ (NSURLSessionDataTask *)packageFileListOn:(NSString *)aPath WithBlock:(void (^)(NSArray<NSString *> *fileList, NSError *error))block;

@property (nonatomic, readonly) NSString *packageName;
@property (nonatomic, readonly) NSArray<NSString*> *family;
@property (nonatomic, readonly) NSString *packageDescription;
@property (nonatomic, readonly) NSString *version;
@property (nonatomic, readonly) NSString *modeType;
@property (nonatomic, readonly) NSArray<NSString *> *ports;
@property (nonatomic, readonly) NSArray<NSString *> *libraryDependencies;
@property (nonatomic, readonly) NSArray<NSString *> *bigPkgDependencies;
@property (nonatomic, readonly) NSArray<NSString *> *masterFilePath;
@property (nonatomic, readonly) NSArray<NSString *> *nodeFilePath;

+ (BOOL)makeIntermediateDirectories:(NSString *)aPath;

+ (void)downloadFileFromURL:(NSString *)URL
                   basePath:(NSString *)aBasePath
                 completion:(void (^)(NSURL *filePath))completionBlock
                    onError:(void (^)(NSError *error))errorBlock;
@end

@interface PCPackageMeta (NSCoding) <NSCoding>
@end
