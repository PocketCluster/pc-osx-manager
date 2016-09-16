//
//  VboxCom.m
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/12/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#import "VboxCom.h"
#import "GlobalMacro.h"
#import "libvboxcom.h"

#define ERROR_MSG_SIZE 512

NSString * const kLibVboxComErrorDomain = @"io.pocketcluster.osx.libvboxcom";

NSString * const STORAGE_CONTROLLER_NAME = @"SATA";

@interface VboxCom()
+(NSError *)_buildErrorMessageWithResult:(VBRESULT)result errorMessage:(char*)errorMessage suggestion:(NSString *)suggestion;
+(NSError *)_buildErrorMessageWithResult:(VBRESULT)result failureReason:(NSString *)failureReason suggestion:(NSString * const)suggestion;
+(NSError *)_uninitVboxError;

@end

@implementation VboxCom {
    void* _vbox_client;        // vbox client
    void* _virtualbox;         // vbox machine
    void* _vbox_session;       // vbox session
    void* _vbox_machine;       // vbox machine
    char* _error_message;
    char* _base_folder;        // vbox setting file path for example "/Users/almightykim/VirtualBox VMs/pc-master/pc-master.vbox"
    VBRESULT _init_status;     // vbox init status
    
}
@dynamic baseFolder;

-(instancetype)init {
    @throw([NSException exceptionWithName:@"InvalidMethodInvocation" reason:@"One should not use -[VboxCom init]" userInfo:nil]);
}

-(instancetype)initWithError:(NSError **)error {
    self = [super init];
    if (self) {
        _error_message  = (char *)malloc(sizeof(char) * ERROR_MSG_SIZE);
        _vbox_client    = NULL;
        _vbox_session   = NULL;
        _virtualbox     = NULL;
        _vbox_machine   = NULL;
        _base_folder    = NULL;
        memset(_error_message, ERROR_MSG_SIZE, sizeof(char));
        
        _init_status = vbox_init(_error_message);
        if (_init_status != GOOD) {
            *error = [VboxCom _buildErrorMessageWithResult:_init_status errorMessage:_error_message suggestion:@"Please restart computer or upgrade Virtualbox to the latest version."];
            strcpy(_error_message, "");
        }
        
    }
    return self;
}

-(void)dealloc {
    free(_error_message);
    vbox_term();
    _init_status = FATAL;
}

#pragma mark - UTIL
+(NSError *)_uninitVboxError {
    return [VboxCom _buildErrorMessageWithResult:FATAL failureReason:@"VboxCom is not perperly initialized" suggestion:@"Please restart computer to properly start"];
}

+(NSError *)_buildErrorMessageWithResult:(VBRESULT)result errorMessage:(char*)errorMessage suggestion:(NSString * const)suggestion {
    return [VboxCom
            _buildErrorMessageWithResult:result
            failureReason:[NSString stringWithCString:errorMessage encoding:NSASCIIStringEncoding]
            suggestion:suggestion];
}

+(NSError *)_buildErrorMessageWithResult:(VBRESULT)result failureReason:(NSString *)failureReason suggestion:(NSString * const)suggestion {
    Assert(!IS_EMPTY_STRING(failureReason) ,@"Failure reason message cannot be empty");
    Assert(!IS_EMPTY_STRING(suggestion)    ,@"Error suggestion message cannot be empty");
    
    NSError *error = nil;
    NSDictionary *userInfo = nil;
    switch (result) {
        case INFO: {
            userInfo = @{NSLocalizedFailureReasonErrorKey:failureReason};
            break;
        }
        case FATAL: {
            userInfo = @{NSLocalizedFailureReasonErrorKey:failureReason, NSLocalizedRecoverySuggestionErrorKey:suggestion};
            break;
        }
        default:
            break;
    }
    
    if (result != GOOD) {
        error = [NSError
                 errorWithDomain:kLibVboxComErrorDomain
                 code:result
                 userInfo:userInfo];
    }
    
    return error;
}

