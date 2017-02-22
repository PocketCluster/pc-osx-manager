#!/usr/bin/env bash

:<<COMMENT
The *simplest* way to resize disk partition is to do with sfdisk

    sfdisk -d /dev/sdx > partition.layout
    ... do whatever you want ...
    sfdisk /dev/sdx < partition.layout
    resize2fs /dev/sdx2

COMMENT


# We'll create 32mb size image
SEEK=32
SIZE=$((16 * 2048))
IMAGE="./disk.img"

echo "This is to create a disk image which we can manipulate and test"

if [ $(id -u) -ne 0 ]; then
    echo "ERROR! Must be root."
    exit 1
fi

rm ${IMAGE}

# seek images
dd if=/dev/zero of=${IMAGE} bs=1M count=1
dd if=/dev/zero of=${IMAGE} bs=1M count=0 seek=${SEEK}

# set partition
sfdisk -f ${IMAGE} <<EOM
unit: sectors
1 : start=     2048, size=    30720, Id= c, bootable
2 : start=    32768, size=  ${SIZE}, Id=83
3 : start=        0, size=        0, Id= 0
4 : start=        0, size=        0, Id= 0
EOM
