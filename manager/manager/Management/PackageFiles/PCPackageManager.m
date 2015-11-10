//
//  PCPackageManager.m
//  manager
//
//  Created by Almighty Kim on 11/10/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "SynthesizeSingleton.h"
#import "PCPackageManager.h"
#import "PCConstants.h"

@interface PCPackageManager()
@property (nonatomic, strong) NSMutableArray<PCPackageMeta *> *installedPackage;
@end

@implementation PCPackageManager
SYNTHESIZE_SINGLETON_FOR_CLASS_WITH_ACCESSOR(PCPackageManager, sharedManager);

- (instancetype) init {
    
    self = [super init];
    if(self){
        self.installedPackage = [NSMutableArray<PCPackageMeta *> arrayWithCapacity:0];
    }
    return self;
}

- (void)addInstalledPackage:(PCPackageMeta *)aPackage onCluster:(NSString *)aCluster {
    
    // check if a package with same id and cluster relation exists
    @synchronized(self) {
        [self.installedPackage addObject:aPackage];
    }
}

- (void)removeInstalledPackage:(PCPackageMeta *)aPackage {
    @synchronized(self) {
        [self.installedPackage removeObject:aPackage];
    }
}

- (void)loadInstalledPackage {
    @synchronized(self) {
        [self.installedPackage removeAllObjects];
        id data = [[NSUserDefaults standardUserDefaults] dataForKey:kPCInstalledPackageCollection];
        if(data) {
            NSArray *saved = (NSArray *)[NSKeyedUnarchiver unarchiveObjectWithData:data];
            [self.installedPackage addObjectsFromArray:saved];
        }
    }
}

- (void)saveInstalledPackage {
    @synchronized(self) {
        NSMutableArray *rpis = [self installedPackage];
        if(rpis != nil && [self installedPackage] != 0) {
            NSData *data = [NSKeyedArchiver archivedDataWithRootObject:rpis];
            if (data){
                [[NSUserDefaults standardUserDefaults] setObject:data forKey:kPCInstalledPackageCollection];
                [[NSUserDefaults standardUserDefaults] synchronize];
            }
        }
    }
}

@end
