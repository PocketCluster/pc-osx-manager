package slcontext

import (
    "strings"
    "os/exec"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/gravitational/teleport/lib/auth"
)

func UserIdentityPostProcess(uinfo *auth.PocketResponseUserIdentity) error {
    if len(uinfo.LoginName) == 0 {
        return errors.Errorf("invalid login name")
    }
    if len(uinfo.UID) == 0 {
        return errors.Errorf("invalid uid")
    }
    if len(uinfo.GID) == 0 {
        return errors.Errorf("invalid gid")
    }

    cmdtmpl := `/usr/sbin/addgroup --gid [UID] [USERLOGIN];/usr/sbin/adduser --disabled-login --shell /usr/sbin/nologin --gid [UID] --uid [UID] [USERLOGIN];/usr/bin/passwd -l [USERLOGIN]`
    cmdtmpl = strings.Replace(cmdtmpl, "[USERLOGIN]", uinfo.LoginName, -1)
    cmdtmpl = strings.Replace(cmdtmpl, "[UID]", uinfo.UID, -1)
    cmds := strings.Split(cmdtmpl, ";")

    for _, c := range cmds {
        clist := strings.Split(c," ")
        cmd := exec.Command(clist[0], clist[1:]...)
        out, err := cmd.Output()
        if err != nil {
            log.Error(err.Error())
        }
        log.Debug(string(out))
    }

    return nil
}