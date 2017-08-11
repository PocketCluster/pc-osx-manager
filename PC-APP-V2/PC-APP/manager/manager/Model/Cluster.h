//
//  RaspberryCluster.h
//  manager
//
//  Created by Almighty Kim on 11/2/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "Node.h"

@interface Cluster : NSObject
@property (nonatomic, strong, readonly) NSString* ClusterID;
@property (nonatomic, strong, readonly) NSString* ClusterUUID;
@property (nonatomic, strong, readonly) NSString* ClusterDomain;
@property (nonatomic, strong, readonly) NSString* UserMadeName;

@property (nonatomic, strong, readonly) NSArray* Nodes;
@property (nonatomic, strong, readonly) NSArray* Packages;

@end
