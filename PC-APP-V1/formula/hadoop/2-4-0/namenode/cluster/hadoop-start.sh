#!/bin/bash

. /pocket/conf/hadoop/2.4.0/cluster/conf.bashrc
$HADOOP_HOME/sbin/start-all.sh
$HADOOP_HOME/sbin/mr-jobhistory-daemon.sh start historyserver
