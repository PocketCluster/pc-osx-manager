//
//  vbox.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __VBOX_H__
#define __VBOX_H__

HRESULT VboxInit();

void VboxTerm();

unsigned int VboxGetAppVersion();

unsigned int VboxGetApiVersion();

HRESULT VboxClientInitialize(IVirtualBoxClient** client);

HRESULT VboxClientThreadInitialize();

HRESULT VboxClientThreadUninitialize();

void VboxClientUninitialize();

HRESULT VboxClientRelease(IVirtualBoxClient* client);

HRESULT VboxGetVirtualBox(IVirtualBoxClient* client, IVirtualBox** cbox);

HRESULT VboxGetRevision(IVirtualBox* cbox, ULONG* revision);

HRESULT VboxIVirtualBoxRelease(IVirtualBox* cbox);

HRESULT VboxComposeMachineFilename(IVirtualBox* cbox, const char* cname, char* cflags, char* cbaseFolder, char **cpath);

#endif /* vbox_h */
