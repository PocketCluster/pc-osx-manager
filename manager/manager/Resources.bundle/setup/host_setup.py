#!/usr/bin/python
import sys

POCKET_START = '# --------------- POCKETCLUSTER START ---------------'
POCKET_END = '# ---------------  POCKETCLUSTER END  ---------------'
HOSTADDR_FILE = '/etc/hosts'

def handle_pocket_section(hostsline=[], hostaddrs={}):
    # start of pocket section
    hostsline.append(POCKET_START)

    for name, addr in hostaddrs.iteritems():
        hostsline.append("{} {} {}".format(addr, name, name))

    # end of pocket section
    hostsline.append(POCKET_END)


def redef_hostaddrs(hostaddrs={}):

    with open(HOSTADDR_FILE, "r+") as hostsfile:

        hostsline = list()
        is_pocket_defined = False
        is_pocket_editing = False

        for hl in hostsfile:
            line = str(hl).strip()

            if line == POCKET_START:
                is_pocket_defined = True
                is_pocket_editing = True
                handle_pocket_section(hostsline, hostaddrs)
                continue

            if line == POCKET_END:
                is_pocket_editing = False
                continue

            if not is_pocket_editing:
                # this is a special case where you need to skip
                if "127.0.0.1" in line and not "localhost" in line:
                    continue
                else:
                    hostsline.append(line)

        hostsfile.seek(0)

        if not is_pocket_defined:
            handle_pocket_section(hostsline, hostaddrs)

        for h in hostsline:
            hostsfile.write("%s\n" % h)

        hostsfile.truncate()


if __name__ == '__main__':

    hostlen = (len(sys.argv) - 1) / 2
    hostaddrs = dict()
    for i in xrange(0, hostlen):
        hostaddrs.update({sys.argv[i * 2 + 1]: sys.argv[i * 2 + 2]})

    redef_hostaddrs(hostaddrs)
    exit(0)