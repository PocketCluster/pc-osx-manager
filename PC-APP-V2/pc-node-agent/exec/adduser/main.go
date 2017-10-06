package main

import (
    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/gravitational/teleport/lib/auth"
)

func main() {
    log.SetLevel(log.DebugLevel)

    slcontext.UserIdentityPostProcess(&auth.PocketResponseUserIdentity{
        LoginName: "stkim1",
        UID:       "501",
        GID:       "20",
    })
}
