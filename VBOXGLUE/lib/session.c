#include "../include/common.h"
#include "../include/session.h"

HRESULT VboxUnlockMachine(ISession* csession) {
    return ISession_UnlockMachine(csession);
}

HRESULT VboxISessionRelease(ISession* csession) {
    HRESULT result = ISession_Release(csession);
    csession = NULL;
    return result;
}

HRESULT VboxGetSessionConsole(ISession* csession, IConsole** cconsole) {
    return ISession_GetConsole(csession, cconsole);
}

HRESULT VboxGetSessionMachine(ISession* csession, IMachine** cmachine) {
    return ISession_GetMachine(csession, cmachine);
}

HRESULT VboxGetSessionState(ISession* csession, PRUint32* cstate) {
    return ISession_GetState(csession, cstate);
}

HRESULT VboxGetSessionType(ISession* csession, PRUint32* ctype) {
    return ISession_GetType(csession, ctype);
}

HRESULT VboxGetSession(IVirtualBoxClient* client, ISession** csession) {
    return IVirtualBoxClient_GetSession(client, csession);
}

HRESULT VboxLockMachine(IMachine* cmachine, ISession* csession, PRUint32 clock) {
    return IMachine_LockMachine(cmachine, csession, clock);
}
