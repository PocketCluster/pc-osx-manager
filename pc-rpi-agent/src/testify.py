#!/usr/bin/env python

__author__ = 'stkim1'

import socket, struct, bson, six, time, select, ConfigParser, os, sys, atexit, netifaces
from uuid import getnode as get_mac
from random import randint
from signal import SIGTERM
#import RPi.GPIO as GPIO
from subprocess import call
from pocketd.daemonizer.Daemon import Daemon

from pocketd.const import *
from pocketd.conf import *
from pocketd.agent.PocketAgentDiscover import PocketAgentDiscover
from pocketd.hostconf.netiface import redef_ifaces
from pocketd.hostconf.hostaddrs import redef_hostaddrs
from pocketd.hostconf.hostname import redef_hostname
from pocketd.hostconf.timezone import redef_timezone
from pocketd.hostconf.saltminion import redef_salt_minion_id
from pocketd.util import convert_to_ascii_dict, convert_to_hostaddr

if __name__ == "__main__":

    pocsock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)

    #receiver option
    mreq = struct.pack("4sl", socket.inet_aton(POCKETCAST_GROUP), socket.INADDR_ANY)
    pocsock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    pocsock.setsockopt(socket.IPPROTO_IP, socket.IP_ADD_MEMBERSHIP, mreq)

    #sender option
    pocsock.setsockopt(socket.IPPROTO_IP, socket.IP_MULTICAST_TTL, 2)

    #bind socket
    try:
        pocsock.bind(POCKETCAST_RECV)
    except Exception as e:
        print 'port {} in use...'.format(POCKETCAST_RECV)
        pocsock.close()
        exit()

    # load config
    config = loadConfig()

    #buiding and sending the DHCPDiscover packet
    discoverer = PocketAgentDiscover(config)

    #last datetime
    last_datetime = 0

    # basically this loop keeps running until computer stops!
    while True:
        try:
            wait_timer = UNBOUNDED_TIMEOUT

            boundness = config.get('global', 'status')
            if boundness == 'unbound':
                pocsock.sendto(discoverer.unboundBroadcast(), POCKETCAST_SEND)
                wait_timer = UNBOUNDED_TIMEOUT
                print '[UNBOUNDED] POCKET Discover sent waiting for reply...\n'

            elif boundness == 'bounded':
                pocsock.sendto(discoverer.boundedBroadcast(), POCKETCAST_SEND)
                wait_timer = BOUNDED_TIMEOUT
                print '[BOUNDED] POCKET Checker sent waiting for reply...\n'


            start = time.time()
            rd, wr, ex = select.select([pocsock.fileno()], [], [], wait_timer)
            elapse = time.time() - start

            if rd:
                data = pocsock.recv(1024)

                try:
                    cmd = convert_to_ascii_dict(discoverer.unpack(data))

                    # state #1
                    if boundness == 'unbound':
                        if cmd[MASTER_COMMAND_TYPE] == COMMAND_FIX_BOUND and \
                           cmd[SLAVE_NODE_MACADDR] == config.get(SLAVE_SECTION, SLAVE_NODE_MACADDR):

                            print "FIX BOUND"

                            # save config to interface
                            redef_ifaces(NET_IFACE, cmd)

                            hostname = cmd[SLAVE_NODE_NAME]

                            # config hostname
                            redef_hostname(hostname)

                            # salt minion id
                            redef_salt_minion_id(hostname)

                            # config host address
                            hostaddr = convert_to_hostaddr(cmd)
                            redef_hostaddrs(hostaddr)

                            # config timezone
                            redef_timezone(cmd[MASTER_TIMEZONE])

                            # save config
                            conf_dict = buildBoundedConfig(cmd)
                            config = bindConfigValue(config, conf_dict)
                            saveConfig(config)

                            #set timezone
                            os.system("dpkg-reconfigure -f noninteractive tzdata")

                            #set system time
                            master_datetime = cmd[MASTER_DATETIME]
                            os.system("date +%s -s @" + master_datetime)
                            last_datetime = int(master_datetime)

                            # re-partition
                            call("sh /repartition.sh", shell=True)

                            # restart
                            call("reboot", shell=True)

                            wait_timer = BOUNDED_TIMEOUT

                    # state #2
                    elif boundness == 'bounded':
                        # since this is bounded state, there must be a master. check boundness
                        if cmd[MASTER_BOUND_AGENT] == config.get(MASTER_SECTION, MASTER_BOUND_AGENT) and \
                           cmd[SLAVE_NODE_MACADDR] == config.get(SLAVE_SECTION, SLAVE_NODE_MACADDR):

                            # sync time every 60 seconds
                            master_datetime = cmd[MASTER_DATETIME]
                            if 60 < (int(master_datetime) - last_datetime):
                                print "date +%s -s @" + master_datetime
                                os.system("date +%s -s @" + master_datetime)
                                last_datetime = int(master_datetime)

                            # check timezone
                            if cmd[MASTER_TIMEZONE] != config.get(MASTER_SECTION, MASTER_TIMEZONE):
                                #set system timezone
                                redef_timezone(cmd[MASTER_TIMEZONE])

                                #set timezone
                                os.system("dpkg-reconfigure -f noninteractive tzdata")

                                #save changed config
                                config.set(MASTER_SECTION, MASTER_TIMEZONE, cmd[MASTER_TIMEZONE])
                                saveConfig(config)

                            # check interface
                            if cmd[MASTER_IP4_ADDRESS] != config.get(MASTER_SECTION, MASTER_IP4_ADDRESS):
                                hostaddr = convert_to_hostaddr(cmd)
                                redef_hostaddrs(hostaddr)

                                print "MASTER IP ADDRESS HAS CHANGED. CHANGE INTERFACE FILE"

                                # need to save config
                                config.set(MASTER_SECTION, MASTER_IP4_ADDRESS, cmd[MASTER_IP4_ADDRESS])
                                saveConfig(config)

                                # restart network
                                call(["service", "networking", "restart"])
                                # restart salt
                                call(["service", "salt-minion", "restart"])

                            wait_timer = BOUNDED_TIMEOUT

                except Exception as e:
                    print str(e)
                    pass

            # state #3
            else:

                # could not found a pocket agent (master) broadcast about where an agent is
                if boundness == 'bounded':
                    print "WE NEED A MASTER!"
                    wait_timer = UNBOUNDED_TIMEOUT

            if elapse <= (wait_timer - 1.0):
                time.sleep(wait_timer - elapse)

        except socket.timeout as e:
            print str(e)

    pocsock.close()