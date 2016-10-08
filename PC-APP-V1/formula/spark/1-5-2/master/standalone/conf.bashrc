# add the following to the end of file

export JAVA_HOME="$(/usr/libexec/java_home -v 1.8)"
export SCALA_PATH="$(type -p scala)"
SCALA_HOME=""

if [[ -n $SCALA_PATH ]] && [[ -x $SCALA_PATH ]]; 
then
	export SCALA_HOME=${SCALA_PATH/scala/}
fi

export SPARK_HOME=/bigpkg/spark-1.5.2-bin-hadoop2.4
export SPARK_CONF_DIR=/pocket/conf/spark/1.5.2/standalone

export PATH=$PATH:$SPARK_HOME/bin
