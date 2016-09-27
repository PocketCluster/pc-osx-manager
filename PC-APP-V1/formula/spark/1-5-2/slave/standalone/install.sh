#!/bin/bash

# config path
mkdir -p /pocket/conf/spark/1.5.2/standalone

# log file path
mkdir -p /pocket/log/spark/1.5.2/standalone

# intermediate dir path
mkdir -p /pocket/interim/spark/1.5.2/standalone/work
#mkdir -p /pocket/interim/spark/1.5.2/standalone/metastore_db
chmod 777 /pocket/interim/spark/1.5.2/
chown -R pocket:pocket /pocket/

# download package
if [ ! -d "/bigpkg/spark-1.5.2-bin-hadoop2.4" ] ; then
	wget -qO- "http://pc-master:10120/spark-1.5.2-bin-hadoop2.4.tgz" | tar xvz -C /bigpkg 2>&1
	chown -R pocket:pocket "/bigpkg/spark-1.5.2-bin-hadoop2.4"
fi
