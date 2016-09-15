//
//  PCPackageManager.h
//  manager
//
//  Created by Almighty Kim on 11/10/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCPackageMeta.h"

@interface PCPackageManager : NSObject
@property (nonatomic, strong, readonly) NSMutableArray<PCPackageMeta *> *installedPackage;

+ (instancetype)sharedManager;

- (void)addInstalledPackage:(PCPackageMeta *)aPackage;
- (void)removeInstalledPackage:(PCPackageMeta *)aPackage;

- (void)loadInstalledPackage;
- (void)saveInstalledPackage;

#ifdef DEBUG
- (void)clearInstalledPackageMeta;
#endif
@end
