//
//  libvboxcom.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/12/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __LIBVBOXCOM_H__
#define __LIBVBOXCOM_H__

#include <stdbool.h>

typedef void* VBoxGlue;

typedef enum VBGlueResult {
    VBGlue_Ok      = 0,
    VBGlue_Fail
} VBGlueResult;


#pragma mark init & close
VBGlueResult
NewVBoxGlue(VBoxGlue* glue);

VBGlueResult
CloseVBoxGlue(VBoxGlue glue);


#pragma mark app & api version
extern unsigned int VBoxAppVersion(void);
extern unsigned int VBoxApiVersion(void);


#pragma mark system properties
// ** you are resposible for clearing the acquired name string **
VBGlueResult
VBoxHostSearchNetworkInterfaceByName(VBoxGlue glue, const char* queryName, char** fullNameFound);

VBGlueResult
VBoxHostGetMaxGuestCpuCount(VBoxGlue glue, unsigned int* cpuCount);

VBGlueResult
VBoxHostGetMaxGuestMemSize(VBoxGlue glue, unsigned int* memSize);


#pragma mark machine status
typedef enum VBGlueMachineState {
    VBGlueMachine_Illegal       = 0,
    VBGlueMachine_PoweredOff    = 1,
    VBGlueMachine_Saved         = 2,
    VBGlueMachine_Aborted       = 4,
    VBGlueMachine_Running       = 5,
    VBGlueMachine_Paused        = 6,
    VBGlueMachine_Stuck         = 7,
    VBGlueMachine_Starting      = 10,
    VBGlueMachine_Stopping      = 11,
} VBGlueMachineState;

VBGlueMachineState
VBoxMachineGetCurrentState(VBoxGlue glue);

VBGlueResult
VBoxMachineIsSettingChanged(VBoxGlue glue, bool* isMachineChanged);


#pragma mark find, create, build, & release, destroy machine
VBGlueResult
VBoxMachineFindByNameOrID(VBoxGlue glue, const char* machineName);

VBGlueResult
VBoxMachineCreateByName(VBoxGlue glue, const char* baseFolder, const char* machineName);

// option created by this function does not handle deallocation.
// make sure to dealloc it once done
typedef struct VBoxSharedFolder {
    char*    SharedDirName;
    char*    SharedDirPath;
} VBoxSharedFolder;

typedef struct VBoxBuildOption {
    int                CpuCount;
    int                MemSize;
    const char*        HostInterface;
    const char*        BootImagePath;
    const char*        HddImagePath;
    VBoxSharedFolder** Sharedfolders;
    int                SFoldersCount;
} VBoxBuildOption;

VBoxBuildOption*
VBoxMakeBuildOption(int cpu, int mem, const char* host, const char* boot, const char* hdd, VBoxSharedFolder** sfolders, int sflen);

VBGlueResult
VBoxMachineBuildWithOption(VBoxGlue glue, VBoxBuildOption* option);

VBGlueResult
VBoxMachineModifyWithOption(VBoxGlue glue, VBoxBuildOption* option);

VBGlueResult
VBoxMachineDiscardSettings(VBoxGlue glue);

VBGlueResult
VBoxMachineRelease(VBoxGlue glue);

VBGlueResult
VBoxMachineDestory(VBoxGlue glue);


#pragma mark start & stop machine
VBGlueResult
VBoxMachineHeadlessStart(VBoxGlue glue);

VBGlueResult
VBoxMachineForceDown(VBoxGlue glue);

VBGlueResult
VBoxMachineAcpiDown(VBoxGlue glue);


#pragma mark utils
VBGlueResult
VBoxTestErrorMessage(VBoxGlue glue);

const char*
VBoxGetErrorMessage(VBoxGlue glue);

const char*
VBoxGetSettingFilePath(VBoxGlue glue);

const char*
VBoxGetMachineID(VBoxGlue glue);

#endif /* __LIBVBOXCOM_H__ */
