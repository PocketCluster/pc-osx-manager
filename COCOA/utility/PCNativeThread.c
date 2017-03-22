//
//  PCNativeThread.c
//  static-core
//
//  Created by Almighty Kim on 3/22/17.
//  Copyright © 2017 PocketCluster. All rights reserved.
//

#include <stdlib.h>
#include <pthread.h>
#include "PCNativeThread.h"

unsigned long long
PCNativeThreadID() {
    unsigned long long id;
    if (pthread_threadid_np(pthread_self(), &id)) {
        abort();
    }
    return id;
}