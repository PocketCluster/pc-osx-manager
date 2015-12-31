__author__ = 'stkim1'

from pocketd.const import *
import netifaces

# write back interface configuration file to /etc/networking/interfaces
def handle_pocket_section(ifacelines=[], ifaceredefs={}):
    # start of pocket section
    ifacelines.append(POCKET_START)

    ifacelines.append("iface eth0 inet static")
    for k, v in ifaceredefs.iteritems():
        if k in IFACE_KEYS:
            ifacelines.append("{} {}".format(k, v))

    # end of pocket section
    ifacelines.append(POCKET_END)


def redef_ifaces(filepath, redefs):
    with open(filepath, 'r+') as ifacefile:

        ifacelines = list()
        eth = False

        is_pocket_defiend = False
        is_pocket_editing = False

        # first scan
        for l in ifacefile:
            line = str(l).strip()

            if line == POCKET_START:
                is_pocket_defiend = True
                is_pocket_editing = True
                handle_pocket_section(ifacelines, redefs)
                continue

            if line == POCKET_END:
                is_pocket_editing = False
                continue

            if not is_pocket_editing:
                ifacelines.append(line)


        # second scan in case there is no pocket section
        if not is_pocket_defiend:

            ifacefile.seek(0)

            #empty iface lines
            del ifacelines[:]

            for l in ifacefile:
                line = str(l).strip()

                if line.startswith('iface eth0 inet'):
                    ifacelines.append(POCKET_START)
                    ifacelines.append("iface eth0 inet static")
                    for k, v in redefs.iteritems():
                        if k in IFACE_KEYS:
                            ifacelines.append("{} {}".format(k, v))
                    ifacelines.append(POCKET_END)
                    continue

                if line.split(' ')[0] in IFACE_KEYS:
                    continue

                ifacelines.append(line)


        ifacefile.seek(0)

        for l in ifacelines:
            ifacefile.write("%s\n" % l)

        ifacefile.truncate()


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
