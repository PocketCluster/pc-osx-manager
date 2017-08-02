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

typedef enum VBGlueResult {
    VBGlue_Ok      = 0,
    VBGlue_Fail
} VBGlueResult;

typedef void* VBoxGlue;

#pragma mark init & close
VBGlueResult
NewVBoxGlue(VBoxGlue* glue);

VBGlueResult
CloseVBoxGlue(VBoxGlue glue);


#pragma mark app & api version
extern unsigned int VBoxAppVersion(void);
extern unsigned int VBoxApiVersion(void);


#pragma mark machine status
VBGlueResult
VBoxIsMachineSettingChanged(VBoxGlue glue, bool* isMachineChanged);


#pragma mark find, build & destroy machine
VBGlueResult
VBoxFindMachineByNameOrID(VBoxGlue glue, const char* machine_name);

VBGlueResult
VBoxCreateMachineByName(VBoxGlue glue, const char* base_folder, const char* machine_name);

VBGlueResult
VBoxReleaseMachine(VBoxGlue glue);


#pragma mark destroy machine
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
