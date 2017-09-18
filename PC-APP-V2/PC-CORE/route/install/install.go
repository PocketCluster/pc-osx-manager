package install

import (
    "encoding/json"
    "fmt"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/model"
)

const (
    timeout = time.Duration(10 * time.Second)
)

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

            pkg        model.Package
            rpProgress string = routepath.RpathPackageInstallProgress()
            pkgID      string = ""
            repoList          = []string{}
            stopC             = make(chan struct{})
            client            = newClient(timeout, false)
        )

        // 0. find registry destination first
        regDir, err := context.SharedHostContext().ApplicationRepositoryDirectory()
        if err != nil {
            return errors.WithStack(err)
        }

        // 1. parse input package id
        err = json.Unmarshal([]byte(payload), &struct {
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


        // --- --- --- --- --- download meta first --- --- --- --- ---
        _ = makeMessageFeedBack(feeder, rpProgress, "Downloading package information...")
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
        _ = makeMessageFeedBack(feeder, rpProgress, "Checking image repositories...")
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
        _ = makeMessageFeedBack(feeder, rpProgress, "Downloading core image...")
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
        err = execSync(feeder, cActionPack, stopC, rpProgress, regDir)
        if err != nil {
            return feedError(errors.WithMessage(err, "unable to sync core image"))
        }

        //  --- --- --- --- --- download node sync --- --- --- --- ---
        _ = makeMessageFeedBack(feeder, rpProgress, "Downloading node image...")
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
        err = execSync(feeder, nActionPack, stopC, rpProgress, regDir)
        if err != nil {
            return feedError(errors.WithMessage(err, "unable to sync node image"))
        }

        return feedError(errors.Errorf("ALL DONE OK"))
    })
}
