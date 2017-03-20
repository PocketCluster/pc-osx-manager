#!/bin/bash

# Exit if any command fails
set -e

# Environmental variables
GO_BUILD="${PROJECT_DIR}/../../../PC-APP-V2/.build/"

echo "Project directory : ${PROJECT_DIR}"
echo "Project source : ${SRCROOT}"
echo "Go binary destination : ${GO_BUILD}"

# Copy header file
cp "${GO_BUILD}/_cgo_export.h" "${PROJECT_DIR}/static-core/"

# Copy PC-CORE binary
if [[ -f "${PROJECT_DIR}/static-core/pc-core.a" ]]; then
    rm "${PROJECT_DIR}/static-core/pc-core.a"
fi
cp "${GO_BUILD}/pc-core.a" "${PROJECT_DIR}/static-core/"

echo "GO PC-CORE copied!"