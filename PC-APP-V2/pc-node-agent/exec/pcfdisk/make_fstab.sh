#!/usr/bin/env bash

rm -rf /tmp/etc/ && mkdir -p /tmp/etc/
cat <<EOT >> /tmp/etc/fstab
proc            /proc           proc    defaults          0       0
/dev/mmcblk0p2  /               ext4   defaults,noatime   0       1
/dev/mmcblk0p1  /boot/          vfat    defaults          0       2
/dev/sda1       /work           ext4   defaults,noatime   0       1
EOT