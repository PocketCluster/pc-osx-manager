#!/usr/bin/env python

__author__ = 'stkim1'

import socket, struct, bson, six, time, select, ConfigParser, os, sys, atexit, netifaces
from uuid import getnode as get_mac
from random import randint
from signal import SIGTERM
import RPi.GPIO as GPIO

PC_PROTO = 'pc_ver'
VERSION = '1.0.0'

ADDRESS = 'address'
NETMASK = 'netmask'
BROADCS = 'broadcast'
GATEWAY = 'gateway'
NAMESRV = 'dns-nameservers'
IFACE_KEYS = [ADDRESS, NETMASK, BROADCS, GATEWAY, NAMESRV]

MASTER_COMMAND_TYPE = "pc_ma_ct"
COMMAND_FIX_BOUND = "ct_fix_bound"

MASTER_BOUND_AGENT = "pc_ma_ba"
SLAVE_LOOKUP_AGENT = "pc_sl_la"

SLAVE_NODE_MACADDR = "pc_sl_nm"





POCKETCAST_GROUP = '239.193.127.127'
PAGENT_SEND_PORT = 10060
PAGENT_RECV_PORT = 10061

POCKETCAST_SEND = (POCKETCAST_GROUP, PAGENT_SEND_PORT)
POCKETCAST_RECV = (POCKETCAST_GROUP, PAGENT_RECV_PORT)

TIMEOUT = 3

CONFIG_PATH = '/etc/pocket/conf.ini'
NET_IFACE = '/etc/network/interfaces'

#---*---*---*---*---*---*---*---*---*DAE*---*---*---*---*---*---*---*---*---*---
class Daemon:
    """
    A generic daemon class.

    Usage: subclass the Daemon class and override the run() method
    """
    def __init__(self, pidfile, stdin='/dev/null', stdout='/dev/null', stderr='/dev/null'):
        self.stdin = stdin
        self.stdout = stdout
        self.stderr = stderr
        self.pidfile = pidfile


    def daemonize(self):
        """
        do the UNIX double-fork magic, see Stevens' "Advanced
        Programming in the UNIX Environment" for details (ISBN 0201563177)
        http://www.erlenstar.demon.co.uk/unix/faq_2.html#SEC16
        """
        try:
            pid = os.fork()
            if pid > 0:
                # exit first parent
                sys.exit(0)
        except OSError, e:
            sys.stderr.write("fork #1 failed: %d (%s)\n" % (e.errno, e.strerror))
            sys.exit(1)

        # decouple from parent environment
        os.chdir("/")
        os.setsid()
        os.umask(0)

        # do second fork
        try:
            pid = os.fork()
            if pid > 0:
                # exit from second parent
                sys.exit(0)
        except OSError, e:
            sys.stderr.write("fork #2 failed: %d (%s)\n" % (e.errno, e.strerror))
            sys.exit(1)

        # redirect standard file descriptors
        sys.stdout.flush()
        sys.stderr.flush()
        si = file(self.stdin, 'r')
        so = file(self.stdout, 'a+')
        se = file(self.stderr, 'a+', 0)
        os.dup2(si.fileno(), sys.stdin.fileno())
        os.dup2(so.fileno(), sys.stdout.fileno())
        os.dup2(se.fileno(), sys.stderr.fileno())

        # write pidfile
        atexit.register(self.delpid)
        pid = str(os.getpid())
        file(self.pidfile,'w+').write("%s\n" % pid)


    def delpid(self):
        os.remove(self.pidfile)

    #   0 if daemon has been started
    #   1 if daemon was already running
    #   2 if daemon could not be started
    def start(self):
        """
        Start the daemon
        """
        # Check for a pidfile to see if the daemon already runs
        try:
            pf = file(self.pidfile,'r')
            pid = int(pf.read().strip())
            pf.close()
        except IOError:
            pid = None

        if pid:
            message = "pidfile %s already exist. Daemon already running?\n"
            sys.stderr.write(message % self.pidfile)
            #sys.exit(1)
            return 1

        # Start the daemon
        self.daemonize()
        self.run()
        return 0

    #   0 if daemon has been stopped
    #   1 if daemon was already stopped
    #   2 if daemon could not be stopped
    #   other if a failure occurred
    def stop(self):
        """
        Stop the daemon
        """
        # Get the pid from the pidfile
        try:
            pf = file(self.pidfile, 'r')
            pid = int(pf.read().strip())
            pf.close()
        except IOError:
            pid = None

        if not pid:
            message = "pidfile %s does not exist. Daemon not running?\n"
            sys.stderr.write(message % self.pidfile)
            return 1 # not an error in a restart

        # Try killing the daemon process
        try:
            while 1:
                os.kill(pid, SIGTERM)
                time.sleep(0.1)

            # have stopped
            return 0

        except OSError, err:

            err = str(err)
            if err.find("No such process") > 0:
                if os.path.exists(self.pidfile):
                    os.remove(self.pidfile)
            else:
                print str(err)
                #sys.exit(1)
            return 1

    def restart(self):
        """
        Restart the daemon
        """
        self.stop()
        return self.start()

    def run(self):
        """
        You should override this method when you subclass Daemon. It will be called after the process has been
        daemonized by start() or restart().
        """




