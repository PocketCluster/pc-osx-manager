__author__ = 'stkim1'

import ConfigParser, os
from pocketd.const import *
from pocketd.util import getMacHexStr
from pocketd.hostconf.hostname import get_hostname
from pocketd.hostconf.timezone import get_timezone

def saveConfig(config):
    # Writing our configuration file to 'example.cfg'
    with open(CONFIG_PATH, 'wb') as configfile:
        config.write(configfile)

def loadConfig():

    config = ConfigParser.RawConfigParser()

    if os.path.exists(CONFIG_PATH):
        config.read(CONFIG_PATH)
    else:
        try:
            os.makedirs('/etc/pocket', 0644)
        except:
            pass

        try:
            config = buildInitConfig(config)
            saveConfig(config)
        except Exception as e2:
            #import traceback;traceback.print_exc()
            #print str(e2)
            pass

    return config

def buildBoundedConfig(cmd):
    value = dict()
    value['global']       = {PC_PROTO: VERSION,
                             'status': 'bounded'}

    value[MASTER_SECTION] = {MASTER_IP4_ADDRESS: cmd[MASTER_IP4_ADDRESS],
                             MASTER_IP6_ADDRESS: cmd[MASTER_IP6_ADDRESS],
                             MASTER_BOUND_AGENT: cmd[MASTER_BOUND_AGENT],
                             MASTER_HOSTNAME: cmd[MASTER_HOSTNAME],
                             MASTER_TIMEZONE: cmd[MASTER_TIMEZONE]}

    value[SLAVE_SECTION]  = {SLAVE_NODE_MACADDR: cmd[SLAVE_NODE_MACADDR],
                             SLAVE_NODE_NAME: cmd[SLAVE_NODE_NAME],
                             ADDRESS: cmd[ADDRESS],
                             NETMASK: cmd[NETMASK],
                             BROADCS: cmd[BROADCS],
                             GATEWAY: cmd[GATEWAY],
                             NAMESRV: cmd[NAMESRV]}
    return value


def bindConfigValue(config, value):
    for sk, sv in value.iteritems():
        for nk, nv in sv.iteritems():
            config.set(sk, nk, nv)
    return config


def buildInitConfig(config):
    # When adding sections or items, add them in the reverse order of
    # how you want them to be displayed in the actual file.
    # In addition, please note that using RawConfigParser's and the raw
    # mode of ConfigParser's respective set functions, you can assign
    # non-string values to keys internally, but will receive an error
    # when attempting to write to a file or when you get it in non-raw
    # mode. SafeConfigParser does not allow such assignments to take place.
    config.add_section('global')
    config.set('global', PC_PROTO, VERSION)
    config.set('global', 'status', 'unbound')

    config.add_section(MASTER_SECTION)
    config.set(MASTER_SECTION, MASTER_IP4_ADDRESS, '')
    config.set(MASTER_SECTION, MASTER_IP6_ADDRESS, '')
    config.set(MASTER_SECTION, MASTER_BOUND_AGENT, '')
    config.set(MASTER_SECTION, MASTER_HOSTNAME, get_hostname())
    config.set(MASTER_SECTION, MASTER_TIMEZONE, get_timezone())

    config.add_section(SLAVE_SECTION)
    config.set(SLAVE_SECTION, SLAVE_NODE_MACADDR, getMacHexStr())
    config.set(SLAVE_SECTION, SLAVE_NODE_NAME, '')
    config.set(SLAVE_SECTION, ADDRESS, '')
    config.set(SLAVE_SECTION, NETMASK, '')
    config.set(SLAVE_SECTION, BROADCS, '')
    config.set(SLAVE_SECTION, GATEWAY, '')
    config.set(SLAVE_SECTION, NAMESRV, '')

    return config


