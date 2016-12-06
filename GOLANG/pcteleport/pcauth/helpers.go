package pcauth

import (
    "fmt"
    "time"
    "io/ioutil"
    "strings"
    "path/filepath"

    "github.com/gravitational/teleport/lib/auth"
    log "github.com/Sirupsen/logrus"
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

// keysPath returns two full file paths: to the host.key and host.cert
func keysPath(dataDir string, id auth.IdentityID) (key string, cert string) {
    return filepath.Join(dataDir, fmt.Sprintf("%s.key", strings.ToLower(string(id.Role)))),
        filepath.Join(dataDir, fmt.Sprintf("%s.cert", strings.ToLower(string(id.Role))))
}

// writeKeys saves the key/cert pair for a given domain onto disk. This usually means the
// domain trusts us (signed our public key)
func writeKeys(dataDir string, id auth.IdentityID, key []byte, cert []byte) error {
    kp, cp := keysPath(dataDir, id)
    log.Debugf("write key to %v, cert from %v", kp, cp)

    if err := ioutil.WriteFile(kp, key, 0600); err != nil {
        return err
    }
    if err := ioutil.WriteFile(cp, cert, 0600); err != nil {
        return err
    }
    return nil
}

func apiEndpoint(params ...string) string {
    return fmt.Sprintf("http://stub:0/%s/%s", PocketApiVersion, strings.Join(params, "/"))
}

