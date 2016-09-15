
with open('interfaces', 'r+') as iface:

    rplc = list()
    eth = False

    for line in iface:
        l = str(line).strip()

        if l.startswith('#') or not len(l):
            rplc.append(l)
            continue

        else:

            if l.startswith('iface eth0 inet'):
                rplc.append("iface eth0 inet static")
                rplc.append("address 192.168.1.234")
                rplc.append("netmask 255.255.255.0")
                rplc.append("gateway 192.168.1.1")
                rplc.append("dns-nameservers 8.8.8.8")
                eth = True

            else:

                if eth:
                    if l.startswith('address') or l.startswith('netmask') or l.startswith('gateway') or l.startswith('dns-nameservers'):
                        continue
                    else:
                        rplc.append(l)
                else:
                    rplc.append(l)

    iface.seek(0)
    for l in rplc:
        iface.write("%s\n" % l)

    iface.truncate()