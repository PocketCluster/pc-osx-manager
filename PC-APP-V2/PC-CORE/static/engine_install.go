package main

import (
    "bytes"
    "encoding/binary"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/Redundancy/go-sync"
    //"github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/route/routepath"
    "github.com/stkim1/pc-core/model"
)

// reads the file headers and checks the magic string, then the semantic versioning
// return : in order of 'filesize', 'blocksize', 'blockcount', 'rootHash', 'error'
func readHeadersAndCheck(r io.Reader) (int64, uint32, uint32, []byte, error) {
    var (
        bMagic                []byte = make([]byte, len(gosync.PocketSyncMagicString))
        major, minor, patch   uint16 = 0, 0, 0
        filesize              int64  = 0
        blocksize, blockcount uint32 = 0, 0
        hLen                  uint32 = 0
        rootHash              []byte = nil
    )
    // magic string
    if _, err := r.Read(bMagic); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    } else if string(bMagic) != gosync.PocketSyncMagicString {
        return 0, 0, 0, nil, errors.New("meta header does not confirm. Not a valid meta")
    }

    // version
    for _, v := range []*uint16{&major, &minor, &patch} {
        if err := binary.Read(r, binary.LittleEndian, v); err != nil {
            return 0, 0, 0, nil, errors.WithStack(err)
        }
    }
    if major != gosync.PocketSyncMajorVersion || minor != gosync.PocketSyncMinorVersion || patch != gosync.PocketSyncPatchVersion {
        return 0, 0, 0, nil, errors.Errorf("The acquired version (%v.%v.%v) does not match the tool (%v.%v.%v).",
            major, minor, patch,
            gosync.PocketSyncMajorVersion, gosync.PocketSyncMinorVersion, gosync.PocketSyncPatchVersion)
    }

    if err := binary.Read(r, binary.LittleEndian, &filesize); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    if err := binary.Read(r, binary.LittleEndian, &blocksize); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    if err := binary.Read(r, binary.LittleEndian, &blockcount); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    if err := binary.Read(r, binary.LittleEndian, &hLen); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    rootHash = make([]byte, hLen)
    if _, err := r.Read(rootHash); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    return filesize, blocksize, blockcount, rootHash, nil
}

