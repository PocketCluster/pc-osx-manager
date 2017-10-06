package slcontext

import (
    log "github.com/Sirupsen/logrus"

    "github.com/gravitational/teleport/lib/auth"
)

func UserIdentityPostProcess(uinfo *auth.PocketResponseUserIdentity) error {
    log.Debugf("login %v | uid %v | gid %v", uinfo.LoginName, uinfo.UID, uinfo.GID)
    return nil
}