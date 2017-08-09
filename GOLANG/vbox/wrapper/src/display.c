#include "VBoxCAPIGlue.h"
#include "common.h"

HRESULT VboxIDisplayRelease(IDisplay* cdisplay) {
    return IDisplay_Release(cdisplay);
}
