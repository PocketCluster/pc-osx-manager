//
//  PCMeta.m
//  AFNETPROTO
//
//  Created by Almighty Kim on 11/4/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "NullStringChecker.h"
#import "Package.h"

NSString * const kPackageDescription = @"description";
NSString * const kPackageID          = @"package-id";
NSString * const kInstalled          = @"installed";
NSString * const kMenuName           = @"menu-name";
NSString * const kCoreImageSize      = @"core-image-size";
NSString * const kNodeImageSize      = @"node-image-size";
NSString * const kConsoles           = @"consoles";
NSString * const kAddress            = @"address";
NSString * const kName               = @"name";

@interface Console()
- (instancetype) initWithConsoleDict:(NSDictionary *)aConsoleDict;
@property (nonatomic, strong, readwrite) NSString *address;
@property (nonatomic, strong, readwrite) NSString *name;
@end

@implementation Console
- (instancetype) initWithConsoleDict:(NSDictionary *)aConsoleDict {
    self = [super init];
    if (self) {
        self.address = [aConsoleDict objectForKey:kAddress];
        self.name    = [aConsoleDict objectForKey:kName];
    }
    return self;
}
@end

@interface Package()
@property (nonatomic, readwrite) BOOL installed;
@property (nonatomic, strong, readwrite) NSString *packageDescription;
@property (nonatomic, strong, readwrite) NSString *packageID;
@property (nonatomic, strong, readwrite) NSString *menuName;
@property (nonatomic, strong, readwrite) NSString *totalImageSize;
@property (nonatomic, strong, readwrite) NSArray<Console *>* consoles;
@property (nonatomic, readwrite) ExecState execState;
@end

@implementation Package

+ (NSArray<Package *> *)packagesFromList:(NSArray<NSDictionary *> *)aList {
    static const double unit_gigabyte = 1073741824.0;
    static const double unit_megabyte = 1048576.0;

    if (aList == nil || [aList count] == 0) {
        return nil;
    }
    
    NSMutableArray *packages = [NSMutableArray<Package *> new];
    
    for (NSDictionary *dict in aList) {
        Package *pkg = [Package new];
        pkg.packageDescription = [dict objectForKey:kPackageDescription];
        pkg.packageID = [dict objectForKey:kPackageID];
        pkg.menuName  = [dict objectForKey:kMenuName];
        pkg.installed = [[dict objectForKey:kInstalled] boolValue];

        // core image size
        NSString *cSize = [dict objectForKey:kCoreImageSize];
        // node imaeg size
        NSString *nSize = [dict objectForKey:kNodeImageSize];
        // total image size
        if (!(ISNULL_STRING(cSize) || ISNULL_STRING(nSize))) {
            double tBytes = [cSize doubleValue] + [nSize doubleValue];
            NSString *sTotal;
            if (unit_gigabyte < tBytes) {
                sTotal = [NSString stringWithFormat:@"%.1lf GB", (tBytes / unit_gigabyte)];
            } else {
                sTotal = [NSString stringWithFormat:@"%.1lf MB", (tBytes / unit_megabyte)];
            }
            pkg.totalImageSize = sTotal;
        } else {
            pkg.totalImageSize = nil;
        }

        // package consoles
        NSMutableArray<Console *> *pConsoles = [NSMutableArray<Console *> arrayWithCapacity:0];
        for (NSDictionary *c in (NSArray<NSDictionary*>*)[dict objectForKey:kConsoles]) {
            [pConsoles addObject:[[Console alloc] initWithConsoleDict:c]];
        }
        pkg.consoles = pConsoles;

        [packages addObject:pkg];
    }

    return packages;
}

- (void) updateWithPackage:(Package *)aPackage {
    self.packageDescription = aPackage.packageDescription;
    self.menuName           = aPackage.menuName;
    self.installed          = aPackage.installed;

    if (!ISNULL_STRING(aPackage.totalImageSize)) {
        self.totalImageSize = aPackage.totalImageSize;
    }
}

// execution state should only be handled separately
- (void) updateExecState:(ExecState)state {
    self.execState = state;
}
@end
