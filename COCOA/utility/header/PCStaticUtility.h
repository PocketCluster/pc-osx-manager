//
//  PCStaticUtility.h
//  LinkDetector
//
//  Created by Almighty Kim on 21/10/2016.
//  Copyright © 2016 PocketCluster.io. All rights reserved.
//

#ifndef __PCSTATICUTILITY_H__
#define __PCSTATICUTILITY_H__

#include <stdio.h>
#include <CoreFoundation/CoreFoundation.h>

const char*
copy_string(const char* str_src);

const char*
CFStringCopyToCString(CFStringRef string);

#endif //__PCSTATICUTILITY_H__
