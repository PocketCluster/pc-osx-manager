#include "VBoxCAPIGlue.h"

// NOTE: Including the C file is a sketchy but working method for getting it
//       compiled and linked with the Go source. The C must only be included in
//       one Go file. By convention, this is the file that wraps the
//       ClientInitialize() function.
#include "VBoxCAPIGlue.c"

HRESULT VboxInit() {
    return VBoxCGlueInit();
}

void VboxTerm() {
    VBoxCGlueTerm();
}

HRESULT VboxFAILED(HRESULT result) {
    return FAILED(result);
}

HRESULT VboxArrayOutFree(void* carray) {
    return g_pVBoxFuncs->pfnArrayOutFree(carray);
}

void VboxUtf8Free(char* cstring) {
    g_pVBoxFuncs->pfnUtf8Free(cstring);
    cstring = NULL;
}

unsigned int VboxGetAppVersion() {
    return g_pVBoxFuncs->pfnGetVersion();
}

unsigned int VboxGetApiVersion() {
    return g_pVBoxFuncs->pfnGetAPIVersion();
}

HRESULT VboxClientInitialize(IVirtualBoxClient** client) {
    return g_pVBoxFuncs->pfnClientInitialize(NULL, client);
}

HRESULT VboxClientThreadInitialize() {
    return g_pVBoxFuncs->pfnClientThreadInitialize();
}

HRESULT VboxClientThreadUninitialize() {
    return g_pVBoxFuncs->pfnClientThreadUninitialize();
}

void VboxClientUninitialize() {
    g_pVBoxFuncs->pfnClientUninitialize();
}

HRESULT VboxClientRelease(IVirtualBoxClient* client) {
    HRESULT result = IVirtualBoxClient_Release(client);
    client = NULL;
    return result;
}

HRESULT VboxGetVirtualBox(IVirtualBoxClient* client, IVirtualBox** cbox) {
    return IVirtualBoxClient_GetVirtualBox(client, cbox);
}

HRESULT VboxGetRevision(IVirtualBox* cbox, ULONG* revision) {
    return IVirtualBox_GetRevision(cbox, revision);
}

HRESULT VboxIVirtualBoxRelease(IVirtualBox* cbox) {
    HRESULT result = IVirtualBox_Release(cbox);
    cbox = NULL;
    return result;
}

HRESULT VboxComposeMachineFilename(IVirtualBox* cbox, const char* cname, char* cflags, char* cbaseFolder, char **cpath) {
    BSTR wname;
    HRESULT result = g_pVBoxFuncs->pfnUtf8ToUtf16(cname, &wname);
    if (FAILED(result)) {
        return result;
    }

    BSTR wflags = NULL;
    result = g_pVBoxFuncs->pfnUtf8ToUtf16(cflags, &wflags);
    if (FAILED(result)) {
        g_pVBoxFuncs->pfnUtf16Free(wname);
    }

    BSTR wbaseFolder;
    result = g_pVBoxFuncs->pfnUtf8ToUtf16(cbaseFolder, &wbaseFolder);
    if (FAILED(result)) {
        g_pVBoxFuncs->pfnUtf16Free(wflags);
        g_pVBoxFuncs->pfnUtf16Free(wname);
        return result;
    }

    BSTR wpath = NULL;
    result = IVirtualBox_ComposeMachineFilename(cbox, wname, NULL, wflags, wbaseFolder, &wpath);
    g_pVBoxFuncs->pfnUtf16Free(wbaseFolder);
    g_pVBoxFuncs->pfnUtf16Free(wflags);
    g_pVBoxFuncs->pfnUtf16Free(wname);
    if (FAILED(result)) {
        return result;
    }

    g_pVBoxFuncs->pfnUtf16ToUtf8(wpath, cpath);
    g_pVBoxFuncs->pfnComUnallocString(wpath);
    return result;
}
