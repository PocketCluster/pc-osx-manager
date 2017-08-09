//
//  machine.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __MACHINE_H__
#define __MACHINE_H__

HRESULT VboxGetMachineName(IMachine* cmachine, char** cname);

HRESULT VboxGetMachineOSTypeId(IMachine* cmachine, char** cosTypeId);

HRESULT VboxGetMachineSettingsFilePath(IMachine* cmachine, char** cpath);

HRESULT VboxGetMachineMemorySize(IMachine* cmachine, PRUint32* cram);

HRESULT VboxSetMachineMemorySize(IMachine* cmachine, PRUint32 cram);

HRESULT VboxGetMachineVRAMSize(IMachine* cmachine, PRUint32* cvram);

HRESULT VboxSetMachineVRAMSize(IMachine* cmachine, PRUint32 cvram);

HRESULT VboxGetMachinePointingHIDType(IMachine* cmachine, PRUint32* ctype);

HRESULT VboxSetMachinePointingHIDType(IMachine* cmachine, PRUint32 ctype);

HRESULT VboxGetMachineSettingsModified(IMachine* cmachine, PRBool* cmodified);

HRESULT VboxMachineSaveSettings(IMachine* cmachine);

HRESULT VboxMachineUnregister(IMachine* cmachine, PRUint32 cleanupMode, IMedium*** cmedia, ULONG* mediaCount);

HRESULT VboxMachineDeleteConfig(IMachine* cmachine, PRUint32 mediaCount, IMedium** cmedia, IProgress** cprogress);

HRESULT VboxMachineAttachDevice(IMachine* cmachine, const char* cname, PRInt32 cport, PRInt32 cdevice, PRUint32 ctype, IMedium* cmedium);

HRESULT VboxMachineUnmountMedium(IMachine* cmachine, const char* cname, PRInt32 cport, PRInt32 cdevice, PRBool cforce);

HRESULT VboxMachineGetMedium(IMachine* cmachine, char* cname, PRInt32 cport, PRInt32 cdevice, IMedium** cmedium);

HRESULT VboxMachineCreateSharedFolder(IMachine *cmachine, const char* cName, const char* cHostPath, BOOL cWritable, BOOL cAutomount);

HRESULT VboxMachineGetNetworkAdapter(IMachine *cmachine, PRUint32 slotNumber, INetworkAdapter **cAdapter);

HRESULT VboxIMachineRelease(IMachine* cmachine);

HRESULT VboxCreateMachine(IVirtualBox* cbox, char* cSettingsFile, const char* cname, char* cosTypeId, char* cflags, IMachine** cmachine);

HRESULT VboxMachineGetID(IMachine* cmachine, char** cMachineId);

HRESULT VboxFindMachine(IVirtualBox* cbox, const char* cnameOrId, IMachine** cmachine);

HRESULT VboxGetMachines(IVirtualBox* cbox, IMachine*** cmachines, ULONG* machineCount);

HRESULT VboxRegisterMachine(IVirtualBox* cbox, IMachine* cmachine);

#pragma mark State/Launch

HRESULT VboxMachineGetState(IMachine* cmachine, MachineState *state);

HRESULT VboxMachineLaunchVMProcess(IMachine* cmachine, ISession* csession, char* cuiType, char* cenvironment, IProgress** cprogress);

#endif /* machine_h */
