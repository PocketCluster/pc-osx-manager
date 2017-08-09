#include "VBoxCAPIGlue.h"
#include "common.h"

HRESULT VboxGetSystemProperties(IVirtualBox* cbox, ISystemProperties **cprops) {
    return IVirtualBox_GetSystemProperties(cbox, cprops);
}

HRESULT VboxGetSystemPropertiesMaxGuestRAM(ISystemProperties* cprops, ULONG *cmaxRam) {
    return ISystemProperties_GetMaxGuestRAM(cprops, cmaxRam);
}

HRESULT VboxGetSystemPropertiesMaxGuestVRAM(ISystemProperties* cprops, ULONG *cmaxVram) {
    return ISystemProperties_GetMaxGuestVRAM(cprops, cmaxVram);
}

HRESULT VboxGetSystemPropertiesMaxGuestCpuCount(ISystemProperties* cprops, ULONG *cmaxCpus) {
    return ISystemProperties_GetMaxGuestVRAM(cprops, cmaxCpus);
}

HRESULT VboxISystemPropertiesRelease(ISystemProperties* cprops) {
    return ISystemProperties_Release(cprops);
}