-(NSString *)baseFolder {
    if (_base_folder == NULL || strlen(_base_folder) == 0) {
        return nil;
    }
    return [NSString stringWithCString:_base_folder encoding:NSASCIIStringEncoding];
}

-(BOOL)checkAppVersion {
    return 5000000 <= vbox_app_version();
}

#pragma mark - MACHINE PROPERTY

-(NSString *)retrieveMachineId:(__autoreleasing NSError **)error {
    NSString* machineId = nil;
    if (_init_status != GOOD) {
        *error = [VboxCom _uninitVboxError];
        return machineId;
    }
    
    char* vbox_machine_id = NULL;
    VBRESULT result = vbox_machine_getid(&_vbox_machine, &vbox_machine_id, _error_message);
    if (result != GOOD || vbox_machine_id == NULL) {
        *error = [VboxCom _buildErrorMessageWithResult:result errorMessage:_error_message suggestion:@"Please restart computer to create new master node."];
        strcpy(_error_message, "");
        return machineId;
    }
    
    machineId = [NSString stringWithCString:vbox_machine_id encoding:NSASCIIStringEncoding];
    VboxUtf8Free(vbox_machine_id);
    return machineId;
}

-(BOOL)isMachineSettingChanged:(__autoreleasing NSError **)error {
    if (_init_status != GOOD) {
        *error = [VboxCom _uninitVboxError];
        return NO;
    }
    if (_vbox_machine == NULL) {
        *error = [VboxCom _buildErrorMessageWithResult:FATAL errorMessage:_error_message suggestion:@"Please open or create a machine before checking"];
        return NO;
    }

    int result = vbox_machine_is_setting_changed(&_vbox_machine, &_vbox_session, _error_message);
    if (strlen(_error_message) != 0) {
        *error = [VboxCom _buildErrorMessageWithResult:FATAL errorMessage:_error_message suggestion:@"Please restart to property check machine setting modification"];
    }

    return (BOOL)result;
}



#pragma mark - INIT & CLOSE

/**
 XPCOM session initializer. This is to
 1) check if the version of Virtualbox is what we are looking for,
 2) check if there is an error happens when initialize session variables
 
 @param error   an error message
 */
-(BOOL)openSession:(__autoreleasing NSError **)error{
    if (_init_status != GOOD) {
        *error = [VboxCom _uninitVboxError];
        return NO;
    }
    
    VBRESULT result = vbox_session_init(&_vbox_client, &_vbox_session, &_virtualbox, _error_message);
    if (result != GOOD) {
        *error = [VboxCom _buildErrorMessageWithResult:result errorMessage:_error_message suggestion:@"Please restart computer to create new master node."];
        strcpy(_error_message, "");
        return NO;
    }
    return YES;
}

/**
 XPCOM session deinitializer. This is to simply close session variables.
 
 @param error   an error message
 */
-(void)closeSession:(__autoreleasing NSError **)error {
    if (_init_status != GOOD) {
        *error = [VboxCom _uninitVboxError];
        return;
    }

    VBRESULT result = vbox_session_close(&_vbox_client, &_vbox_session, &_virtualbox, _error_message);
    if (result != GOOD) {
        *error = [VboxCom _buildErrorMessageWithResult:result errorMessage:_error_message suggestion:@"THIS IS INFORMATIVE"];
        strcpy(_error_message, "");
    }
    _vbox_client    = NULL;
    _vbox_session   = NULL;
    _virtualbox     = NULL;
    _vbox_machine   = NULL;
}

#pragma mark - FIND, BUILD & DESTROY MACHINE
/**
 Find a machine by name
 
 @param machineName    a name to find a machine
 */
