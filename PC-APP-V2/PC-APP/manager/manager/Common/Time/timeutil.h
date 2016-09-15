//
//  timeutil.h
//  OSXGLEssentials
//
//  Created by Sung Taek Kim on 5/12/14.
//
//

#ifndef _TIMEUTIL_H_
#define _TIMEUTIL_H_
#include <time.h>
#include <sys/time.h>

extern void init_time_marker(void);
extern double get_epoch_since_midnight(const struct timeval* time);
extern double timeval2double(const struct timeval *time);

#endif
