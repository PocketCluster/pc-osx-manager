#!/bin/bash

TOTAL_NODES=3

while getopts t: option
do
	case "${option}"
	in
		t) TOTAL_NODES=${OPTARG};;
	esac
done

function setupHosts {
	echo "modifying /etc/hosts file"

	echo "127.0.0.1     localhost" >> /etc/nhosts

	echo "10.211.55.1   pc-master" >> /etc/nhosts
	echo "10.211.55.1   salt" >> /etc/nhosts

	for i in $(seq 1 $TOTAL_NODES)
	do 
		entry="10.211.55.20${i} pc-node${i}"
		echo "adding ${entry}"
		echo "${entry}" >> /etc/nhosts
	done

	mv /etc/nhosts /etc/hosts
}


function setupAccount {

    # Set up default user
    adduser --gecos "Pocket Cluster User" --add_extra_groups --disabled-password pocket
    usermod -a -G sudo,adm -p $(echo "pocket" | openssl passwd -1 -stdin) pocket
    echo "pocket ALL=(ALL) NOPASSWD:ALL" | tee "/etc/sudoers.d/pocket"
    chmod 440 /etc/sudoers.d/pocket
}

echo "setup node's hosts file..."
setupHosts

echo "setup PocketCluster account..."
setupAccount
