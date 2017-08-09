//
//  GlobalMacro.h
//
//  Created by Almighty Kim on 8/23/11.
//  Copyright 2011 __MyCompanyName__. All rights reserved.
//

#ifdef DEBUG
	#define Log(args...)	         NSLog(@"%@",[NSString stringWithFormat:args])
	#define Assert(cond,desc...)	 NSAssert(cond, @"%@", [NSString stringWithFormat: desc])
#else
	#define Log(args...)
	#define Assert(cond,desc...)
#endif

#define IS_EMPTY_STRING(__POINTER) \
    (__POINTER == nil || \
    __POINTER == (NSString *)[NSNull null] || \
    ![__POINTER isKindOfClass:[NSString class]] || \
    ![__POINTER length])
