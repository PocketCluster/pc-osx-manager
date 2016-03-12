
//#define LOCAL_DEBUG
#include "debug.h"

#include "config.h"

#include "acbuf.h"
#include "fileio.h"
#include <unistd.h>
#include <algorithm>

bool acbuf::setsize(size_t c) {
    if (m_nCapacity == c)
        return true;
    
    char *p = (char*) realloc(m_buf, c+1);
    if (!p)
        return false;
    
    m_buf = p;
    m_nCapacity = c;
    m_buf[c] = 0; // terminate to make string operations safe
    return true;
}

bool acbuf::initFromFile(const char *szPath)
{
    struct stat statbuf;
    
    if (0 != stat(szPath, &statbuf))
        return false;
    
    int fd=open(szPath, O_RDONLY);
    if (fd < 0)
        return false;
    
    clear();
    
    if(!setsize(statbuf.st_size))
        return false;
    
    while (freecapa() > 0)
    {
        if (sysread(fd) < 0)
        {
            forceclose(fd);
            return false;
        }
    }
    forceclose(fd);
    return true;
}

ssize_t acbuf::syswrite(int fd, size_t maxlen)
{
    size_t todo(std::min(maxlen, size()));
    ssize_t n;
    do
    {
        n = ::write(fd, rptr(), todo);
    } while (n < 0 && errno == EINTR);
    
    if (n < 0 && errno == EAGAIN)
        n = 0;
    if (n < 0)
        return -errno;
    drop(n);
    return n;
}

ssize_t acbuf::sysread(int fd)
{
    ssize_t n;
    do {
        n = ::read(fd, m_buf + w, m_nCapacity - w);
    } while((n < 0 && EINTR == errno)/* || (EAGAIN == errno && n <= 0)*/); // cannot handle EAGAIN here, let the caller check errno
    if (n < 0)
        return - errno;
    if (n > 0)
        w += n;
    return n;
}

tSS& tSS::addEscaped(const char *fmt)
{
    if(!fmt || !*fmt)
        return *this;
    size_t nl = strlen(fmt);
    reserve(length() + nl);
    char *p = wptr();
    
    for(; *fmt ; fmt++)
    {
        if(*fmt == '\\')
            *(p++) = unEscape(*(++fmt));
        else
            *(p++) = *fmt;
    }
    got(p-wptr());
    return *this;
}
