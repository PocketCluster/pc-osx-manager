package status

import (
    "os"
)

func HostName() (string, error) {
    return os.Hostname()
}