func initInstallRoutePath() {
    const (
        timeout = time.Duration(10 * time.Second)
    )
    var (
        newRequest = func(url string, isBinaryReq bool) (*http.Request, error) {
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
        newClient = func(timeout time.Duration, noCompress bool) *http.Client {
            return &http.Client {
                Timeout: timeout,
                Transport: &http.Transport {
                    DisableCompression: noCompress,
                },
            }
        }
        readRequest = func(req *http.Request, client *http.Client) ([]byte, error) {
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
    )

    // get the list of available packages
    theApp.GET(routepath.RpathPackageList(), func(_, rpath, _ string) error {
        var (
            feedError = func(errMsg string) error {
                data, err := json.Marshal(ReponseMessage{
                    "package-list": {
                        "status": false,
                        "error" : errMsg,
                    },
                })
                if err != nil {
                    log.Debugf(err.Error())
                    return errors.WithStack(err)
                }
                err = FeedResponseForGet(rpath, string(data))
                return errors.WithStack(err)
            }

            pkgList = []map[string]interface{}{}
        )

        req, err :=  newRequest("https://api.pocketcluster.io/service/v014/package/list", false)
        if err != nil {
            log.Debugf(errors.WithStack(err).Error())
            return feedError("Unable to access package list. Reason : " + errors.WithStack(err).Error())
        }
        client := newClient(timeout, true)
        resp, err := client.Do(req)
        if err != nil {
            log.Debugf(errors.WithStack(err).Error())
            return feedError("Unable to access package list. Reason : " + errors.WithStack(err).Error())
        }
        defer resp.Body.Close()

        if resp.StatusCode != 200 {
            return feedError(errors.Errorf("Service unavailable. Status : %d", resp.StatusCode).Error())
        }
        var pkgs = []*model.Package{}
        err = json.NewDecoder(resp.Body).Decode(&pkgs)
        if err != nil {
            log.Debugf(errors.WithStack(err).Error())
            return feedError("Unable to access package list. Reason : " + errors.WithStack(err).Error())
        }
        if len(pkgs) == 0 {
            return feedError("No package avaiable. Contact us at Slack channel.")
        } else {
            // update package doesn't return error when there is packages to update.
            model.UpdatePackages(pkgs)
        }

        for i, _ := range pkgs {
            pkgList = append(pkgList, map[string]interface{} {
                "package-id" : pkgs[i].PkgID,
                "description": pkgs[i].Description,
                "installed": false,
            })
        }
        data, err := json.Marshal(ReponseMessage{
            "package-list": {
                "status": true,
                "list":   pkgList,
            },
        })
        if err != nil {
            log.Debugf(err.Error())
            return errors.WithStack(err)
        }
        err = FeedResponseForGet(rpath, string(data))
        if err != nil {
            return errors.WithStack(err)
        }
        return nil
    })


    // install a package
    theApp.POST(routepath.RpathPackageInstall(), func(_, rpath, payload string) error {
        var (
            feedError = func(errMsg string) error {
                log.Errorf(errMsg)
                data, err := json.Marshal(ReponseMessage{
                    "package-install": {
                        "status": false,
                        "error" : errMsg,
                    },
                })
                if err != nil {
                    log.Debugf(err.Error())
                    return errors.WithStack(err)
                }
                err = FeedResponseForPost(rpath, string(data))
                return errors.WithStack(err)
            }
            client = newClient(timeout, false)
            repoList = []string{}
            pkgID string = ""
        )
        err := json.Unmarshal([]byte(payload), &struct {
            PkgID *string `json:"pkg-id"`
        }{&pkgID})
        if err != nil {
            feedError(errors.WithStack(err).Error())
            return err
        }

        pkgs, _ := model.FindPackage("pkg_id = ?", pkgID)
        if len(pkgs) == 0 {
            errMsg := fmt.Sprintf("selected package %s is not available", pkgID)
            feedError(errMsg)
            return errors.Errorf(errMsg)
        }
        // pick up the right package
        pkg := pkgs[0]

        // --- --- --- --- --- download meta first --- --- --- --- ---
        metaReq, err := newRequest(fmt.Sprintf("https://api.pocketcluster.io%s", pkg.MetaURL), false)
        if err != nil {
            return feedError("Unable to access package meta data. Reason : " + errors.WithStack(err).Error())
        }
        _, err = readRequest(metaReq, client)
        if err != nil {
            return feedError("Unable to access package meta data. Reason : " + errors.WithStack(err).Error())
        }

        //  --- --- --- --- --- download repo list --- --- --- --- ---
        repoReq, err := newRequest("https://api.pocketcluster.io/service/v014/package/repo", false)
        if err != nil {
            return feedError("Unable to access repository list. Reason : " + errors.WithStack(err).Error())
        }
        repoData, err := readRequest(repoReq, client)
        if err != nil {
            return feedError("Unable to access repository list. Reason : " + errors.WithStack(err).Error())
        }
        err = json.Unmarshal(repoData, &repoList)
        if err != nil {
            return feedError("Unable to access repository list. Reason : " + errors.WithStack(err).Error())
        }
        if len(repoList) == 0 {
            return feedError("Unable to access repository list. Reason : image repository unavailable")
        }

        //  --- --- --- --- --- download core sync --- --- --- --- ---
        cSyncReq, err := newRequest(fmt.Sprintf("https://api.pocketcluster.io%s", pkg.CoreImageSync), true)
        if err != nil {
            return feedError("Unable to access core image sync data. Reason : " + errors.WithStack(err).Error())
        }
        cSyncData, err := readRequest(cSyncReq, client)
        if err != nil {
            return feedError("Unable to access core image sync data. Reason : " + errors.WithStack(err).Error())
        }
        filesize, blocksize, blockcount, rootHash, err := readHeadersAndCheck(bytes.NewBuffer(cSyncData))
        if err != nil {
            return feedError("Unable to access core image sync data. Reason : " + errors.WithStack(err).Error())
        }
        log.Debugf("FileSize %v | BlockSize %v | BlockCount %v | RootHash %v ", filesize, blocksize, blockcount, rootHash)
        var coreImgRepoList []string = []string{}
        for _, r := range repoList {
            coreImgRepoList = append(coreImgRepoList, fmt.Sprintf("https://%s%s", r, pkg.CoreImageURL))
        }

        //  --- --- --- --- --- download node sync --- --- --- --- ---
        nSyncReq, err := newRequest(fmt.Sprintf("https://api.pocketcluster.io%s", pkg.NodeImageSync), true)
        if err != nil {
            return feedError("Unable to access node image sync data. Reason : " + errors.WithStack(err).Error())
        }
        nSyncData, err := readRequest(nSyncReq, client)
        if err != nil {
            return feedError("Unable to access node image sync data. Reason : " + errors.WithStack(err).Error())
        }
        filesize, blocksize, blockcount, rootHash, err = readHeadersAndCheck(bytes.NewBuffer(nSyncData))
        if err != nil {
            return feedError("Unable to access node image sync data. Reason : " + errors.WithStack(err).Error())
        }
        log.Debugf("FileSize %v | BlockSize %v | BlockCount %v | RootHash %v ", filesize, blocksize, blockcount, rootHash)
        var nodeImgRepoList []string = []string{}
        for _, r := range repoList {
            nodeImgRepoList = append(nodeImgRepoList, fmt.Sprintf("https://%s%s", r, pkg.NodeImageURL))
        }


        feedError("This is test error")
        return nil
    })
}