package disk

func ActivateSwapPartition() error {
/*
# make swap space
mkswap /dev/mmcblk0p3

# turn swap space
swapon /dev/mmcblk0p3

# you can use UUID
blkid /dev/mmcblk0p3

vi /etc/fstab
# add the following line
/dev/mmcblk0p3  none    swap    sw                  0       0
 */
    return nil
}