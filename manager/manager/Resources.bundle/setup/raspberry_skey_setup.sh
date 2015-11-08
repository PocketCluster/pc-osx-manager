#!/bin/bash

BASE_BUNDLE_PATH=$1
NUM_NODES=$2

echo "USER_SETUP_STEP_2"

salt-key -y --accept="pc-master"

sleep 1

for ((i=1;i<=${NUM_NODES};i++));
do
    salt-key -y --accept="pc-node${i}"
    sleep 1
done

salt 'pc-node*' state.sls base/setup
salt 'pc-node*' state.sls base/ssh-login
salt 'pc-node*' cmd.run 'sh /repartition.sh'
salt 'pc-node*' cmd.run 'reboot'

sleep 3

echo "USER_SETUP_DONE"

exit 0
