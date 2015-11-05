import subprocess

if __name__ == "__main__":

    lve = subprocess.call("VBoxManage list vms > /dev/null", shell=True)
    if lve:exit(lve)

    vboxnet = dict()
    for iface in [il.strip() for il in subprocess.check_output("ifconfig -l", shell=True).split(" ")]:
        if iface.startswith("vboxnet"):
            vboxnet[iface] = ""
            for p in [pl.strip() for pl in subprocess.check_output("ifconfig " + iface, shell=True).split("\n")]:
                if p.startswith("inet"):
                    inet = p.split(" ")[1]
                    vboxnet[iface] = inet
                    if inet == "10.211.55.1":
                        print iface
                        exit(0)

    ## by the time you reach at this point, you need a new interface
    nIface = "vboxnet{}".format(len(vboxnet))
    subprocess.call("VBoxManage hostonlyif create", shell=True)
    subprocess.call("VBoxManage hostonlyif ipconfig {} --ip 10.211.55.1 --netmask 255.255.255.0".format(nIface), shell=True)
    print nIface

    exit(0)
