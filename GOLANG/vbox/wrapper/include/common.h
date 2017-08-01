//
//  common.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __COMMON_H__
#define __COMMON_H__

HRESULT VboxFAILED(HRESULT result);

// Wrappers declared in vbox.c
HRESULT VboxArrayOutFree(void* array);

void VboxUtf8Free(char* cstring);

#endif /* COMMON_H */
