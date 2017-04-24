//
//  ConstPackage.h
//  PocketCluster
//
//  Created by Almighty Kim on 8/23/16.
//  Copyright 2011 __MyCompanyName__. All rights reserved.
//

#ifndef __GLOBAL_MACRO_H__
#define __GLOBAL_MACRO_H__

#ifdef DEBUG
	#define Log(args...)	NSLog(@"%@",[NSString stringWithFormat:args])
	#define Assert(cond,desc...)	NSAssert(cond, @"%@", [NSString stringWithFormat: desc])
	#define SAFE_DEALLOC_CHECK(__POINTER) { Log(@"%@ dealloc",self);}
	#define ASSURE_VIEW_REMOVAL(__VIEW_POINTER) {[__VIEW_POINTER removeFromSuperview]; __VIEW_POINTER = nil;}
	//#define SAFE_DEALLOC_CHECK(__POINTER) { [super dealloc]; }
#else
	#define Log(args...)
	#define Assert(cond,desc...)
	#define SAFE_DEALLOC_CHECK(__POINTER)
	#define ASSURE_VIEW_REMOVAL(__VIEW_POINTER) {[__VIEW_POINTER removeFromSuperview]; __VIEW_POINTER = nil;}
#endif

#define ASSURE_DEALLOC(__POINTER) { __POINTER = nil; }

#define CHECK_DELEGATE_EXECUTION(__POINTER,__PROTOCOL,__SELECTOR) \
    ((__POINTER != nil) && [__POINTER conformsToProtocol:__PROTOCOL] && [__POINTER respondsToSelector:__SELECTOR])

#define CHECK_FUNCTION_EXEC() Log(@"%s",__PRETTY_FUNCTION__)

#endif //__GLOBAL_MACRO_H__