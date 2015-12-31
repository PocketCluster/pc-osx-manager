__author__ = 'stkim1'

import socket, bson

from pocketd.const import *
from pocketd.util import keyValueDict
from pocketd.hostconf.netiface import get_iface_stat
from pocketd.hostconf.timezone import get_timezone

class PocketAgentDiscover:
    def __init__(self, config):
        self._config = config

    def unpack(self, package):
        return bson.loads(package)

    def unboundBroadcast(self):

        packet = dict()
        packet.update(keyValueDict(PC_PROTO, VERSION))

        # master
        packet.update(keyValueDict(MASTER_BOUND_AGENT, SLAVE_LOOKUP_AGENT))

        # slave
        node_mac = self._config.get(SLAVE_SECTION, SLAVE_NODE_MACADDR)
        host_name = socket.gethostname()
        packet.update(keyValueDict(SLAVE_NODE_MACADDR, node_mac))
        packet.update(keyValueDict(SLAVE_NODE_NAME, host_name))

        # get current interface status
        istat = get_iface_stat()
        for k, v in istat.iteritems():
            packet.update(keyValueDict(k,v))

        return bson.dumps(packet)

    def boundedBroadcast(self):
        packet = dict()
        packet.update(keyValueDict(PC_PROTO, VERSION))

        # master
        bounded_master = self._config.get(MASTER_SECTION, MASTER_BOUND_AGENT)
        packet.update(keyValueDict(MASTER_BOUND_AGENT, bounded_master))

        # slave
        node_mac = self._config.get(SLAVE_SECTION, SLAVE_NODE_MACADDR)
        host_name = socket.gethostname()
        packet.update(keyValueDict(SLAVE_NODE_MACADDR, node_mac))
        packet.update(keyValueDict(SLAVE_NODE_NAME, host_name))
        packet.update(keyValueDict(SLAVE_TIMEZONE,self._config.get(MASTER_SECTION, MASTER_TIMEZONE)))

        # get current interface status
        istat = get_iface_stat()
        for k, v in istat.iteritems():
            packet.update(keyValueDict(k,v))

        return bson.dumps(packet)
