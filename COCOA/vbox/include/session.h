//
//  session.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __SESSION_H__
#define __SESSION_H__

HRESULT VboxUnlockMachine(ISession* csession);

HRESULT VboxISessionRelease(ISession* csession);

HRESULT VboxGetSessionConsole(ISession* csession, IConsole** cconsole);

HRESULT VboxGetSessionMachine(ISession* csession, IMachine** cmachine);

HRESULT VboxGetSessionState(ISession* csession, PRUint32* cstate);

HRESULT VboxGetSessionType(ISession* csession, PRUint32* ctype);

HRESULT VboxGetSession(IVirtualBoxClient* client, ISession** csession);

HRESULT VboxLockMachine(IMachine* cmachine, ISession* csession, PRUint32 clock);

#endif /* session_h */
