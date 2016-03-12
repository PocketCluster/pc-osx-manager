
//#define LOCAL_DEBUG
#include "debug.h"

#include "acfg.h"

#include "header.h"
#include "config.h"
#include <acbuf.h>

#include <cstdio>
#include <iostream>
#include <string.h>
#include <unistd.h>

#include "fileio.h"
#include "filereader.h"

#include <map>

using namespace std;

#if 1
struct eHeadPos2label
{
    header::eHeadPos pos;
    const char *str;
};

eHeadPos2label mapId2Headname[] =
{
    { header::LAST_MODIFIED, "Last-Modified"},
    { header::CONTENT_LENGTH, "Content-Length"},
    { header::CONNECTION, "Connection"},
    { header::CONTENT_TYPE, "Content-Type"},
    { header::IF_MODIFIED_SINCE, "If-Modified-Since"},
    { header::RANGE, "Range"},
    { header::IFRANGE, "If-Range"},
    { header::CONTENT_RANGE, "Content-Range"},
    { header::PROXY_CONNECTION, "Proxy-Connection"},
    { header::TRANSFER_ENCODING, "Transfer-Encoding"},
    { header::AUTHORIZATION, "Authorization"},
    { header::LOCATION, "Location" },
    { header::XFORWARDEDFOR, "X-Forwarded-For"},
    { header::XORIG, "X-Original-Source"}
};
#endif
#if 0 // nonsense... save a penny, waste an hour
struct tHeadLabelMap
{
    class noCaseComp
    {
        //	bool operator<(const tStringRef &a) { return strncasecmp(a.first, first, second)<0; }
    };
    map<pair<const char*,size_t>, header::eHeadPos> lookup;
    tHeadLabelMap()
    {
        //tHeadLabelMap &x=*this;
        insert(make_pair(tStringRef(NAMEWLEN("foo")), header::XORIG));
    }
} label_map;
#endif

header::header(const header &s)
:type(s.type),
frontLine(s.frontLine),
m_nEstimLength(s.m_nEstimLength)
{
    for (UINT i = 0; i < HEADPOS_MAX; i++)
        h[i] = s.h[i] ? strdup(s.h[i]) : NULL;
}

header& header::operator=(const header& s)
{
    type=s.type;
    frontLine=s.frontLine;
    m_nEstimLength=s.m_nEstimLength;
    for (UINT i = 0; i < HEADPOS_MAX; ++i)
    {
        if (h[i])
            free(h[i]);
        h[i] = s.h[i] ? strdup(s.h[i]) : NULL;
    }
    return *this;
}

header::~header()
{
    for(auto& p:h)
        free(p);
}

void header::clear()
{
    for(UINT i=0; i<HEADPOS_MAX; i++)
        del((eHeadPos) i);
    frontLine.clear();
    type=INVALID;
    m_nEstimLength=0;
}

void header::del(eHeadPos i)
{
    free(h[i]);
    h[i]=0;
}

inline ssize_t header::Load(const char* const in, size_t maxlen)
{
    if (maxlen < 9)
        return 0;
    
    if (!in)
        return -1;
    if (!strncmp(in,  "HTTP/1.", 7))
        type = ANSWER;
    else if (!strncmp(in, "GET ", 4))
        type = GET;
    else if (!strncmp(in, "HEAD ", 5))
        type = HEAD;
    else if (!strncmp(in, "POST ", 5))
        type = POST;
    else if (!strncmp(in, "CONNECT ", 8))
        type = CONNECT;
    else
        return -1;
    
    const char *posNext = in;
    
    while (true)
    {
        const char *szBegin = posNext;
        size_t pos = static_cast<size_t>(szBegin - in);
        const char *end = (const char*) memchr(szBegin, '\r', maxlen - pos);
        if (!end)
            return 0;
        
        if (end + 1>= in + maxlen)
            return 0; // one newline must fit there, always
        
        if (szBegin == end)
        {
            if (end[1] =='\n')
            {
                m_nEstimLength = end + 2 - in;
                return m_nEstimLength; // end detected
            }
            
            return - 1; // looks like crap
        }
        posNext = end+2;
        
        while (isspace((UINT)*end))
            end--;
        end++;
        
        if (frontLine.empty())
        {
            frontLine.assign(in, end-in);
            trimBack(frontLine);
            continue;
        }
        
        // end is on the last relevant char now
        const char *sep=(const char*) memchr(szBegin, ':', end-szBegin);
        if (!sep)
            return -1;
        
        const char *key = szBegin;
        size_t keyLen=sep-szBegin;
        
        sep++;
        while (sep < end && isspace((UINT)*sep))
            sep++;
        
        for (const auto& id2key : mapId2Headname)
        {
            if (strncasecmp(id2key.str, key, keyLen))
                continue;
            
            unsigned int l = static_cast<unsigned int>(end - sep);
            if (!(h[id2key.pos] = (char*) realloc(h[id2key.pos], l + 1)))
                continue;
            
            memcpy(h[id2key.pos], sep, l);
            h[id2key.pos][l]='\0';
            break;
        }
    }
    return -2;
}

