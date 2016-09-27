#!/usr/bin/env bash

export JAVA_HOME=${JAVA_HOME}
export HADOOP_CONF_DIR=/pocket/conf/hadoop/2.4.0/cluster
export SPARK_LOCAL_IP="{{ grains['nodename'] }}"
export SPARK_WORKER_INSTANCES=1
export SPARK_EXECUTOR_MEMORY=1g
export SPARK_WORKER_MEMORY=1g
export SPARK_WORKER_OPTS="$SPARK_WORKER_OPTS -Djava.awt.headless=true"
export SPARK_MASTER_IP="pc-master"
export SPARK_LOG_DIR=/pocket/log/spark/1.5.2/standalone
export SPARK_LOCAL_DIRS=/pocket/interim/spark/1.5.2/standalone/work