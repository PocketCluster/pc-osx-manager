//
//  timeutil.c
//  OSXGLEssentials
//
//  Created by Sung Taek Kim on 5/12/14.
//
//

#include <stdio.h>
#include "timeutil.h"

static time_t time_since_midnight;

void init_time_marker(void)
{
    // time since midnight
    time_t t = time(NULL);
    struct tm tm = *localtime(&t);
    tm.tm_sec = 0;
    tm.tm_min = 0;
    tm.tm_hour = 0;
    time_since_midnight = mktime(&tm);
    
    printf("\n *** time_since_midnight %lf ***\n", (double)time_since_midnight);
}

double get_epoch_since_midnight(const struct timeval* time)
{
    return timeval2double(time) - (double)time_since_midnight;
}


double timeval2double(const struct timeval *time)
{
    return (double)(time->tv_sec + time->tv_usec * 0.000001);
}

#if 0
void printTime()
{
    struct timeval tv;
    time_t nowtime;
    struct tm *nowtm;
    char tmbuf[64], buf[64];
    
    gettimeofday(&tv, NULL);
    nowtime = tv.tv_sec;
    nowtm = localtime(&nowtime);
    strftime(tmbuf, sizeof tmbuf, "%Y-%m-%d %H:%M:%S", nowtm);
    snprintf(buf, sizeof buf, "%s.%06d", tmbuf, tv.tv_usec);
}
#endif