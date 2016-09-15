#ifndef __DEBUG_H__
#define __DEBUG_H__

#include "acfg.h"
#include "aclogger.h"
#include "meta.h"

#ifdef DEBUG
    #include <assert.h>
    #define ASSERT(x) assert(x)
#else
    #define ASSERT(x)
#endif

#ifndef DEBUG
    #define ldbgvl(v, x)
    #define dbglvl(v)
    #define ldbg(x)
    #define dbgline
    #define ASSERT(x)
    #define LOG(x)
    #define LOGSTART(x)
    #define LOGSTARTs(x)
    #define LOGSTART2(x,y)
    #define LOGSTART2s(x,y)
    #define DBGQLOG(x)
    // strip away
    inline void dump_proc_status()
    {};
#else
    #include <fstream>
    #include <iostream>
    #define LOGLVL(n, x) if(acfg::debug&n) \
    { \
        __logobj.GetFmter() << x; __logobj.Write(__FILE__,__LINE__); \
    }
    #define LOG(x) LOGLVL(LOG_DEBUG, x)
    #define LOGSTART(x) t_logger __logobj(x, this);
    #define LOGSTARTs(x) t_logger __logobj(x, NULL);
    #define LOGSTART2(x, y) t_logger __logobj(x, this); LOGLVL(LOG_DEBUG, y)
    #define LOGSTART2s(x, y) t_logger __logobj(x, NULL); LOGLVL(LOG_DEBUG, y)
    #define ldbg(x) LOG(x)
    #define dbgline ldbg("mark")
    #define DBGQLOG(x) {aclog::err(tSS()<< x);}
    inline void dump_proc_status()
    {
        using namespace std;
        ifstream sf("/proc/self/status");
        while (sf)
        {
            string s;
            getline(sf, s);
            cerr << s << endl;
        }
    };
#endif

#endif // __DEBUG_H__