#---*---*---*---*---*---*---*---*---*UTIL*---*---*---*---*---*---*---*---*---*---

class PocketDaemon(Daemon):

    def run(self):
        while True:
            time.sleep(3)

def getMacHexStr():
    return str(hex(get_mac()))[2:-1]

def keyValueDict(key, value):
    return {key:six.text_type(value, encoding='utf-8', errors='replace')}

def convert(input):
    if isinstance(input, dict):
        return {convert(key): convert(value) for key, value in input.iteritems()}
    elif isinstance(input, list):
        return [convert(element) for element in input]
    elif isinstance(input, unicode):
        return input.encode('utf-8')
    else:
        return input

# write back interface configuration file to /etc/networking/interfaces
def redef_ifaces(filepath, redefs):
    with open(filepath, 'r+') as iface:
        rplc = list()
        eth = False
        for line in iface:
            l = str(line).strip()

            if l.startswith('#') or not len(l):
                rplc.append(l)
                continue

            else:

                if l.startswith('iface eth0 inet'):
                    eth = True
                    rplc.append("iface eth0 inet static")
                    for k, v in redefs.iteritems():
                        if k in IFACE_KEYS:
                            rplc.append("{} {}".format(k, v))

                else:
                    if eth:
                        if l.startswith('address') or l.startswith('netmask') or l.startswith('gateway') or l.startswith('dns-nameservers') or l.startswith('broadcast'):
                            continue
                        else:
                            rplc.append(l)
                    else:
                        rplc.append(l)

        iface.seek(0)
        for l in rplc:
            iface.write("%s\n" % l)

        iface.truncate()


def get_iface_stat():
    estat = netifaces.ifaddresses('eth0')
    gws = netifaces.gateways()
    gw, iface = gws['default'][netifaces.AF_INET]

    redefs = dict()

    redefs.update({ADDRESS: estat[netifaces.AF_INET][0]['addr']})
    redefs.update({NETMASK: estat[netifaces.AF_INET][0]['netmask']})
    redefs.update({BROADCS: estat[netifaces.AF_INET][0]['broadcast']})
    redefs.update({GATEWAY: gw})
    redefs.update({NAMESRV: '8.8.8.8'})

    return redefs


def saveConfig(config):
    # Writing our configuration file to 'example.cfg'
    with open(CONFIG_PATH, 'wb') as configfile:
        config.write(configfile)


