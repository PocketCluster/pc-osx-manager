#!/bin/bash

. /pocket/conf/hadoop/2.4.0/cluster/conf.bashrc
. /pocket/conf/spark/1.5.2/standalone/conf.bashrc

HADOOP_STATUS="$(jps | grep NameNode)"

for ((i=1;i<=${NUM_NODES};i++));
do
	echo "pc-node${i}" >> "${SPARK_CONF_DIR}"/slaves
done

if [ -z "${HADOOP_STATUS}" ] ; 
then
	$HADOOP_HOME/sbin/start-dfs.sh
	$HADOOP_HOME/bin/hdfs dfsadmin -safemode wait
	$HADOOP_HOME/bin/hdfs dfs -mkdir -p /sparklog
	$HADOOP_HOME/bin/hdfs dfs -chmod -R 1777 /sparklog
	$HADOOP_HOME/sbin/stop-dfs.sh
else
    $HADOOP_HOME/bin/hadoop dfsadmin -safemode wait
	$HADOOP_HOME/bin/hdfs dfs -mkdir -p /sparklog
	$HADOOP_HOME/bin/hdfs dfs -chmod -R 1777 /sparklog
fi
