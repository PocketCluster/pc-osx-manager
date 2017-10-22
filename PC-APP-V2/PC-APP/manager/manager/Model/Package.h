//
//  Package
//  PocketCluster
//
//  Created by Almighty Kim on 11/4/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

extern NSString * const kPackageDescription;
extern NSString * const kPackageID;
extern NSString * const kInstalled;

@interface Package : NSObject

@property (nonatomic, readonly) BOOL installed;
@property (nonatomic, readonly) NSString *packageDescription;
@property (nonatomic, readonly) NSString *packageID;

+ (NSArray<Package *> *)packagesFromList:(NSArray<NSDictionary *> *)aList;

// update all package status except `package id` with a package w/ same id.
- (void) updateWithPackage:(Package *)newSamePackage;
@end
