#include "../include/common.h"
#include "../include/host.h"

HRESULT VboxGetHostFromVirtualbox(IVirtualBox *cVbox, IHost** host) {
    return IVirtualBox_get_Host(cVbox, host);
}

HRESULT VboxHostSearchNetworkInterfacesFromList(IHost* cHost, const char* cQueryName, char** cNameFound) {

    size_t queryNameLength = strlen(cQueryName) + 1; // query name + :
    PRUint32 wSize                      = 0;
    PRBool isNameFound                  = PR_FALSE;
    HRESULT result;

    IHostNetworkInterface** interfaces  = NULL;
    IHostNetworkInterface* anIface      = NULL;
    BSTR wFullname                      = NULL;
    char* cFullname                     = NULL;
    char* cNameToCompare                = NULL;

    result = cHost->lpVtbl->GetNetworkInterfaces(cHost, &wSize, &interfaces);
    if ( FAILED(result) ) {
        return result;
    }
    if ( wSize == 0 ) {
        return NS_ERROR_NOT_AVAILABLE;
    }
    
    // build name to compare ( query name + ':\0' )
    cNameToCompare = (char *)calloc(queryNameLength + 1, sizeof(char));
    sprintf(cNameToCompare, "%s:", cQueryName);
    
    // iterate through to find out the target interface by name
    for ( int i = 0; i < wSize; i++ ) {
        anIface = interfaces[i];

        // get the full name of an interface
        result = IHostNetworkInterface_get_Name(anIface, &wFullname);
        if ( FAILED(result) ) {
            return result;
        }
        
        // get UTF8 string & free it asap
        result = g_pVBoxFuncs->pfnUtf16ToUtf8(wFullname, &cFullname);
        g_pVBoxFuncs->pfnUtf16Free(wFullname);
        
        // compare the interface name
        if ( strncmp(cNameToCompare, cFullname, queryNameLength) == 0 ) {
            *cNameFound = cFullname;
            isNameFound = PR_TRUE;
            break;
        }
        
        // no match. free UTF8 C string
        free(cFullname);
    }

    // free comparison buffer
    free(cNameToCompare);
    // free interface arrays
    VboxArrayOutFree(interfaces);

    if (isNameFound != PR_TRUE) {
        return NS_ERROR_NOT_AVAILABLE;
    }
    
    return NS_OK;
}

HRESULT VboxHostFindNetworkInterfaceByID(IHost* cHost, const char* cIfID, IHostNetworkInterface** interfaces) {
    HRESULT result;
    
    BSTR wIfID;
    result = g_pVBoxFuncs->pfnUtf8ToUtf16(cIfID, &wIfID);
    if (FAILED(result)) {
        return result;
    }

    result = IHost_FindHostNetworkInterfaceById(cHost, wIfID, interfaces);
    
    g_pVBoxFuncs->pfnUtf16Free(wIfID);
    return result;
}

// this needs exact name to match
HRESULT VboxHostFindNetworkInterfaceByName(IHost* cHost, const char* cName, IHostNetworkInterface** interfaces) {
    HRESULT result;
    
    BSTR wName;
    result = g_pVBoxFuncs->pfnUtf8ToUtf16(cName, &wName);
    if (FAILED(result)) {
        return result;
    }
    
    result = IHost_FindHostNetworkInterfaceByName(cHost, wName, interfaces);
    
    g_pVBoxFuncs->pfnUtf16Free(wName);
    return result;
}

HRESULT VboxGetHostNetworkInterfaceName(IHostNetworkInterface* cIface, char** cFullname) {
    HRESULT result;
    
    BSTR wFullname;
    result = IHostNetworkInterface_get_Name(cIface, &wFullname);
    if (FAILED(result)) {
        return result;
    }

    result = g_pVBoxFuncs->pfnUtf16ToUtf8(wFullname, cFullname);
    g_pVBoxFuncs->pfnUtf16Free(wFullname);
    return result;
}

#if defined(VBOX_BIND_BUG_FIXED)
HRESULT VboxHostFindNetworkInterfaceByType(IHost* cHost, HostNetworkInterfaceType type, PRUint32* cSize, IHostNetworkInterface** interfaces) {
    // invalid argument numbers prevent this to be compiled
    return IHost_FindHostNetworkInterfacesOfType(cHost, type, cSize, interfaces);
}
#endif