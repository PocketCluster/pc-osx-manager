package install

/*
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

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/service"
)
*/

import (
    "archive/tar"
    "io"
    "os"
    "path/filepath"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "xi2.org/x/xz"
)

const (
    timeout = time.Duration(10 * time.Second)

    irvicePackageSyncPatch   string = "irvice.package.sync.patch"
    irvicePackageSyncMonitor string = "irvice.package.sync.monitor"
    irvicePackageSyncUnarch  string = "irvice.package.sync.unarch"
    iventPackageSyncStop     string = "ivent.package.sync.stop"
    iventPackageSyncError    string = "ivent.package.sync.error"
    iventSyncActionPackage   string = "ivent.sync.action.package"
)

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

/*
func InitInstallPackageRoutePath(appLife route.Router, feeder route.ResponseFeeder) {
    // install a package
    appLife.POST(routepath.RpathPackageInstall(), func(_, rpath, payload string) error {
        var (
            feedError = func(irr error) error {
                log.Error(irr.Error())
                data, frr := json.Marshal(route.ReponseMessage{
                    "package-install": {
                        "status": false,
                        "error" : irr.Error(),
                    },
                })
                // this should never happen
                if frr != nil {
                    log.Error(frr.Error())
                }
                frr = feeder.FeedResponseForPost(rpath, string(data))
                if frr != nil {
                    log.Error(frr.Error())
                }
                return irr
            }

            rpathPkgInstProgress string = routepath.RpathPackageInstallProgress()
            pkgID string = ""
            pkg model.Package

            unarchC  = make(chan service.Event)
            monitorC = make(chan service.Event)
            errorC   = make(chan service.Event)
            stopC    = make(chan service.Event)
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

        // unarch func
        appLife.RegisterServiceWithFuncs(
            irvicePackageSyncUnarch,
            func() error {
                ae := <-unarchC
                action, ok := ae.Payload.(*patchActionPack)
                if !ok {
                    return errors.Errorf("[ERR] invalid patchActionPack type")
                }

                rDir, err := context.SharedHostContext().ApplicationRepositoryDirectory()
                if err != nil {
                    return err
                }
                err = xzUncompressor(action.reader, rDir)

                return nil
            },
            service.BindEventWithService(iventSyncActionPackage, unarchC))

        // monitor func
        appLife.RegisterServiceWithFuncs(
            irvicePackageSyncMonitor,
            func() error {

                ae := <-monitorC
                action, ok := ae.Payload.(*patchActionPack)
                if !ok {
                    return errors.Errorf("[ERR] invalid patchActionPack type")
                }

                for {
                    select {
                        case <- appLife.StopChannel(): {
                            closeAction(action)
                            feedError(errors.Errorf("core image sync halt"))
                            return nil
                        }
                        case <- stopC: {
                            closeAction(action)
                            feedError(errors.Errorf("core image sync halt"))
                            return nil
                        }
                        case e := <- errorC: {
                            if e.Payload != nil {
                                closeAction(action)
                            }
                            return nil
                        }
                        case rpt := <- action.report: {
                            data, err := json.Marshal(ReponseMessage{
                                "package-progress": {
                                    "total-size":   rpt.TotalSize,
                                    "received":     rpt.Received,
                                    "remaining":    rpt.Remaining,
                                    "speed":        rpt.Speed,
                                    "done-percent": rpt.DonePercent,
                                },
                            })
                            if err != nil {
                                log.Errorf(err.Error())
                                continue
                            }
                            err = feeder.FeedResponseForPost(rpathPkgInstProgress, string(data))
                            if err != nil {
                                log.Errorf(err.Error())
                                continue
                            }
                        }
                    }
                }

                return nil
            },
            service.BindEventWithService(iventSyncActionPackage, monitorC),
            service.BindEventWithService(iventPackageSyncError,  errorC),
            service.BindEventWithService(iventPackageSyncStop,   stopC))

        // register service to run
        appLife.RegisterServiceWithFuncs(
            irvicePackageSyncPatch,
            func() error {
                var (
                    client     = newClient(timeout, false)
                    repoList   = []string{}
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
                cActionPack, err := prepSync(repoList, cSyncData, pkg.CoreImageChksum, pkg.CoreImageURL)
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync core image"))
                }
                appLife.BroadcastEvent(service.Event{
                    Name:iventSyncActionPackage,
                    Payload:cActionPack})

                cAerr := cActionPack.msync.Patch()
                appLife.BroadcastEvent(service.Event{Name:iventPackageSyncError, Payload:cAerr})
                if cAerr != nil {
                    return feedError(errors.WithMessage(cAerr, "unable to sync core image"))
                }
                // FIXME: what if this is closed by user? no error means we cannot classify!
                closeAction(cActionPack)

                //  --- --- --- --- --- download node sync --- --- --- --- ---
                nSyncReq, err := newRequest(fmt.Sprintf("https://api.pocketcluster.io%s", pkg.NodeImageSync), true)
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync node image"))
                }
                nSyncData, err := readRequest(nSyncReq, client)
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync node image"))
                }
                nActionPack, err := prepSync(repoList, nSyncData, pkg.NodeImageChksum, pkg.NodeImageURL)
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync node image"))
                }
                err = execSync(nActionPack, stopC, rpathPkgInstProgress)
                if err != nil {
                    return feedError(errors.WithMessage(err, "unable to sync core image"))
                }

                return nil
            })

        return nil
    })
}

*/