def buildInitConfig(config):
    #config = ConfigParser.RawConfigParser()

    # When adding sections or items, add them in the reverse order of
    # how you want them to be displayed in the actual file.
    # In addition, please note that using RawConfigParser's and the raw
    # mode of ConfigParser's respective set functions, you can assign
    # non-string values to keys internally, but will receive an error
    # when attempting to write to a file or when you get it in non-raw
    # mode. SafeConfigParser does not allow such assignments to take place.
    config.add_section('global')
    config.set('global', 'version', '1.0.0')
    config.set('global', 'state', 'unbound')

    config.add_section('agent')
    config.set('master', 'ip4', '')
    config.set('master', 'bound-id', '')
    config.set('master', 'hostname', '')

    config.add_section('node')
    config.set('node', 'mac', getMacHexStr())
    config.set('node', 'name', '')
    config.set('node', 'ip4', '')
    config.set('node', 'gateway', '')
    config.set('node', 'netmask', '')

    saveConfig(config)


def loadConfig():

    config = ConfigParser.RawConfigParser()
    try:
        config.read(CONFIG_PATH)
    except Exception as e:
        print str(e)
        os.makedirs('/etc/pocket', 0644)
        buildInitConfig(config)

    finally:
        return config


class PocketAgentDiscover:
    def __init__(self, config):
        self.ma = config.get('node', 'mac')
        self.agentcert = config.get('agent', 'cert')

    def unboundBroadcast(self):
        packet = dict()
        packet.update(keyValueDict(PC_PROTO, VERSION))
        packet.update(keyValueDict(SLAVE_NODE_MACADDR, self.ma))
        packet.update(keyValueDict(MASTER_BOUND_AGENT, SLAVE_LOOKUP_AGENT))

        # get current interface status
        istat = get_iface_stat()
        for k, v in istat.iteritems():
            packet.update(keyValueDict(k,v))

        return bson.dumps(packet)

    def boundBroadcast(self):
        packet = dict()
        packet.update(keyValueDict(PC_PROTO, VERSION))
        packet.update(keyValueDict(SLAVE_NODE_MACADDR, self.ma))
        packet.update(keyValueDict(MASTER_BOUND_AGENT, self.agentcert))
        return bson.dumps(packet)


class PocketNodeDiscover:
    def __init__(self, data, transID):
        self.data = data
        self.DHCPServerIdentifier = ''
        self.transID = transID
        self.offerIP = ''
        self.leaseTime = ''
        self.router = ''
        self.subnetMask = ''
        self.DNS = []
        self.unpack()

    def unpack(self):
        b = bson.loads(self.data)
        print b
        #if self.data[4:8] == self.transID:

    def printOffer(self):
        key = ['DHCP Server', 'Offered IP address', 'subnet mask', 'lease time (s)' , 'default gateway']
        val = [self.DHCPServerIdentifier, self.offerIP, self.subnetMask, self.leaseTime, self.router]
        for i in range(4):
            print('{0:20s} : {1:15s}'.format(key[i], val[i]))

        #print('{0:20s}'.format('DNS Servers') + ' : ', end='')
        print "{0:20s}".format('DNS Servers')
        if self.DNS:
            print('{0:15s}'.format(self.DNS[0]))
        if len(self.DNS) > 1:
            for i in range(1, len(self.DNS)):
                print('{0:22s} {1:15s}'.format(' ', self.DNS[i]))