-(BOOL)acquireMachineByNameOrID:(NSString *)machineNameOrID error:(__autoreleasing NSError **)error {
    Assert(!IS_EMPTY_STRING(machineNameOrID), @"machine name cannot be empty.");
    if (_init_status != GOOD) {
        *error = [VboxCom _uninitVboxError];
        return NO;
    }
    if (_vbox_machine != NULL) {
        *error = [VboxCom _buildErrorMessageWithResult:FATAL errorMessage:_error_message suggestion:@"Please release opened machine before find one"];
        return NO;
    }
    
    VBRESULT result = vbox_machine_find(&_virtualbox, &_vbox_machine, [machineNameOrID UTF8String], _error_message);
    if (result != GOOD) {
        *error = [VboxCom _buildErrorMessageWithResult:result errorMessage:_error_message suggestion:@"Please restart computer to create new master node."];
        strcpy(_error_message, "");

        if (result == FATAL || _vbox_machine == NULL)
            return NO;
    }

    result = vbox_machine_setting_path(&_vbox_machine, &_base_folder, _error_message);
    if (result != GOOD) {
        *error = [VboxCom _buildErrorMessageWithResult:result errorMessage:_error_message suggestion:@"Please restart computer to property check master node."];
        strcpy(_error_message, "");
        
        if (result == FATAL)
            return NO;
    }    
    return YES;
}

-(BOOL)createMachineWithName:(NSString *)machineName error:(__autoreleasing NSError **)error {
    Assert(!IS_EMPTY_STRING(machineName), @"machine name cannot be empty.");
    if (_init_status != GOOD) {
        *error = [VboxCom _uninitVboxError];
        return NO;
    }
    
    VBRESULT result = vbox_machine_create(&_virtualbox, &_vbox_machine, [machineName UTF8String], &_base_folder, _error_message);
    if (result != GOOD) {
        *error = [VboxCom _buildErrorMessageWithResult:result errorMessage:_error_message suggestion:@"Please restart computer to create new master node."];
        strcpy(_error_message, "");
        return NO;
    }
    return YES;
}

-(void)releaseMachine:(__autoreleasing NSError **)error {
    Assert(_vbox_machine != NULL, @"vbox machine cannot be null");
    if (_init_status != GOOD) {
        *error = [VboxCom _uninitVboxError];
        return;
    }
    if (_vbox_machine == NULL) {
        *error = [VboxCom _buildErrorMessageWithResult:FATAL errorMessage:_error_message suggestion:@"Please open or create a machine before releasing"];
        return;
    }

    VBRESULT result = vbox_machine_release(&_vbox_machine, _base_folder, _error_message);
    if (result != GOOD) {
        *error = [VboxCom _buildErrorMessageWithResult:result errorMessage:_error_message suggestion:@"THIS IS INFORMATIVE"];
        strcpy(_error_message, "");
    }
}

#pragma mark - BUILD & DESTORY MACHINE

