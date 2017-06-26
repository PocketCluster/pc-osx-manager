#!/bin/bash

# Exit if any command fails
set -e

# Figure out where things are coming from and going to
export GOROOT="${HOME}/Workspace/POCKETPKG/DEPREPO/GOARCHIVE/go-1.7.6"
export GOREPO=${GOREPO:-"${HOME}/Workspace/POCKETPKG"}
export GOWORKPLACE=${GOWORKPLACE:-"${HOME}/Workspace/GOPLACE"}
export GOPATH="${GOREPO}:${GOWORKPLACE}"
export GO=${GOROOT}/bin/go
export GG_BUILD="${PWD}/../../.build"
export ARCHIVE="${GG_BUILD}/pc-core.a"
#PATH=${PATH:-"$GEM_HOME/ruby/2.0.0/bin:$HOME/.util:$GOROOT/bin:$GOREPO/bin:$GOWORKPLACE/bin:$HOME/.util:$NATIVE_PATH"}
export PATH="$GOROOT/bin:$GOREPO/bin:$GOWORKPLACE/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
export VERBOSE=${VERBOSE:-0}

# Clean old directory
if [ -d ${GG_BUILD} ]; then
    rm -rf ${GG_BUILD} && mkdir -p ${GG_BUILD}
fi

echo "--- --- --- --- --- --- --- --- --- --- --- --- GO ENVIRONMENTS --- --- --- --- --- --- --- --- --- --- --- ---"
echo $(GO version)
GO env
echo "--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---"

echo "Make the temp folders for go objects"
mkdir -p ${GG_BUILD}

echo "Generate _cgo_export.h and copy into source folder"
${GO} tool cgo -objdir ${GG_BUILD} native_*.go main.go

echo "Compile and produce object files"
# [Default mode] First trial
#CGO_ENABLED=1 CC=clang ${GO} build -ldflags '-tmpdir '${GG_BUILD}' -linkmode external' ./...

# [Default mode] External clang linker
#CGO_ENABLED=1 CC=clang ${GO} build -v -x -ldflags '-v -tmpdir '${GG_BUILD}' -linkmode external -extld clang' ./...

# [Archive mode]
#CGO_ENABLED=1 CC=clang ${GO} build -v -x -buildmode=c-archive -ldflags '-v -tmpdir '${GG_BUILD}' -linkmode external' ./...

# [Shared mode] go.dwarf file
#CGO_ENABLED=1 CC=clang ${GO} build -v -x -buildmode=c-shared -ldflags '-v -tmpdir '${GG_BUILD}' -linkmode external' ./...

# [Archive mode] prevents go.dwarf generated (-w), strip symbol (-s)
#CGO_ENABLED=1 CC=clang ${GO} build -v -x -buildmode=c-archive -ldflags '-v -w -s -tmpdir '${GG_BUILD}' -linkmode external' ./...

# [Default mode] default mode (we need main() function), disable go.dwarf generation (-w), strip symbol (-s)
if [[ ${VERBOSE} -eq 1 ]]; then
    CGO_ENABLED=1 CC=clang ${GO} build -v -x -ldflags '-v -s -w -tmpdir '${GG_BUILD}' -linkmode external' ./...
else
    CGO_ENABLED=1 CC=clang ${GO} build -ldflags '-s -w -tmpdir '${GG_BUILD}' -linkmode external' ./...
fi

echo "Combine the object files into a static library"
ar rcs ${ARCHIVE} ${GG_BUILD}/*.o
mv ${GG_BUILD}/_cgo_export.h ${GG_BUILD}/pc-core.h
rm static*
echo "${ARCHIVE} generated!"
