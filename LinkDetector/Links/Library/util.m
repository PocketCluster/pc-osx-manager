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
#include <CoreFoundation/CoreFoundation.h>

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


const char*
CFStringCopyToCString(CFStringRef string) {
    CFIndex length = CFStringGetLength(string);
    if (string == NULL || length == 0) {
        return NULL;
    }
    // due to termination character, we need +1 space at the end
    char *str_buffer = (char *) malloc (sizeof(char) * (length + 1));
    CFStringGetCString(string, str_buffer, length + 1, kCFStringEncodingUTF8);
    return (const char*)str_buffer;
}
