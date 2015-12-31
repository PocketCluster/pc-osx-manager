#!/usr/bin/env python

__author__ = 'almightykim'

import sys

if __name__ == '__main__':

    if len(sys.argv) == 2:
        if 'start' == sys.argv[1]:
            exit(0)
        elif 'stop' == sys.argv[1]:
            sys.exit(2)
        elif 'restart' == sys.argv[1]:
            sys.exit(3)
        else:
            print "Unknown command"
            sys.exit(2)
        sys.exit(0)
    else:
        print "usage: %s start|stop|restart" % sys.argv[0]
        sys.exit(2)

    sys.exit(0)