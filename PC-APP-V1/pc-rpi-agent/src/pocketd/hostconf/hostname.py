__author__ = 'stkim1'

from pocketd.const import *

def redef_hostname(hostname="pc-node"):
    with open(HOSTNAME_FILE, "w") as hostfile:
        hostfile.write(hostname)

def get_hostname():
    with open(HOSTNAME_FILE, "r") as hostfile:
        return hostfile.read()
