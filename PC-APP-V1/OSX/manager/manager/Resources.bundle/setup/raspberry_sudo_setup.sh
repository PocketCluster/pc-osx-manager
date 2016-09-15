#!/bin/bash

BASE_BUNDLE_PATH=$1
MASTER_ADDRESS=$2
ARGV=("${@}")
# ARGV=("${BASH_ARGV[*]}")
ARGC=$#
NODELIST=""

for (( i=2; i<$ARGC; i++ ));
do
    NODENUM=$(( $i - 1 ))
    NODELIST+=" pc-node${NODENUM} "${ARGV[$i]}
done

echo "SUDO_SETUP_STEP_0"

# setup root directories
mkdir -p /{pocket,bigpkg}
mkdir -p /bigpkg/archive

# copy vagrant files & modify
chown -R $SUDO_USER:staff /pocket
chown -R $SUDO_USER:staff /bigpkg
chown $SUDO_USER:admin /usr/local

# copy salt config files
mkdir -p /etc/salt

if [[ -n "/etc/salt/master" ]]; then
    mv /etc/salt/master /etc/salt/master.backup
fi

if [[ -n "/etc/salt/minion" ]]; then
    mv /etc/salt/minion /etc/salt/minion.backup
fi

cp -f "${BASE_BUNDLE_PATH}"/etc/salt/* /etc/salt/
sed -i '' 's|PC_USER|'$SUDO_USER'|g' /etc/salt/*

# check if this really works!
sed -i '' 's|10.211.55.1|0.0.0.0|g' /etc/salt/master

# change hosts
python "${BASE_BUNDLE_PATH}"/setup/host_setup.py salt ${MASTER_ADDRESS} pc-master ${MASTER_ADDRESS} ${NODELIST}

echo "SUDO_SETUP_DONE"

exit 0