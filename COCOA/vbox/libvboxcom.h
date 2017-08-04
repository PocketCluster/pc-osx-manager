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

typedef struct VBoxBuildOption {
    int            CpuCount;
    int            MemSize;
    const char*    HostInterface;
    const char*    BootImagePath;
    const char*    HddImagePath;
    const char*    SharedDirPath;
    const char*    SharedDirName;
} VBoxBuildOption;


#pragma mark init & close
VBGlueResult
NewVBoxGlue(VBoxGlue* glue);

VBGlueResult
CloseVBoxGlue(VBoxGlue glue);


#pragma mark app & api version
extern unsigned int VBoxAppVersion(void);
extern unsigned int VBoxApiVersion(void);

#pragma mark host network interface
// ** you are resposible for clearing the acquired name string **
VBGlueResult
VBoxSearchHostNetworkInterfaceByName(VBoxGlue glue, const char* queryName, char** fullNameFound);


#pragma mark machine status
VBGlueResult
VBoxIsMachineSettingChanged(VBoxGlue glue, bool* isMachineChanged);


#pragma mark find, create, & release machine
VBGlueResult
VBoxFindMachineByNameOrID(VBoxGlue glue, const char* machineName);

VBGlueResult
VBoxCreateMachineByName(VBoxGlue glue, const char* baseFolder, const char* machineName);

VBGlueResult
VBoxReleaseMachine(VBoxGlue glue);


#pragma mark build & destroy machine
// option created by this function does not handle deallocation.
// make sure to dealloc it once done
VBoxBuildOption*
VBoxMakeBuildOption(int cpu, int mem, const char* host, const char* boot, const char* hdd, const char* spath, const char* sname);

VBGlueResult
VBoxBuildMachine(VBoxGlue glue, VBoxBuildOption* option);

VBGlueResult
VBoxDestoryMachine(VBoxGlue glue);


#if 0
#pragma mark start & stop machine
VBGlueResult
VBoxStartMachine(VBoxGlue glue);

VBGlueResult
VBoxStopMachine(VBoxGlue glue);
#endif


#pragma mark utils
VBGlueResult
VBoxTestErrorMessage(VBoxGlue glue);

const char*
VBoxGetErrorMessage(VBoxGlue glue);

const char*
VboxGetSettingFilePath(VBoxGlue glue);

const char*
VboxGetMachineID(VBoxGlue glue);

#endif /* __LIBVBOXCOM_H__ */
