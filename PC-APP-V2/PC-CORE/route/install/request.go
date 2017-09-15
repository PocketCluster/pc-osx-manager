package install

import (
    "io/ioutil"
    "net/http"
    "time"

    "github.com/pkg/errors"
)

func newRequest(url string, isBinaryReq bool) (*http.Request, error) {
    req, err :=  http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    //req.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
    req.Header.Add("User-Agent", "PocketCluster/0.1.4 (OSX)")
    if isBinaryReq {
        req.Header.Set("Content-Type", "application/octet-stream")
    } else {
        req.Header.Set("Content-Type", "application/json; charset=utf-8")
    }
    req.ProtoAtLeast(1, 1)
    return req, nil
}

func newClient(timeout time.Duration, noCompress bool) *http.Client {
    return &http.Client {
        Timeout: timeout,
        Transport: &http.Transport {
            DisableCompression: noCompress,
        },
    }
}

func readRequest(req *http.Request, client *http.Client) ([]byte, error) {
    resp, err := client.Do(req)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        return nil, errors.Errorf("protocol status : %d", resp.StatusCode)
    }
    return ioutil.ReadAll(resp.Body)
}