-(BOOL)buildMachineWithCPU:(int)cpuCount
                memorySize:(int)memorySize
             hostInterface:(NSString *)hostInterface
          sharedFolderName:(NSString *)sharedFolderName
          sharedFolderPath:(NSString *)sharedFolderPath
             bootImagePath:(NSString *)bootImagePath
              hardDiskSize:(int)hardDiskSize
                  progress:(void (^)(int progress))buildProgress
                     error:(__autoreleasing NSError **)error
{
    Assert(0 < cpuCount && cpuCount <= 16            ,@"CPU count cannot be smaller than 0 or greater than 16");
    Assert(2048 <= memorySize && memorySize <= 16384 ,@"Memory size cannot be smaller than 2G and greater than 16G");
    Assert(!IS_EMPTY_STRING(hostInterface)           ,@"Host interface cannot be empty");
    Assert(!IS_EMPTY_STRING(sharedFolderName)        ,@"Shared folder name cannot be empty");
    Assert(!IS_EMPTY_STRING(sharedFolderPath)        ,@"Shared folder path cannot be empty");
    Assert(!IS_EMPTY_STRING(bootImagePath)           ,@"Boot image path cannot be empty");
    //TODO: check hard disk size
    if (_init_status != GOOD) {
        *error = [VboxCom _uninitVboxError];
        return NO;
    }

    
    //TODO : always find before building
    
    // build basic machine with bios & motherboard settings
    VBRESULT result = vbox_machine_build(&_virtualbox, &_vbox_machine, cpuCount, memorySize, _error_message);
    if (result != GOOD) {
        *error = [VboxCom _buildErrorMessageWithResult:result errorMessage:_error_message suggestion:@"Please restart computer to create new master node."];
        strcpy(_error_message, "");

        if (result == FATAL)
            return NO;
    }

    // add bridged network
    result = vbox_machine_add_bridged_network(&_vbox_machine, &_vbox_session, [hostInterface UTF8String], _error_message);
    if (result != GOOD) {
        *error = [VboxCom _buildErrorMessageWithResult:result errorMessage:_error_message suggestion:@"Please restart computer to create new master node."];
        strcpy(_error_message, "");
        
        if (result == FATAL)
            return NO;
    }

    // add shared folder
    result = vbox_machine_add_shared_folder(&_vbox_machine, &_vbox_session, [sharedFolderName UTF8String], [sharedFolderPath UTF8String], _error_message);
    if (result != GOOD) {
        *error = [VboxCom _buildErrorMessageWithResult:result errorMessage:_error_message suggestion:@"Please restart computer to create new master node."];
        strcpy(_error_message, "");
        
        if (result == FATAL)
            return NO;
    }

    // add storage controller
    result = vbox_machine_add_storage_controller(&_vbox_machine, &_vbox_session, [STORAGE_CONTROLLER_NAME UTF8String], _error_message);
    if (result != GOOD) {
        *error = [VboxCom _buildErrorMessageWithResult:result errorMessage:_error_message suggestion:@"Please restart computer to create new master node."];
        strcpy(_error_message, "");
        
        if (result == FATAL)
            return NO;
    }

    // add boot image
    result = vbox_machine_add_boot_image(&_virtualbox, &_vbox_machine, &_vbox_session, [STORAGE_CONTROLLER_NAME UTF8String], [bootImagePath UTF8String], _error_message);
    if (result != GOOD) {
        *error = [VboxCom _buildErrorMessageWithResult:result errorMessage:_error_message suggestion:@"Please restart computer to create new master node."];
        strcpy(_error_message, "");
        
        if (result == FATAL)
            return NO;
    }
 
    // add hard drive
    // TODO : we need to replace this with pure c function in libvboxcom.h/c.
    NSString *hddMediumPath = [[NSString stringWithCString:_base_folder encoding:NSASCIIStringEncoding]
                               stringByReplacingOccurrencesOfString:@"pc-master.vbox"
                               withString:@"pc-master-hdd.vmdk"];

    result = vbox_machine_add_hard_disk(&_virtualbox, &_vbox_machine, &_vbox_session, [STORAGE_CONTROLLER_NAME UTF8String], [hddMediumPath UTF8String], hardDiskSize,
                                        ^(int progress, int done) {
                                   
                                            if (buildProgress != nil) {
                                                buildProgress(progress);
                                            }

                                        }, _error_message);
    if (result != GOOD) {
        *error = [VboxCom _buildErrorMessageWithResult:result errorMessage:_error_message suggestion:@"Please restart computer to create new master node."];
        strcpy(_error_message, "");
        
        if (result == FATAL)
            return NO;
    }
    
    return YES;
}

-(BOOL)destoryCurrentMachine:(__autoreleasing NSError **)error {
    // TODO : when you destory a machine, you also need to destory everything include folder, residue and/or whatnots
    if (_init_status != GOOD) {
        *error = [VboxCom _uninitVboxError];
        return NO;
    }
    if (_vbox_machine == NULL) {
        *error = [VboxCom _buildErrorMessageWithResult:FATAL errorMessage:_error_message suggestion:@"Please open or create a machine before checking"];
        return NO;
    }
    
    BOOL dont_remove_boot_image = NO;
    VBRESULT result = vbox_machine_destroy(&_vbox_machine, _base_folder, [STORAGE_CONTROLLER_NAME UTF8String], (int)dont_remove_boot_image, ^(int progress, int done) {
    }, _error_message);
    
    if (result != GOOD) {
        *error = [VboxCom _buildErrorMessageWithResult:result errorMessage:_error_message suggestion:@"Please restart computer to delete machine."];
        strcpy(_error_message, "");
        
        if (result == FATAL)
            return NO;
    }
    return YES;
}

@end
