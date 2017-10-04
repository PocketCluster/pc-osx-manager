//
//  consol.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __CONSOL_H__
#define __CONSOL_H__

HRESULT VboxGetConsoleDisplay(IConsole* cconsole, IDisplay** cdisplay);

HRESULT VboxGetConsoleKeyboard(IConsole* cconsole, IKeyboard** ckeyboard);

HRESULT VboxGetConsoleMouse(IConsole* cconsole, IMouse** cmouse);

HRESULT VboxGetConsoleMachine(IConsole* cconsole, IMachine** cmachine);

HRESULT VboxConsoleForcePowerDown(IConsole* cconsole, IProgress** cprogress);

HRESULT VboxConsoleAcpiPowerDown(IConsole* cconsole);

HRESULT VboxConsoleRelease(IConsole* cconsole);

#endif /* consol_h */
