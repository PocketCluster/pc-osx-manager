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

HRESULT VboxMachineAttachDevice(IMachine* cmachine, char* cname, PRInt32 cport, PRInt32 cdevice, PRUint32 ctype, IMedium* cmedium);

HRESULT VboxMachineUnmountMedium(IMachine* cmachine, char* cname, PRInt32 cport, PRInt32 cdevice, PRBool cforce);

HRESULT VboxMachineGetMedium(IMachine* cmachine, char* cname, PRInt32 cport, PRInt32 cdevice, IMedium** cmedium);

HRESULT VboxMachineGetNetworkAdapter(IMachine *cmachine, PRUint32 slotNumber, INetworkAdapter **cAdapter);

HRESULT VboxMachineLaunchVMProcess(IMachine* cmachine, ISession* csession, char* cuiType, char* cenvironment, IProgress** cprogress);

HRESULT VboxIMachineRelease(IMachine* cmachine);

HRESULT VboxCreateMachine(IVirtualBox* cbox, char* cSettingsFile, char* cname, char* cosTypeId, char* cflags, IMachine** cmachine);

HRESULT VboxFindMachine(IVirtualBox* cbox, char* cnameOrId, IMachine** cmachine);

HRESULT VboxGetMachines(IVirtualBox* cbox, IMachine*** cmachines, ULONG* machineCount);

HRESULT VboxRegisterMachine(IVirtualBox* cbox, IMachine* cmachine);

#endif /* machine_h */
