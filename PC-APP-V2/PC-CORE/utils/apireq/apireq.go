package apireq

import (
    "io/ioutil"
    "net/http"
    "time"

    "github.com/pkg/errors"
)

const (
    ConnTimeout = time.Duration(5 * time.Second)
)

func NewRequest(url string, isBinaryReq bool) (*http.Request, error) {
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

func NewClient(timeout time.Duration, noCompress bool) *http.Client {
    return &http.Client {
        Timeout: timeout,
        Transport: &http.Transport {
            DisableCompression: noCompress,
        },
    }
}

func ReadRequest(req *http.Request, client *http.Client) ([]byte, error) {
    resp, err := client.Do(req)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    if resp.StatusCode == 200 {
        if resp.Body == nil {
            return nil, errors.Errorf("invalid null service response")
        } else {
            defer resp.Body.Close()
            return ioutil.ReadAll(resp.Body)
        }
    }

    if resp.Body == nil {
        return nil, errors.Errorf("service response with code %d", resp.StatusCode)
    } else {
        defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
        return body, errors.Errorf("service response with code %d", resp.StatusCode)
    }
}