if __name__ == '__main__':


    """
    daemon = PocketDaemon('/var/run/pocket.pid')
    if len(sys.argv) == 2:

        if 'start' == sys.argv[1]:
            ret = daemon.start()
            sys.exit(ret)

        elif 'stop' == sys.argv[1]:
            ret = daemon.stop()
            sys.exit(ret)

        elif 'restart' == sys.argv[1]:
            ret = daemon.restart()
            sys.exit(ret)

        else:
            print "Unknown command"
            sys.exit(2)

    else:
        print "usage: %s start|stop|restart" % sys.argv[0]
        sys.exit(2)
    """


    pocketagent = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)

    #receiver option
    mreq = struct.pack("4sl", socket.inet_aton(POCKETCAST_GROUP), socket.INADDR_ANY)
    pocketagent.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    pocketagent.setsockopt(socket.IPPROTO_IP, socket.IP_ADD_MEMBERSHIP, mreq)

    #sender option
    pocketagent.setsockopt(socket.IPPROTO_IP, socket.IP_MULTICAST_TTL, 2)

    #bind socket
    try:
        pocketagent.bind(POCKETCAST_RECV)
    except Exception as e:
        print 'port', POCKETCAST_RECV, 'in use...'
        pocketagent.close()
        exit()


    # load config
    config = loadConfig()

    #buiding and sending the DHCPDiscover packet
    discoverPacket = PocketAgentDiscover(config)


    # in case we are not bounded yet
    while True:
        try:

            pocketagent.sendto(discoverPacket.unboundBroadcast(), POCKETCAST_SEND)
            print 'POCKET Discover sent waiting for reply...\n'

            start = time.time()
            rd, wr, ex = select.select([pocketagent.fileno()], [], [], TIMEOUT)
            elapse = time.time() - start

            if rd:
                data = pocketagent.recv(1024)


                try:
                    cmd = convert(bson.loads(data))

                    print cmd

                    if cmd[MASTER_COMMAND_TYPE] == COMMAND_FIX_BOUND:
                        print "FIX BOUND"

                        # save config to interface
                        redef_ifaces("/home/ubuntu/interfaces", cmd)

                        # config salt

                        # config consul

                        # save config





                except Exception as e:
                    print str(e)
                    pass



            if elapse <= 2.0:
                time.sleep(3.0 - elapse)

        except socket.timeout as e:
            print(e)

    exit(0)




























    # when agent is bound.
    if len(config.get('agent', 'cert')):

        #check if agent is alive and accept connection
        if len(config.get('agent', 'ip4')):

            while True:
                try:

                    pocketagent.sendto(discoverPacket.boundBroadcast(), POCKETCAST_SEND)

                    print('POCKET Discover sent waiting for reply...\n')

                    rd, wr, ex = select.select([pocketagent.fileno()], [], [], TIMEOUT)
                    if rd:
                        data = pocketagent.recv(1024)
                        print "data found"
                        print bson.loads(data)

                        """
                        offer = PocketNodeDiscover(data, discoverPacket.transactionID)
                        if offer.offerIP:
                            offer.printOffer()
                            break
                        """
                        time.sleep(3)

                except socket.timeout as e:
                    print(e)

        else:

            while True:
                try:

                    pocketagent.sendto(discoverPacket.boundBroadcast(), POCKETCAST_SEND)

                    print('POCKET Discover sent waiting for reply...\n')

                    rd, wr, ex = select.select([pocketagent.fileno()], [], [], TIMEOUT)
                    if rd:
                        data = pocketagent.recv(1024)
                        print "data found"
                        print bson.loads(data)

                        """
                        offer = PocketNodeDiscover(data, discoverPacket.transactionID)
                        if offer.offerIP:
                            offer.printOffer()
                            break
                        """

                    time.sleep(3)

                except socket.timeout as e:
                    print(e)

    # when no agent is bound.
    else:

        while True:
            try:

                pocketagent.sendto(discoverPacket.unboundBroadcast(), POCKETCAST_SEND)

                print('POCKET Discover sent waiting for reply...\n')

                rd, wr, ex = select.select([pocketagent.fileno()], [], [], TIMEOUT)
                if rd:
                    data = pocketagent.recv(1024)
                    print bson.loads(data)

                    # check offer, check conf then write it to conf and setup everything
                    #offer = PocketNodeDiscover(data, discoverPacket.transactionID)
                    break

                time.sleep(3)

            except socket.timeout as e:
                print(e)


    pocketagent.close()   #we close the socket



    print "Config complete"

    exit(0)
