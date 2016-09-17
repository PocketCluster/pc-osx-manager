//
//  common.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __COMMON_H__
#define __COMMON_H__

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <signal.h>
#include <unistd.h>
#include <sys/poll.h>
#include "VBoxCAPIGlue.h"

// Wrappers declared in vbox.c
HRESULT VboxArrayOutFree(void* array);

void VboxUtf8Free(char* cstring);

#endif /* COMMON_H */
