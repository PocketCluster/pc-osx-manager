package main

import (
    "archive/tar"
    "bytes"
    "encoding/base64"
    "encoding/binary"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "path/filepath"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/Redundancy/go-sync"
    "github.com/Redundancy/go-sync/blockrepository"
    "github.com/Redundancy/go-sync/blocksources"
    "github.com/Redundancy/go-sync/chunks"
    "github.com/Redundancy/go-sync/filechecksum"
    "github.com/Redundancy/go-sync/index"
    "github.com/Redundancy/go-sync/patcher"
    "github.com/Redundancy/go-sync/patcher/multisources"
    "github.com/Redundancy/go-sync/showpipe"
    "xi2.org/x/xz"

    //"github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/route/routepath"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/service"
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

func checkMetaChksum(data []byte, refData string) error {
    if len(data) == 0 {
        return errors.Errorf("invalid data to check")
    }
    if len(refData) == 0 {
        return errors.Errorf("invalid length of reference checksum")
    }
    refChksum, err := base64.URLEncoding.DecodeString(refData)
    if err != nil {
        return errors.WithStack(err)
    }
    hasher := filechecksum.DefaultStrongHashGenerator()
    hasher.Write(data)
    if bytes.Compare(hasher.Sum(nil), refChksum) != 0 {
        return errors.Errorf("invalid checksum value")
    }
    return nil
}

func isTwoChksumSame(chksum []byte, refData string) error {
    if len(chksum) == 0 {
        return errors.Errorf("invalid length of checksum to compare")
    }
    if len(refData) == 0 {
        return errors.Errorf("invalid length of reference checksum to compare")
    }
    refChksum, err := base64.URLEncoding.DecodeString(refData)
    if err != nil {
        return errors.WithStack(err)
    }
    if bytes.Compare(chksum, refChksum) != 0 {
        return errors.Errorf("invalid checksum value")
    }
    return nil
}

func readIndex(rd io.Reader, blocksize, blockcount uint, rootHash []byte) (*index.ChecksumIndex, error) {
    var (
        generator    = filechecksum.NewFileChecksumGenerator(blocksize)
        idx          *index.ChecksumIndex = nil
    )

    readChunks, err := chunks.CountedLoadChecksumsFromReader(
        rd,
        blockcount,
        generator.GetWeakRollingHash().Size(),
        generator.GetStrongHash().Size(),
    )
    if err != nil {
        return nil, errors.WithStack(err)
    }

    idx = index.MakeChecksumIndex(readChunks)
    cRootHash, err := idx.SequentialChecksumList().RootHash()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    if bytes.Compare(cRootHash, rootHash) != 0 {
        return nil, errors.Errorf("[ERR] mismatching integrity checksum")
    }

    return idx, nil
}

func xzUncompressor(archiveReader io.Reader, uncompPath string) error {
    var (
        xreader   *xz.Reader
        unarchive *tar.Reader
        err       error
    )

    // Check that the server actually sent compressed data
    xreader, err = xz.NewReader(archiveReader, 0)
    if err != nil {
        return errors.WithStack(err)
    }

    unarchive = tar.NewReader(xreader)
    for {
        header, err := unarchive.Next()
        if err == io.EOF {
            break
        } else if err != nil {
            return errors.WithStack(err)
        }

        path := filepath.Join(uncompPath, header.Name)
        info := header.FileInfo()
        if info.IsDir() {
            if err = os.MkdirAll(path, info.Mode()); err != nil {
                return errors.WithStack(err)
            }
            continue
        }
        file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
        if err != nil {
            return errors.WithStack(err)
        }
        written, err := io.Copy(file, unarchive)
        if err != nil {
            file.Close()
            return errors.WithStack(err)
        } else {
            log.Debugf("written %v", written)
        }
        err = file.Close()
        if err != nil {
            return errors.WithStack(err)
        }
    }
    return nil
}

func initInstallRoutePath() {
    const (
        timeout = time.Duration(10 * time.Second)
    )

    // get the list of available packages
    theApp.GET(routepath.RpathPackageList(), func(_, rpath, _ string) error {
        var (
            feedError = func(irr error) error {
                data, frr := json.Marshal(ReponseMessage{
                    "package-list": {
                        "status": false,
                        "error" : irr.Error(),
                    },
                })
                if frr != nil {
                    log.Debugf(frr.Error())
                }
                frr = FeedResponseForGet(rpath, string(data))
                if frr != nil {
                    log.Debugf(frr.Error())
                }
                return errors.WithStack(irr)
            }

            pkgList = []map[string]interface{}{}
            pkgs    = []*model.Package{}
        )

        req, err :=  newRequest("https://api.pocketcluster.io/service/v014/package/list", false)
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to access package list"))
        }
        client := newClient(timeout, true)
        resp, err := readRequest(req, client)
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to access package list"))
        }
        err = json.Unmarshal(resp, &pkgs)
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to access package list"))
        }
        if len(pkgs) == 0 {
            return feedError(errors.Errorf("No package avaiable. Contact us at Slack channel."))
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
            return feedError(errors.WithMessage(err, "Unable to access package list"))
        }
        err = FeedResponseForGet(rpath, string(data))
        if err != nil {
            return errors.WithStack(err)
        }
        return nil
    })

    // install a package
    theApp.POST(routepath.RpathPackageInstall(), func(_, rpath, payload string) error {
        const (
            irvicePackageSyncPatch      string = "irvice.package.sync.patch"
            irvicePackageSyncControl    string = "irvice.package.sync.control"
            iventPackageSyncReportCore  string = "ivent.package.sync.report.core"
            iventPackageSyncReportNode  string = "ivent.package.sync.report.node"
        )
        type rptPack struct {
            reader    io.Reader
            msrc      *multisources.MultiSourcePatcher
            report    chan showpipe.PipeProgress
        }
        var (
            feedError = func(irr error) error {
                log.Error(irr.Error())
                data, frr := json.Marshal(ReponseMessage{
                    "package-install": {
                        "status": false,
                        "error" : irr.Error(),
                    },
                })
                // this should never happen
                if frr != nil {
                    log.Error(frr.Error())
                }
                frr = FeedResponseForPost(rpath, string(data))
                if frr != nil {
                    log.Error(frr.Error())
                }
                return irr
            }

            pkgID string = ""
            pkg model.Package
        )

        // 1. parse input package id
        err := json.Unmarshal([]byte(payload), &struct {
            PkgID *string `json:"pkg-id"`
        }{&pkgID})
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to specify package package"))
        }

        // TODO 2. check if service is already running

        // 3. find appropriate model
        pkgs, _ := model.FindPackage("pkg_id = ?", pkgID)
        if len(pkgs) == 0 {
            return feedError(errors.Errorf("selected package %s is not available", pkgID))
        }

        // 4. pick up the first package & we are ready to patch.
        pkg = pkgs[0]

        // register service to run
        theApp.RegisterServiceWithFuncs(
            irvicePackageSyncPatch,
            func() error {
                var (
                    client = newClient(timeout, false)
                    repoList = []string{}
                )

                // --- --- --- --- --- download meta first --- --- --- --- ---
                metaReq, err := newRequest(fmt.Sprintf("https://api.pocketcluster.io%s", pkg.MetaURL), false)
                if err != nil {
                    return feedError(errors.WithMessage(err, "Unable to access package meta data"))
                }
                metaData, err := readRequest(metaReq, client)
                if err != nil {
                    return feedError(errors.WithMessage(err, "Unable to access package meta data"))
                }
                err = checkMetaChksum(metaData, pkg.MetaChksum)
                if err != nil {
                    return feedError(errors.WithMessage(err, "Unable to access package meta data"))
                }
                // TODO : save meta

                //  --- --- --- --- --- download repo list --- --- --- --- ---
                repoReq, err := newRequest("https://api.pocketcluster.io/service/v014/package/repo", false)
                if err != nil {
                    return feedError(errors.WithMessage(err, "Unable to access repository list"))
                }
                repoData, err := readRequest(repoReq, client)
                if err != nil {
                    return feedError(errors.WithMessage(err, "Unable to access repository list"))
                }
                err = json.Unmarshal(repoData, &repoList)
                if err != nil {
                    return feedError(errors.WithMessage(err, "Unable to access repository list"))
                }
                if len(repoList) == 0 {
                    return feedError(errors.WithMessage(err, "Unable to access repository list"))
                }

                //  --- --- --- --- --- download core sync --- --- --- --- ---
                cSyncReq, err := newRequest(fmt.Sprintf("https://api.pocketcluster.io%s", pkg.CoreImageSync), true)
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync core image"))
                }
                cSyncData, err := readRequest(cSyncReq, client)
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync core image"))
                }
                cFilesize, cBlocksize, cBlockcount, cRootHash, err := readHeadersAndCheck(bytes.NewBuffer(cSyncData))
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync core image"))
                }
                err = isTwoChksumSame(cRootHash, pkg.CoreImageChksum)
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync core image"))
                }
                cChksumIdx, err := readIndex(bytes.NewBuffer(cSyncData), uint(cBlocksize), uint(cBlockcount), cRootHash)
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync core image"))
                }

                var (
                    cReader, cWriter, cRprtr = showpipe.PipeWithReport(uint64(cFilesize))
                    cResolver                = blockrepository.MakeKnownFileSizedBlockResolver(int64(cBlocksize), cFilesize)
                    cVerifier                = &filechecksum.HashVerifier{
                        Hash:                filechecksum.DefaultStrongHashGenerator(),
                        BlockSize:           uint(cBlocksize),
                        BlockChecksumGetter: cChksumIdx,
                    }
                    cRepoList []patcher.BlockRepository = nil
                )
                for rID, r := range repoList {
                    cRepoList = append(cRepoList,
                        blockrepository.NewBlockRepositoryBase(
                            uint(rID),
                            blocksources.NewRequesterWithTimeout(fmt.Sprintf("https://%s%s", r, pkg.CoreImageURL), "PocketCluster/0.1.4 (OSX)", true, timeout),
                            cResolver,
                            cVerifier))
                }
                cSync, err := multisources.NewMultiSourcePatcher(cWriter, cRepoList, cChksumIdx)
                theApp.BroadcastEvent(service.Event{
                    Name:iventPackageSyncReportCore,
                    Payload: rptPack{
                            reader: cReader,
                            msrc:   cSync,
                            report: cRprtr,
                    }})

                err = cSync.Patch()
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync core image"))
                }


                //  --- --- --- --- --- download node sync --- --- --- --- ---
                nSyncReq, err := newRequest(fmt.Sprintf("https://api.pocketcluster.io%s", pkg.NodeImageSync), true)
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync node image"))
                }
                nSyncData, err := readRequest(nSyncReq, client)
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync node image"))
                }
                nFilesize, nBlocksize, nBlockcount, nRootHash, err := readHeadersAndCheck(bytes.NewBuffer(nSyncData))
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync node image"))
                }
                nChksumIdx, err := readIndex(bytes.NewBuffer(nSyncData), uint(nBlocksize), uint(nBlockcount), nRootHash)
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync node image"))
                }

                var (
                    nReader, nWriter, nRprtr = showpipe.PipeWithReport(uint64(nFilesize))
                    nResolver                = blockrepository.MakeKnownFileSizedBlockResolver(int64(nBlocksize), nFilesize)
                    nVerifier                = &filechecksum.HashVerifier{
                        Hash:                filechecksum.DefaultStrongHashGenerator(),
                        BlockSize:           uint(nBlocksize),
                        BlockChecksumGetter: nChksumIdx,
                    }
                    nRepoList []patcher.BlockRepository = nil
                )
                for rID, r := range repoList {
                    nRepoList = append(nRepoList,
                        blockrepository.NewBlockRepositoryBase(
                            uint(rID),
                            blocksources.NewRequesterWithTimeout(fmt.Sprintf("https://%s%s", r, pkg.NodeImageURL), "PocketCluster/0.1.4 (OSX)", true, timeout),
                            nResolver,
                            nVerifier))
                }
                nSync, err := multisources.NewMultiSourcePatcher(nWriter, nRepoList, nChksumIdx)
                theApp.BroadcastEvent(service.Event{
                    Name:iventPackageSyncReportNode,
                    Payload: rptPack{
                        reader: nReader,
                        msrc:   nSync,
                        report: nRprtr,
                    }})

                err = nSync.Patch()
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync node image"))
                }

                return nil
            })

        return nil
    })
}