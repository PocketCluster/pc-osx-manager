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

typedef void** VOID_DPTR;

typedef enum VBGlueResult {
    VBGlue_Ok      = 0,
    VBGlue_Fail
} VBGlueResult;

typedef void VBoxGlue;


#pragma mark init & close
VBGlueResult
NewVBoxGlue(VBoxGlue** glue);

VBGlueResult
CloseVBoxGlue(VBoxGlue* glue);


#pragma mark app & api version
extern unsigned int VBoxAppVersion(void);
extern unsigned int VBoxApiVersion(void);


#pragma mark machine meta
VBGlueResult
VBoxGetMachineID(VBoxGlue* glue, char** machine_id);

VBGlueResult
VBoxGetMachineSettingFilePath(VBoxGlue* glue, char** setting_file_path);

bool
VBoxIsMachineSettingChanged(VBoxGlue* glue);


#pragma mark find, build & destroy machine
VBGlueResult
VBoxFindMachineByNameOrID(VBoxGlue* glue, const char* machine_name);

VBGlueResult
VBoxCreateMachineByName(VBoxGlue* glue, const char* machine_name, char** setting_file_path);

VBGlueResult
VBoxReleaseMachine(VBoxGlue* glue);




#pragma mark utils
VBGlueResult
VBoxTestErrorMessage(VBoxGlue* glue);

const char*
VBoxGetErrorMessage(VBoxGlue* glue);

#endif /* __LIBVBOXCOM_H__ */
