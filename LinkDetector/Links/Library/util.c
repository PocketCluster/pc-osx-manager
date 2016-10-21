//
//  util.c
//  LinkDetector
//
//  Created by Almighty Kim on 21/10/2016.
//  Copyright © 2016 PocketCluster.io. All rights reserved.
//

#include "util.h"
#include <string.h>
#include <stdlib.h>

const char*
copy_string(const char* str_src) {
    if (str_src == NULL || strlen(str_src) == 0) {
        return NULL;
    }
    
    size_t str_len = strlen(str_src);
    char* str_dst = (char*) malloc (sizeof(char) * str_len);
    memcpy (str_dst, str_src, str_len);
    return (const char*) str_dst;
}
