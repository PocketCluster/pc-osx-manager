__author__ = 'stkim1'

import six
from uuid import getnode as get_mac
from pocketd.const import *

def getMacHexStr():
    return str(hex(get_mac()))[2:-1]


def keyValueDict(key, value):
    return {key:six.text_type(value, encoding='utf-8', errors='replace')}


def convert_to_ascii_dict(inData):

    def convert(itrObject):
        if isinstance(itrObject, dict):
            return {convert(key): convert(value) for key, value in itrObject.iteritems()}
        elif isinstance(itrObject, list):
            return [convert(element) for element in itrObject]
        elif isinstance(itrObject, unicode):
            return itrObject.encode('utf-8')
        else:
            return itrObject

    return convert(inData)


def convert_to_hostaddr(hostsinfo):

    snn = hostsinfo[SLAVE_NODE_NAME]
    hosts = dict()
    hosts[SALT_MASTER] = hostsinfo[MASTER_IP4_ADDRESS]
    hosts[PC_MASTER] = hostsinfo[MASTER_IP4_ADDRESS]
    hosts[snn] = hostsinfo[ADDRESS]
    cms = hostsinfo[SLAVE_CLUSTER_MEMBERS]
    for hname, ip in cms.iteritems():
        if snn == hname:
            continue
        else:
            hosts[hname] = ip
    return hosts
