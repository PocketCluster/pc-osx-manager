#!/usr/bin/env python

import subprocess as sp
import os, pwd, sys, getpass


def hosts_cleanup_old(hosts_file="/etc/hosts"):
    with open(hosts_file, "r+") as hosts:
        hosts_content = map(lambda l: str(l).strip(), hosts)
        lines = filter(lambda h: not len(filter(lambda e: e in h, ["pc-node", "salt", "pc-master", "POCKETCLUSTER"])), hosts_content)
        pcnodes = filter(lambda h: "pc-node" in h, hosts_content)
        #hosts.seek(0)
        #map(lambda l: hosts.write("%s\n" % l), lines)
        #hosts.truncate()
        rlist = map(lambda l: l.split(" ")[0:2], pcnodes)
        # http://stackoverflow.com/questions/952914/making-a-flat-list-out-of-list-of-lists-in-python
        #print sum(rlist, [])
        print reduce(lambda l, r: l+r,rlist)


def hosts_cleanup(hosts_file="hosts"):
    with open(hosts_file, "r+") as hosts:
        hosts_content = map(lambda l: str(l).strip(), hosts)
        hosts.seek(0)
        map(lambda l: hosts.write("%s\n" % l),
            filter(lambda h: not len(filter(lambda e: e in h,
                                            ["pc-node", "salt", "pc-master", "POCKETCLUSTER"])), hosts_content))
        hosts.truncate()
        return reduce(lambda l, r: l+r,
                      map(lambda l: l.split(" ")[0:2],
                          filter(lambda h: len(filter(lambda e: e in h, ["pc-node", "salt", "pc-master"])),
                                 hosts_content)))


def remove_files(target_path):
    sp.call("rm -rf " + target_path, shell=True)


def remove_residue():
    print os.environ['SUDO_USER']
    pwd.getpwuid(os.getuid()).pw_dir
    sp.call("rm " + os.environ['HOME'] + "/Library/Preferences/io.pocketcluster.manager.plist", shell=True)
    sp.call("defaults delete io.pocketcluster", shell=True)
    sp.call("defaults delete io.pocketcluster.manager", shell=True)
    sp.call("defaults delete io.pocketcluster.pocketcluster", shell=True)


def clean_hosts_registry(registry_file="known_hosts", remove_target=[]):
    with open(registry_file, "r+") as registry:
        registry_content = map(lambda l: str(l).strip(), registry)
        registry.seek(0)
        map(lambda l: registry.write("%s\n" % l),
            filter(lambda r: not r.split(" ")[0] in remove_target, registry_content))
        registry.truncate()


def clean_host_config(config_file="config", remove_target=[]):
    with open(config_file, "r+") as config:
        config_content = map(lambda l: str(l), config)
        config_group = list()
        cgroup = list()
        config_group.append(cgroup)
        for c in config_content:
            if "host " in c.lower():
                cgroup = list()
                config_group.append(cgroup)
            cgroup.append(c)
        config.seek(0)
        map(lambda l: config.write("%s" % l),
            sum(filter(lambda g: not len(filter(lambda r: r in " ".join(map(lambda l: l.strip(), g)), remove_target)),
                       config_group), []))
        config.truncate()


def klll_aux_process():
    sp.call("ps -efw | grep salt | grep -v grep | awk '{print $2}' | xargs kill", shell=True)


def cleanup_vbox():
    """
    sp.call("sudo route -nv add -net 10.211.55 -interface vboxnet0", shell=True)
    vboxmanage hostonlyif remove vboxnet1
    sudo route -nv add -net 10.211.55 -interface vboxnet0
    VBoxManage list vms > /dev/null
    VBoxManage hostonlyif ipconfig vboxnet0 --ip 10.211.55.1 --netmask 255.255.255.0
    """


def cleanup_ssh():
    """
    restore ssh directory
    restore /etc/hosts
    """


if __name__ == '__main__':
    if getpass.getuser() != "root":
        print "root priviledge is necessary to clean up\nPlease run 'sudo -s chmod +x ./uninstall && ./uninstall'"
        exit(0)

    real_user_home = os.path.expanduser("~" + os.environ['SUDO_USER'])

    hosts_cleanup("./hosts")
    #remove_files("./dummy")
    #klll_aux_process()

    exit(0)