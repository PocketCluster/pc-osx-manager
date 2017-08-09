#!/usr/bin/env bash

# Figure out where things are coming from and going to
export GOROOT="${HOME}/Workspace/POCKETPKG/DEPREPO/GOARCHIVE/go-1.7.6"
export GOREPO=${GOREPO:-"${HOME}/Workspace/POCKETPKG"}
export GOPATH=${GOPATH:-"${GOREPO}:${GOWORKPLACE}"}
export PATH="$GOROOT/bin:$GOREPO/bin:$GOWORKPLACE/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
echo $(go version)
