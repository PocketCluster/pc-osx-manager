/* Created by "go tool cgo" - DO NOT EDIT. */

/* package main */

/* Start of preamble from import "C" comments.  */


#line 3 "/Users/almightykim/Workspace/pc-osx-manager/PC-APP-V2/PC-CORE/static/native_debug.go"




#line 4 "/Users/almightykim/Workspace/pc-osx-manager/PC-APP-V2/PC-CORE/static/native_feedresponse.go"




#include "PCResponseHandle.h"


#line 9 "/Users/almightykim/Workspace/pc-osx-manager/PC-APP-V2/PC-CORE/static/native_lifecycle.go"




#include <pthread.h>
#include "PCLifeCycle.h"
#include "PCNativeThread.h"



#line 4 "/Users/almightykim/Workspace/pc-osx-manager/PC-APP-V2/PC-CORE/static/native_netnoti.go"




#include "SCNetworkTypes.h"
#include "PCInterfaceTypes.h"



#line 3 "/Users/almightykim/Workspace/pc-osx-manager/PC-APP-V2/PC-CORE/static/native_route.go"





/* End of preamble from import "C" comments.  */


/* Start of boilerplate cgo prologue.  */

#ifndef GO_CGO_PROLOGUE_H
#define GO_CGO_PROLOGUE_H

typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef __SIZE_TYPE__ GoUintptr;
typedef float GoFloat32;
typedef double GoFloat64;
typedef float _Complex GoComplex64;
typedef double _Complex GoComplex128;

/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of GoInt.
*/
typedef char _check_for_64_bit_pointer_matching_GoInt[sizeof(void*)==64/8 ? 1:-1];

typedef struct { const char *p; GoInt n; } GoString;
typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;

#endif

/* End of boilerplate cgo prologue.  */

#ifdef __cplusplus
extern "C" {
#endif


extern void OpsCmdBaseServiceStart();

extern void OpsCmdBaseServiceStop();

extern void OpsCmdStorageStart();

extern void OpsCmdStorageStop();

extern void OpsCmdTeleportRootAdd();

extern void OpsCmdTeleportUserAdd();

extern void OpsCmdDebug();

extern void StartResponseFeed();

extern void StopResponseFeed();

extern void lifecycleDead();

extern void lifecycleAlive();

extern void lifecycleVisible();

extern void lifecycleFocused();

extern void lifecycleAwaken();

extern void lifecycleSleep();

extern void crashEmergentExit();

extern void engineDebugOutput(int p0);

extern void NetworkChangeNotificationInterface(PCNetworkInterface** p0, unsigned int p1);

extern void NetworkChangeNotificationGateway(SCNIGateway** p0, unsigned int p1);

extern void RouteRequestGet(char* p0);

extern void RouteRequestPost(char* p0, char* p1);

extern void RouteRequestPut(char* p0, char* p1);

extern void RouteRequestDelete(char* p0);

#ifdef __cplusplus
}
#endif
