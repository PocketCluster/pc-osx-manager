#!/bin/bash

. /pocket/conf/hadoop/2.4.0/cluster/conf.bashrc
. /pocket/conf/spark/1.5.2/standalone/conf.bashrc

$SPARK_HOME/sbin/stop-all.sh