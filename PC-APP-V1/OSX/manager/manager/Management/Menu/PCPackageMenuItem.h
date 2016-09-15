//
//  PCPackageMenuItem.h
//  manager
//
//  Created by Almighty Kim on 11/11/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "PCPackageMeta.h"

@interface PCPackageMenuItem : NSObject
@property (nonatomic, strong, readonly) NSMenuItem *packageItem;

- (instancetype)initWithMetaPackage:(PCPackageMeta *)aMetaPackage;
- (void)destoryMenuItem;
- (void)refreshProcStatus;
@end
