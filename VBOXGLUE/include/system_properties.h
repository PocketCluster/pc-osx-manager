//
//  system_properties.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __SYSTEM_PROPERTIES_H__
#define __SYSTEM_PROPERTIES_H__

HRESULT VboxGetSystemProperties(IVirtualBox* cbox, ISystemProperties **cprops);

HRESULT VboxGetSystemPropertiesMaxGuestRAM(ISystemProperties* cprops, ULONG *cmaxRam);

HRESULT VboxGetSystemPropertiesMaxGuestVRAM(ISystemProperties* cprops, ULONG *cmaxVram);

HRESULT VboxGetSystemPropertiesMaxGuestCpuCount(ISystemProperties* cprops, ULONG *cmaxCpus);

HRESULT VboxISystemPropertiesRelease(ISystemProperties* cprops);

#endif /* system_properties_h */
