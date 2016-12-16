package pcauth

import (
    "fmt"
    "time"
    "io/ioutil"
    "strings"

    "github.com/gravitational/teleport/lib/auth"
    log "github.com/Sirupsen/logrus"

    "github.com/stkim1/pcteleport/pcconfig"
)

// enforceTokenTTL deletes the given token if it's TTL is over. Returns 'false'
// if this token cannot be used
func checkTokenTTL(s *auth.AuthServer, token string) bool {
    // look at the tokens in the token storage
    tok, err := s.Provisioner.GetToken(token)
    if err != nil {
        log.Warn(err)
        return true
    }
    // s.clock is replaced with time.Now()
    now := time.Now().UTC()
    if tok.Expires.Before(now) {
        if err = s.DeleteToken(token); err != nil {
            log.Error(err)
        }
        return false
    }
    return true
}

func readToken(token string) (string, error) {
    if !strings.HasPrefix(token, "/") {
        return token, nil
    }
    // treat it as a file
    out, err := ioutil.ReadFile(token)
    if err != nil {
        return "", nil
    }
    return string(out), nil
}

func writeDockerKeyAndCert(cfg *pcconfig.Config, keys *packedAuthKeyCert) error {
    log.Debugf("write slave docker auth to %v, key to %v, cert from %v", cfg.DockerAuthFile, cfg.DockerKeyFile, cfg.DockerCertFile)
    if err := ioutil.WriteFile(cfg.DockerAuthFile, keys.Auth, 0600); err != nil {
        return err
    }
    if err := ioutil.WriteFile(cfg.DockerKeyFile,  keys.Key, 0600); err != nil {
        return err
    }
    if err := ioutil.WriteFile(cfg.DockerCertFile, keys.Cert, 0600); err != nil {
        return err
    }
    return nil
}

func apiEndpoint(params ...string) string {
    return fmt.Sprintf("http://stub:0/%s/%s", PocketApiVersion, strings.Join(params, "/"))
}

