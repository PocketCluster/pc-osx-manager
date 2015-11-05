//
//  PCTask.h
//  manager
//
//  Created by Almighty Kim on 10/23/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <Foundation/Foundation.h>

@class PCTask;

@protocol PCTaskDelegate <NSObject>

@required
-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask;
-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler;
-(BOOL)task:(PCTask *)aPCTask isOutputClosed:(id<PCTaskDelegate>)aDelegate;

@optional
-(void)setTarget:(id)aTarget;
-(void)setTaskCommand:(NSString *)aTaskCommand;
-(void)setTaskAction:(NSString *)aTaskAction;
-(void)setTask:(NSTask *)aTask;
-(void)setSudoCommand:(BOOL)aSudoCommand;
@end

@interface PCTask : NSObject

@property (weak, nonatomic) id<PCTaskDelegate> delegate;

@property (readonly) BOOL isRunning;

@property (strong, nonatomic) id target;
@property (strong, nonatomic) NSString *taskCommand;
@property (strong, nonatomic) NSString *taskAction;
@property (strong, nonatomic) NSTask *task;
@property (strong, nonatomic) NSString *taskUUID;
@property (nonatomic, getter=isSudoCommand) BOOL sudoCommand;

- (void)launchTask;
- (void)cancelTask;

- (void)runTaskWithProgressBlock:(void (^)(NSString *output))progress done:(void (^)(NSTask *task))done;
@end
