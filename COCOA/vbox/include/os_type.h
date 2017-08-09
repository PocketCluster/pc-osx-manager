//
//  os_type.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __OS_TYPE_H__
#define __OS_TYPE_H__

HRESULT VboxGetGuestOSTypes(IVirtualBox* cbox, IGuestOSType*** ctypes, ULONG* typeCount);

HRESULT VboxGetGuestOSTypeId(IGuestOSType* ctype, char** cid);

HRESULT VboxIGuestOSTypeRelease(IGuestOSType* ctype);

#endif /* os_type_h */
