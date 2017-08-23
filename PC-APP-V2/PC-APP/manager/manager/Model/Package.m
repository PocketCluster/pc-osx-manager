//
//  PCMeta.m
//  AFNETPROTO
//
//  Created by Almighty Kim on 11/4/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "Package.h"

NSString * const kPackageDescription = @"description";
NSString * const kPackageID          = @"package-id";
NSString * const kInstalled          = @"installed";

@interface Package()
@property (nonatomic, readwrite) BOOL installed;
@property (nonatomic, strong, readwrite) NSString *packageDescription;
@property (nonatomic, strong, readwrite) NSString *packageID;
@end

@implementation Package

+ (NSArray<Package *> *)packagesFromList:(NSArray<NSDictionary *> *)aList {
    if (aList == nil || [aList count] == 0) {
        return nil;
    }
    
    NSMutableArray *packages = [NSMutableArray<Package *> new];
    
    for (NSDictionary *dict in aList) {
        Package *pkg = [Package new];
        pkg.packageDescription = [dict objectForKey:kPackageDescription];
        pkg.packageID = [dict objectForKey:kPackageID];
        pkg.installed = [[dict objectForKey:kInstalled] boolValue];
        [packages addObject:pkg];
    }

    return packages;
}
@end
