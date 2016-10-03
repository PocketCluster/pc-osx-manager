package status

import (
    "runtime"
)

func CPUCount() int {
    return runtime.NumCPU()
}

