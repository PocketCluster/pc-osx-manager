//
//  medium.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __MEDIUM_H__
#define __MEDIUM_H__

HRESULT VboxCreateHardDisk(IVirtualBox* cbox, char* cformat, char* clocation, DeviceType deviceType, AccessMode accessMode, IMedium** cmedium);

HRESULT VboxOpenMedium(IVirtualBox* cbox, char* clocation, DeviceType cdeviceType, AccessMode caccessType, PRBool cforceNewUuid, IMedium** cmedium);

HRESULT VboxMediumCreateBaseStorage(IMedium* cmedium, PRInt64 size, PRUint32 variantCount, PRUint32* cvariant, IProgress** cprogress);

HRESULT VboxMediumDeleteStorage(IMedium* cmedium, IProgress** cprogress);

HRESULT VboxMediumClose(IMedium* cmedium);

HRESULT VboxGetMediumLocation(IMedium* cmedium, char** clocation);

HRESULT VboxGetMediumState(IMedium* cmedium, PRUint32* cstate);

HRESULT VboxGetMediumSize(IMedium* cmedium, PRInt64* csize);

HRESULT VboxIMediumRelease(IMedium* cmedium);

#endif /* medium_h */
