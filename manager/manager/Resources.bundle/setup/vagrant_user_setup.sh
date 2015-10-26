#!/bin/bash

BASE_BUNDLE_PATH=$1
VAGRANT_PATH="$(which vagrant)"

echo "USER_SETUP_STEP_0"

# vagrant up
#cd /pocket/boxes && $USER $VAGRANT_PATH up 2>&1

echo "USER_SETUP_STEP_1"

# setup ssh login
if [ ! -d "$HOME/.ssh" ]; then
    mkdir -p $HOME/.ssh
fi

#cd $HOME/.ssh
if [[ ! -f "$HOME/.ssh/id_rsa" ]] && [[ ! -f "$HOME/.ssh/id_rsa.pub" ]];then
    cat /dev/zero | ssh-keygen -t rsa -P ""
fi

PK=$(<"$HOME/.ssh/id_rsa.pub")

if ! grep -Fxq "$PK" $HOME/.ssh/authorized_keys
then
    cat $HOME/.ssh/id_rsa.pub >> $HOME/.ssh/authorized_keys
fi

# localhost key
LOC="$(ssh-keyscan -t rsa localhost)"
if ! grep -Fxq "$LOC" $HOME/.ssh/known_hosts
then
    echo "${LOC}" >> $HOME/.ssh/known_hosts
fi

# pc-node1 key
PN1="$(ssh-keyscan -t rsa pc-node1)"
if ! grep -Fxq "${PN1}" $HOME/.ssh/known_hosts
then
    echo "${PN1}" >> $HOME/.ssh/known_hosts
fi

# pc-node2 key
PN2="$(ssh-keyscan -t rsa pc-node2)"
if ! grep -Fxq "${PN2}" $HOME/.ssh/known_hosts
then
    echo "${PN2}" >> $HOME/.ssh/known_hosts
fi

# pc-node3 key
PN3="$(ssh-keyscan -t rsa pc-node3)"
if ! grep -Fxq "${PN3}" $HOME/.ssh/known_hosts
then
    echo "${PN3}" >> $HOME/.ssh/known_hosts
fi

# config
if [[ ! -f "$HOME/.ssh/config" ]]; then
    touch "$HOME/.ssh/config"
fi

# pc-master
if ! grep -q "Host pc-master" $HOME/.ssh/config
then
    echo "\nHost pc-master" >> $HOME/.ssh/config
    echo "\tHostName pc-master" >> $HOME/.ssh/config
    echo "\tUser $USER" >> $HOME/.ssh/config
    echo "\tIdentityFile ~/.ssh/id_rsa" >> $HOME/.ssh/config
fi

# pc-node{1..3}
for i in {1..3}
do
    if ! grep -q "Host pc-node${i}" $HOME/.ssh/config
    then
        echo "\nHost pc-node${i}" >> $HOME/.ssh/config
        echo "\tHostName pc-node${i}" >> $HOME/.ssh/config
        echo "\tUser pocket" >> $HOME/.ssh/config
        echo "\tIdentityFile ~/.ssh/id_rsa" >> $HOME/.ssh/config
    fi
done

chmod 700 $HOME/.ssh
chmod 600 $HOME/.ssh/*

cp -f $HOME/.ssh/* /pocket/salt/states/base/ssh/

echo "USER_SETUP_STEP_2"

salt-key -a "pc-master"
salt-key -a "pc-node1"
salt-key -a "pc-node2"
salt-key -a "pc-node3"

echo "USER_SETUP_DONE"

exit 0
