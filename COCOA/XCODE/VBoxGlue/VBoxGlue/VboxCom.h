//
//  VboxCom.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/12/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#import <Foundation/Foundation.h>

extern NSString * const kVboxComDomain;

@interface VboxCom : NSObject
@property (nonatomic, readonly) NSString *baseFolder;

-(instancetype)initWithError:(__autoreleasing NSError **)error;

-(BOOL)checkAppVersion;


-(NSString *)retrieveMachineId:(__autoreleasing NSError **)error;

-(BOOL)isMachineSettingChanged:(__autoreleasing NSError **)error;


-(BOOL)openSession:(__autoreleasing NSError **)error;

-(void)closeSession:(__autoreleasing NSError **)error;


-(BOOL)acquireMachineByNameOrID:(NSString *)machineNameOrID error:(__autoreleasing NSError **)error;

-(BOOL)createMachineWithName:(NSString *)machineName error:(__autoreleasing NSError **)error;

-(void)releaseMachine:(__autoreleasing NSError **)error;


-(BOOL)buildMachineWithCPU:(int)cpuCount
                memorySize:(int)memorySize
             hostInterface:(NSString *)hostInterface
          sharedFolderName:(NSString *)sharedFolderName
          sharedFolderPath:(NSString *)sharedFolderPath
             bootImagePath:(NSString *)bootImagePath
              hardDiskSize:(int)hardDiskSize
                  progress:(void (^)(int progress))buildProgress
                     error:(__autoreleasing NSError **)error;

-(BOOL)destoryCurrentMachine:(__autoreleasing NSError **)error;
@end