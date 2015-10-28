//
//  ConstPackage.h
//  FlyStarProto
//
//  Created by Almighty Kim on 8/23/11.
//  Copyright 2011 __MyCompanyName__. All rights reserved.
//



#ifdef DEBUG
//#define USE_TESTFLIGHT
//#import "LoggerClient.h"

	#define Log(args...)	NSLog(@"%@",[NSString stringWithFormat:args])
	//#define Log(args...)	LogMessageCompat(@"%@",[NSString stringWithFormat:args]);
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

#define CHECK_DELEGATE_EXECUTION(__POINTER,__SELECTOR) \
	((__POINTER != nil) && [__POINTER respondsToSelector:__SELECTOR])


#define CHECK_FUNCTION_EXEC() Log(@"%s",__PRETTY_FUNCTION__)