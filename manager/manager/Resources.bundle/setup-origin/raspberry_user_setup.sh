#!/bin/bash

BASE_BUNDLE_PATH=$1
NUM_NODES=$2

echo "USER_SETUP_STEP_0"

mkdir -p /pocket/{conf,log,salt}

# copy salt essential files
cp -Rf "${BASE_BUNDLE_PATH}"/saltstack/* /pocket/salt/

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
if ! grep -Fxq "${LOC}" $HOME/.ssh/known_hosts
then
    echo "${LOC}" >> $HOME/.ssh/known_hosts
fi

# pc-master key
PM="$(ssh-keyscan -t rsa pc-master)"
if ! grep -Fxq "${PM}" $HOME/.ssh/known_hosts
then
    echo "${PM}" >> $HOME/.ssh/known_hosts
fi

# pc-node# key
for ((i=1;i<=${NUM_NODES};i++));
do
    PN="$(ssh-keyscan -t rsa pc-node${i})"
    if ! grep -Fxq "${PN}" $HOME/.ssh/known_hosts
    then
        echo "${PN}" >> $HOME/.ssh/known_hosts
    fi
done

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

# pc-node{1..${NUM_NODES}}
for ((i=1;i<=${NUM_NODES};i++));
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

salt-key -y --accept="pc-master"

sleep 1

for ((i=1;i<=${NUM_NODES};i++));
do
    salt-key -y --accept="pc-node${i}"
    sleep 1
done

salt 'pc-node*' state.sls base/setup
salt 'pc-node*' state.sls base/ssh-login

echo "USER_SETUP_DONE"

exit 0
