#!/bin/bash

. /pocket/conf/hadoop/2.4.0/cluster/conf.bashrc
$HADOOP_HOME/sbin/mr-jobhistory-daemon.sh stop historyserver
$HADOOP_HOME/sbin/stop-all.sh
