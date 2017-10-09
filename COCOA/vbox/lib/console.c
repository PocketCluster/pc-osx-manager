#include "../include/common.h"
#include "../include/console.h"

HRESULT VboxGetConsoleDisplay(IConsole* cconsole, IDisplay** cdisplay) {
    return IConsole_GetDisplay(cconsole, cdisplay);
}

HRESULT VboxGetConsoleKeyboard(IConsole* cconsole, IKeyboard** ckeyboard) {
    return IConsole_GetKeyboard(cconsole, ckeyboard);
}

HRESULT VboxGetConsoleMouse(IConsole* cconsole, IMouse** cmouse) {
    return IConsole_GetMouse(cconsole, cmouse);
}

HRESULT VboxGetConsoleMachine(IConsole* cconsole, IMachine** cmachine) {
    return IConsole_GetMachine(cconsole, cmachine);
}

HRESULT VboxConsoleForcePowerDown(IConsole* cconsole, IProgress** cprogress) {
    return IConsole_PowerDown(cconsole, cprogress);
}

HRESULT VboxConsoleAcpiPowerDown(IConsole* cconsole) {
    return IConsole_PowerButton(cconsole);
}

HRESULT VboxConsoleRelease(IConsole* cconsole) {
    return IConsole_Release(cconsole);
}
