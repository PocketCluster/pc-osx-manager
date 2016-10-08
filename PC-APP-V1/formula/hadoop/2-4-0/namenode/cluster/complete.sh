#!/bin/bash

. /pocket/conf/hadoop/2.4.0/cluster/conf.bashrc

for ((i=1;i<=${NUM_NODES};i++));
do
	if [[ $i == 1 ]]; then
		echo "pc-node${i}" > "${HADOOP_CONF_DIR}"/slaves
	else
		echo "pc-node${i}" >> "${HADOOP_CONF_DIR}"/slaves
	fi
done

$HADOOP_HOME/bin/hadoop namenode -format
$HADOOP_HOME/sbin/start-dfs.sh
$HADOOP_HOME/bin/hdfs dfsadmin -safemode wait

$HADOOP_HOME/bin/hdfs dfs -mkdir -p /user/"${USER}"
$HADOOP_HOME/bin/hdfs dfs -mkdir -p /tmp
$HADOOP_HOME/bin/hdfs dfs -mkdir -p /jobhistory/tmp
$HADOOP_HOME/bin/hdfs dfs -mkdir -p /jobhistory/done

$HADOOP_HOME/bin/hdfs dfs -chmod -R 1777 /jobhistory/tmp
$HADOOP_HOME/bin/hdfs dfs -chmod -R 1777 /jobhistory/done

$HADOOP_HOME/sbin/stop-dfs.sh
