#include "VBoxCAPIGlue.h"
#include "common.h"

HRESULT VboxGetBiosSettingsLogoFadeIn(IBIOSSettings* csettings, PRBool* clogoFadeIn) {
    return IBIOSSettings_GetLogoFadeIn(csettings, clogoFadeIn);
}

HRESULT VboxSetBiosSettingsLogoFadeIn(IBIOSSettings* csettings, PRBool clogoFadeIn) {
    return IBIOSSettings_SetLogoFadeIn(csettings, clogoFadeIn);
}

HRESULT VboxGetBiosSettingsLogoFadeOut(IBIOSSettings* csettings, PRBool* clogoFadeOut) {
    return IBIOSSettings_GetLogoFadeOut(csettings, clogoFadeOut);
}

HRESULT VboxSetBiosSettingsLogoFadeOut(IBIOSSettings* csettings, PRBool clogoFadeOut) {
    return IBIOSSettings_SetLogoFadeOut(csettings, clogoFadeOut);
}

HRESULT VboxGetBiosSettingsBootMenuMode(IBIOSSettings* csettings, PRUint32* cmenuMode) {
    return IBIOSSettings_GetBootMenuMode(csettings, cmenuMode);
}

HRESULT VboxSetBiosSettingsBootMenuMode(IBIOSSettings* csettings, PRUint32 cmenuMode) {
    return IBIOSSettings_SetBootMenuMode(csettings, cmenuMode);
}

HRESULT VboxIBiosSettingsRelease(IBIOSSettings* csettings) {
    HRESULT result = IBIOSSettings_Release(csettings);
    csettings = NULL;
    return result;
}

HRESULT VboxGetMachineBIOSSettings(IMachine* cmachine, IBIOSSettings** csettings) {
    return IMachine_GetBIOSSettings(cmachine, csettings);
}
