#!/bin/bash

# Exit if any command fails
set -e

# Figure out where things are coming from and going to
GO=${GOROOT}/bin/go
GG_BUILD="${PWD}/../../.build"
ARCHIVE="${GG_BUILD}/pc-core.a"

# Clean old directory
if [ -d ${GG_BUILD} ]; then
    rm -rf ${GG_BUILD} && mkdir -p ${GG_BUILD}
fi

echo "Make the temp folders for go objects"
mkdir -p ${GG_BUILD}

echo "Generate _cgo_export.h and copy into source folder"
${GO} tool cgo -objdir ${GG_BUILD} *.go

echo "Compile and produce object files"
CGO_ENABLED=1 CC=clang ${GO} build -ldflags '-tmpdir '${GG_BUILD}' -linkmode external' ./...

echo "Combine the object files into a static library"
ar rcs ${ARCHIVE} ${GG_BUILD}/*.o
rm static
echo "${ARCHIVE} generated!"
