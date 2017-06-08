package disk

import (
    "io/ioutil"
    "os"
    "os/exec"
    "path"
    "strings"

    "github.com/pkg/errors"
)

func activateSwap(partition string) error {
    cmd := exec.Command("/sbin/mkswap", partition)
    err := cmd.Start()
    if err != nil {
        return errors.WithStack(err)
    }
    err = cmd.Wait()
    if err != nil {
        return errors.WithStack(err)
    }

    cmd = exec.Command("/sbin/swapon", partition)
    err = cmd.Start()
    if err != nil {
        return errors.WithStack(err)
    }
    return cmd.Wait()
}

func ActivateSystemSwapParition() error {
    return activateSwap("/dev/mmcblk0p3")
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