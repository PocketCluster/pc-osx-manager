package main

import (
    "os"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-core/defaults"
)

// TODO : WE NEED UNIFIED LOGGING FACILITY (Master + Slave)
func setLogger(debug bool) {
    log.StandardLogger().Hooks = make(log.LevelHooks)
    log.SetFormatter(&log.TextFormatter{
        // Let's enable color for now
        //DisableColors:      true,
        TimestampFormat:    defaults.PocketTimeDateFormat,
    })
    if debug {
        log.SetLevel(log.DebugLevel)
    } else {
        log.SetLevel(log.InfoLevel)
    }
    // directing where logs to go
    log.SetOutput(os.Stderr)
}

