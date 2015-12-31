__author__ = 'stkim1'

from pocketd.const import *
import netifaces

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