// +build linux
package utils

import (
    "os/exec"
)

func RebootNow() {
    cmd := exec.Command("/bin/systemctl", "reboot")
    cmd.Run()
}

func Shutdown() {
    cmd := exec.Command("/bin/systemctl", "poweroff")
    cmd.Run()
}
