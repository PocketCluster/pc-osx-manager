#!/bin/bash

# config path
mkdir -p /pocket/conf/spark/1.5.2/standalone

# log file path
mkdir -p /pocket/log/spark/1.5.2/standalone

# intermediate dir path
mkdir -p /pocket/interim/spark/1.5.2/standalone/work
# mkdir -p /pocket/interim/spark/1.5.2/standalone/metastore_db
chmod 777 /pocket/interim/spark/1.5.2/

# homebrew formula update
# brew update && brew upgrade

# install scala
brew update; brew switch scala 2.11.7; brew install scala

# install r-base is removed for now. this will be re-enabled in the future
# brew tap homebrew/science && brew install r && brew untap homebrew/science 


# save package to archive
if [ ! -f "/bigpkg/archive/spark-1.5.2-bin-hadoop2.4.tgz" ] ; then
	curl -o "/bigpkg/archive/spark-1.5.2-bin-hadoop2.4.tgz" "http://apache.mirror.cdnetworks.com/spark/spark-1.5.2/spark-1.5.2-bin-hadoop2.4.tgz"
fi

if [ ! -f "/bigpkg/spark-1.5.2-bin-hadoop2.4/bin/spark-shell" ] ; then
	tar -xvzf "/bigpkg/archive/spark-1.5.2-bin-hadoop2.4.tgz" -C "/bigpkg/"
fi

exit 0