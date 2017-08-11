//
//  RaspberryCluster.m
//  manager
//
//  Created by Almighty Kim on 11/2/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "Cluster.h"
#import "NullStringChecker.h"
#import "PCPackageManager.h"

@interface Cluster()
@property (nonatomic, strong, readwrite) NSString *ClusterID;
@property (nonatomic, strong, readwrite) NSString *ClusterUUID;
@property (nonatomic, strong, readwrite) NSString* ClusterDomain;
@property (nonatomic, strong, readwrite) NSString* UserMadeName;

@property (nonatomic, strong, readwrite) NSMutableArray *Nodes;
@property (nonatomic, strong, readwrite) NSMutableArray *Packages;
@end

@implementation Cluster

@end
