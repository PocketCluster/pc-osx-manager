#ifndef __BIOS_SETTINGS_H__
#define __BIOS_SETTINGS_H__

HRESULT VboxGetBiosSettingsLogoFadeIn(IBIOSSettings* csettings, PRBool* clogoFadeIn);

HRESULT VboxSetBiosSettingsLogoFadeIn(IBIOSSettings* csettings, PRBool clogoFadeIn);

HRESULT VboxGetBiosSettingsLogoFadeOut(IBIOSSettings* csettings, PRBool* clogoFadeOut);

HRESULT VboxSetBiosSettingsLogoFadeOut(IBIOSSettings* csettings, PRBool clogoFadeOut);

HRESULT VboxGetBiosSettingsBootMenuMode(IBIOSSettings* csettings, PRUint32* cmenuMode);

HRESULT VboxSetBiosSettingsBootMenuMode(IBIOSSettings* csettings, PRUint32 cmenuMode);

HRESULT VboxIBiosSettingsRelease(IBIOSSettings* csettings);

HRESULT VboxGetMachineBIOSSettings(IMachine* cmachine, IBIOSSettings** csettings);;

#endif /* bios_settings_h */
