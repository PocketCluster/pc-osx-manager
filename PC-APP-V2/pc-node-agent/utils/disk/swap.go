package disk

import (
    "io/ioutil"
    "os"
    "path"
    "strings"

    "github.com/pkg/errors"
)

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

func AppendSwapPartitionToTable(rootPath string) error {
    var (
        sysFsTbl = rootPath + "/etc/fstab"
    )
    if !path.IsAbs(sysFsTbl) {
        return errors.Errorf("[ERR] file system is not absolute")
    }
    fstab, err := ioutil.ReadFile(sysFsTbl)
    if err != nil {
        return errors.WithStack(err)
    }

    var fsEntries []string = []string{}
    for _, l := range strings.Split(string(fstab), "\n") {
        if !strings.HasPrefix(l, "/dev/mmcblk0p3") {
            nl := strings.TrimSpace(l)
            if len(nl) != 0 {
                fsEntries = append(fsEntries, nl)
            }
        }
    }
    fsEntries = append(fsEntries, "/dev/mmcblk0p3    none    swap    sw    0    0\n")

    var newFstab string = strings.Join(fsEntries, "\n")
    err = ioutil.WriteFile(sysFsTbl, []byte(newFstab), os.FileMode(0644))
    return errors.WithStack(err)
}