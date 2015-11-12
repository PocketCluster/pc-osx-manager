#!/bin/bash

echo "USER_SETUP_STEP_2"

salt-key -y --accept="pc-master"
sleep 1

for i in {1..3}
do
    salt-key -y --accept="pc-node${i}"
    sleep 1
done

salt 'pc-node*' state.sls base/setup
salt 'pc-node*' state.sls base/ssh-login
rm -rf /pocket/salt/states/base/ssh/*

echo "USER_SETUP_DONE"

exit 0