#!/bin/bash

PCIFACE=$1
VBoxManage list vms 2>&1
VBoxManage hostonlyif ipconfig ${PCIFACE} --ip 10.211.55.1 --netmask 255.255.255.0

exit 0