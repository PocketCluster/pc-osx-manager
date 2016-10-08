#!/bin/bash

# config path
mkdir -p /pocket/conf/hadoop/2.4.0/cluster

# namenode path
mkdir -p /pocket/hdfs/2.4.0/cluster/namenode
mkdir -p /pocket/hdfs/2.4.0/cluster/namenode2-checkpoint
mkdir -p /pocket/hdfs/2.4.0/cluster/datanode
mkdir -p /pocket/hdfs/2.4.0/cluster/yarn/nm-local-dir
mkdir -p /pocket/hdfs/2.4.0/cluster/yarn/nm-log-dir/userlogs

# log file path
mkdir -p /pocket/log/hadoop/2.4.0/cluster

# save package to archive
if [ ! -f "/bigpkg/archive/hadoop-2.4.0.tar.gz" ] ; then
	curl -o "/bigpkg/archive/hadoop-2.4.0.tar.gz" "https://archive.apache.org/dist/hadoop/core/hadoop-2.4.0/hadoop-2.4.0.tar.gz" 
fi

if [ ! -f "/bigpkg/hadoop-2.4.0/bin/hdfs" ] ; then
	tar -xvzf "/bigpkg/archive/hadoop-2.4.0.tar.gz" -C "/bigpkg/"
fi

exit 0