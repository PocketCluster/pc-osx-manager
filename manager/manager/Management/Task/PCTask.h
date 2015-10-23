//
//  PCTask.h
//  manager
//
//  Created by Almighty Kim on 10/23/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface PCTask : NSObject

@property (strong, nonatomic) id target;
@property (strong, nonatomic) NSString *taskCommand;
@property (strong, nonatomic) NSString *taskAction;
@property (strong, nonatomic) NSTask *task;
@property (strong, nonatomic) NSString *taskUUID;

- (void)launchTask;
- (void)cancelTask;
@end