ssize_t header::LoadFromBuf(const char * const in, size_t maxlen)
{
    clear();
    ssize_t ret = Load(in, maxlen);
    if (ret < 0)
        clear();
    return ret;
}

ssize_t header::LoadFromFile(const string &sPath)
{
    clear();
#if 0
    filereader buf;
    return buf.OpenFile(sPath, true) && LoadFromBuf(buf.GetBuffer(), buf.GetSize());
#endif
    acbuf buf;
    if (!buf.initFromFile(sPath.c_str()))
        return -1;
    return LoadFromBuf(buf.rptr(), buf.size());
}


void header::set(eHeadPos i, const char *val)
{
    if (h[i])
    {
        free(h[i]);
        h[i]=NULL;
    }
    if(val)
        h[i] = strdup(val);
}

void header::set(eHeadPos i, const char *val, size_t len)
{
    if(!val)
    {
        free(h[i]);
        h[i]=NULL;
        return;
    }
    h[i] = (char*) realloc(h[i], len+1);
    if(h[i])
    {
        memcpy(h[i], val, len);
        h[i][len]='\0';
    }
}

void header::set(eHeadPos key, cmstring &value)
{
    string::size_type l=value.size()+1;
    h[key]=(char*) realloc(h[key], l);
    if(h[key])
        memcpy(h[key], value.c_str(), l);
}

void header::set(eHeadPos key, off_t nValue)
{
    char buf[3*sizeof(off_t)];
    int len=sprintf(buf, OFF_T_FMT, nValue);
    set(key, buf, len);
}

#ifndef MINIBUILD

tSS header::ToString() const
{
    tSS s;
    s<<frontLine << "\r\n";
    for(const auto& pos2key : mapId2Headname)
        if (h[pos2key.pos])
            s << pos2key.str << ": " << h[pos2key.pos] << "\r\n";
    s<< "Date: " << tCurrentTime() << "\r\n\r\n";
    return s;
}

ssize_t header::StoreToFile(cmstring &sPath) const
{
    ssize_t nByteCount(0);
    const char *szPath=sPath.c_str();
    int fd=open(szPath, O_WRONLY|O_CREAT|O_TRUNC, acfg::fileperms);
    if(fd<0)
    {
        fd=-errno;
        // maybe there is something in the way which can be removed?
        if(::unlink(szPath))
            return fd;
        
        fd=open(szPath, O_WRONLY|O_CREAT|O_TRUNC, acfg::fileperms);
        if(fd<0)
            return -errno;
    }
    
    auto hstr=ToString();
    const char *p=hstr.rptr();
    nByteCount = hstr.length();
    
    for(string::size_type pos=0; pos<(UINT)nByteCount;)
    {
        ssize_t ret = write(fd, p + pos, nByteCount-pos);
        if(ret<0)
        {
            if(EAGAIN == errno || EINTR == errno)
                continue;
            if(EINTR == errno)
                continue;
            
            ret = errno;
            forceclose(fd);
            return -ret;
        }
        pos += ret;
    }
    
    while (0 != close(fd))
    {
        if (errno != EINTR)
            return -errno;
    }

    return nByteCount;
}

#endif // MINIBUILD

std::string header::GenInfoHeaders()
{
    string ret = "Date: ";
    ret += tCurrentTime();
    ret += "\r\nServer: Debian Apt-Cacher NG/" ACVERSION "\r\n";
    return ret;
}

static const char* fmts[] =
{
    "%a, %d %b %Y %H:%M:%S GMT",
    "%A, %d-%b-%y %H:%M:%S GMT",
    "%a %b %d %H:%M:%S %Y"
};

bool header::ParseDate(const char *s, struct tm *tm)
{
    if (!s || !tm)
        return false;
    
    for (const auto& fmt : fmts)
        if (::strptime(s, fmt, tm))
            return true;

    return false;
}
