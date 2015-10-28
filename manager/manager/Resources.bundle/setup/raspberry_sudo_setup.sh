#!/bin/bash

BASE_BUNDLE_PATH=$1
#VAGRANT_PATH="$(which vagrant)"

echo "SUDO_SETUP_STEP_0"

# setup root directories
mkdir -p /{pocket,bigpkg}

# setup basic directories
mkdir -p /etc/salt
mkdir -p /pocket/{boxes,conf,hdfs,log,nodes,salt}
mkdir -p /pocket/nodes/pc-node{1..3}

# copy salt essential files
cp -Rf "${BASE_BUNDLE_PATH}"/saltstack/* /pocket/salt/

# copy vagrant essential fils
cp -Rf "${BASE_BUNDLE_PATH}"/vagrant/* /pocket/boxes/

# copy vagrant files & modify
chown -R $SUDO_USER:staff /pocket
chown -R $SUDO_USER:staff /bigpkg

# copy salt config files
cp -f "${BASE_BUNDLE_PATH}"/etc/salt/* /etc/salt/
sed -i '' 's|PC_USER|'$SUDO_USER'|g' /etc/salt/*

# change hosts
python "${BASE_BUNDLE_PATH}"/setup/host_setup.py salt 10.211.55.1 pc-master 10.211.55.1 pc-node1 10.211.55.201 pc-node2 10.211.55.202 pc-node3 10.211.55.203

echo "SUDO_SETUP_DONE"

exit 0