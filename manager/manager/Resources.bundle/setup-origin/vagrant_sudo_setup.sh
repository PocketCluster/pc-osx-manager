#!/bin/bash

BASE_BUNDLE_PATH=$1
#VAGRANT_PATH="$(which vagrant)"

echo "SUDO_SETUP_STEP_0"

# setup root directories
mkdir -p /{pocket,bigpkg}

# change permission
chown -R $SUDO_USER:staff /pocket
chown -R $SUDO_USER:staff /bigpkg
chown $SUDO_USER:admin /usr/local

# copy & modify salt config files
mkdir -p /etc/salt
cp -f "${BASE_BUNDLE_PATH}"/etc/salt/* /etc/salt/
sed -i '' 's|PC_USER|'$SUDO_USER'|g' /etc/salt/*

# change hosts
python "${BASE_BUNDLE_PATH}"/setup/host_setup.py salt 10.211.55.1 pc-master 10.211.55.1 pc-node1 10.211.55.201 pc-node2 10.211.55.202 pc-node3 10.211.55.203

# cd /pocket/boxes && sudo -u $SUDO_USER $VAGRANT_PATH up 2>&1

echo "SUDO_SETUP_DONE"

exit 0