# add the following to the end of file

export JAVA_HOME="$(/usr/libexec/java_home -v 1.8)"

export HADOOP_HOME="/bigpkg/hadoop-2.4.0"
export HADOOP_PREFIX=$HADOOP_HOME
export HADOOP_INSTALL=$HADOOP_HOME
export HADOOP_MAPRED_HOME=$HADOOP_HOME
export HADOOP_COMMON_HOME=$HADOOP_HOME
export HADOOP_HDFS_HOME=$HADOOP_HOME
export HADOOP_COMMON_LIB_NATIVE_DIR=$HADOOP_HOME/lib/native

export HADOOP_CONF_DIR=/pocket/conf/hadoop/2.4.0/cluster
export YARN_HOME=$HADOOP_HOME
export YARN_CONF_DIR=$HADOOP_CONF_DIR

export PATH=$PATH:$HADOOP_HOME/bin:$HADOOP_HOME/sbin