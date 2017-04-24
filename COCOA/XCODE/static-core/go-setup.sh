#!/bin/bash

# Exit if any command fails
set -e

#pushd ${PWD}
#cd ${PROJECT_DIR}/../../../PC-APP-V2/pc-core/static/ && source ./build.sh
#popd

# Environmental variables
GO_BUILD="${PROJECT_DIR}/../../../PC-APP-V2/.build/"

echo "Project directory : ${PROJECT_DIR}"
echo "Project source : ${SRCROOT}"
echo "Go binary destination : ${GO_BUILD}"

# Copy header to goheader
if [[ -f "${PROJECT_DIR}/../../goheader/pc-core.h" ]]; then
    rm "${PROJECT_DIR}/../../goheader/pc-core.h"
fi
cp "${GO_BUILD}/pc-core.h" "${PROJECT_DIR}/../../goheader/"

# Copy PC-CORE binary
if [[ -f "${PROJECT_DIR}/static-core/pc-core.a" ]]; then
    rm "${PROJECT_DIR}/static-core/pc-core.a"
fi
cp "${GO_BUILD}/pc-core.a" "${PROJECT_DIR}/static-core/"

echo "GO PC-CORE copied!"