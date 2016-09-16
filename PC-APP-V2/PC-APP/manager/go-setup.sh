#!/bin/bash

# Exit if any command fails
set -e

# Environmental variables
GO_BUILD="${PROJECT_DIR}/../../PC-CORE/.build"

# Copy header file
cp "${GO_BUILD}/_cgo_export.h" "${SRCROOT}/manager/Application"

# Copy PC-CORE binary
cp "${GO_BUILD}/pc-core.a" "${SRCROOT}/manager/"

echo "GO PC-CORE copy done!"