//
//  PCPkgProc.h
//  manager
//
//  Created by Almighty Kim on 11/10/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCPackageMeta.h"

@interface PCPkgProc : NSObject
@property (nonatomic, weak) PCPackageMeta *package;
@property (nonatomic, readonly) BOOL isAlive;
-(void)refreshProcessStatus;

@end
