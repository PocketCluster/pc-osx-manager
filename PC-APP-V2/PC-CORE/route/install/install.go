package install

import (
    "context"
    "encoding/json"
    "fmt"
    "os/user"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    tervice "github.com/gravitational/teleport/lib/service"
    tclient "github.com/gravitational/teleport/lib/client"

    pcctx "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/utils/dockertool"
)

const (
    timeout = time.Duration(10 * time.Second)
)

func InitInstallPackageRoutePath(appLife route.Router, feeder route.ResponseFeeder, sshCfg *tervice.PocketConfig) {
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

            pkg        *model.Package
            uRoot      *model.UserMeta
            rpProgress string = routepath.RpathPackageInstallProgress()
            pkgID      string = ""
            repoList          = []string{}
            stopC             = make(chan struct{})
            client            = newClient(timeout, false)
        )

        // 0. find registry destination first
        regDir, err := pcctx.SharedHostContext().ApplicationRepositoryDirectory()
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

        // 3. pick up the first package & we are ready to patch.
        pkgs, err := model.FindPackage("pkg_id = ?", pkgID)
        if err != nil {
            return feedError(errors.Errorf("selected package %s is not available", pkgID))
        }
        pkg = pkgs[0]

        // 4. pick root password for devops
        rusers, err := model.FindUserMetaWithLogin("root")
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to install package due to improper permission"))
        }
        uRoot = rusers[0]

        // 5. read local user information for devops
        luname, err := pcctx.SharedHostContext().LoginUserName()
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to install package due to invalid user information"))
        }
        luser, err := user.Lookup(luname)
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to access user information"))
        }
        log.Infof("user name %v, user id %v", luname, luser.Uid)


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
        tmpl, err := model.FindTemplateWithPackageID(pkgID)
        if err != nil {
            if err == model.NoItemFound {
                tmpl = model.NewTemplateMeta()
            } else {
                return feedError(errors.WithMessage(err, "Unable to store package composer meta data"))
            }
        }
        tmpl.PkgID = pkgID
        tmpl.Body  = metaData
        err = tmpl.Update()
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to store package composer meta data"))
        }


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

        // --- --- --- --- --- install image to core --- --- --- --- ---
        _ = makeMessageFeedBack(feeder, rpProgress, "Installing core image...")

        // --- --- --- --- --- install image to nodes --- --- --- --- ---
        // TODO : we can request swarm server to do this job
        _ = makeMessageFeedBack(feeder, rpProgress, "Installing node image...")
        cli, err := dockertool.NewContainerClient("tcp://pc-node1:2376", "1.24")
        if err != nil {
            return feedError(errors.WithMessage(err, "unable to make connection to " + "pc-node1"))
        }
        err = dockertool.InstallImageFromRepository(cli, "pc-master:5000/arm64v8-ubuntu")
        if err != nil {
            return feedError(errors.WithMessage(err, "unable to sync image to " + "pc-node1"))
        }

        // --- --- --- --- --- setup node --- --- --- --- ---
        tNode := "pc-node1"
        c, err := tclient.MakeNewClient(sshCfg, uRoot.Login, tNode)
        if err != nil {
            return feedError(errors.WithMessage(err, "unable to setup package to " + tNode))
        }
        err = c.APISSH(context.TODO(), []string{"mkdir -p /pocket/basic/package"}, "1524rmfo",false)
        if err != nil {
            return feedError(errors.WithMessage(err, "unable to setup package to " + tNode))
        }
        err = c.Logout()
        if err != nil {
            return feedError(errors.WithMessage(err, "unable to setup package to " + tNode))
        }


        // --- --- --- --- --- install image to nodes --- --- --- --- ---
        data, err := json.Marshal(route.ReponseMessage{
            "package-install": {
                "status": true,
                "pkg-id" : pkgID,
            },
        })
        // this should never happen
        if err != nil {
            log.Error(err.Error())
        }
        err = feeder.FeedResponseForPost(rpath, string(data))
        return errors.WithStack(err)
    })
}
