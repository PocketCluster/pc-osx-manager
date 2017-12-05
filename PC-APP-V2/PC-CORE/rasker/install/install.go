package install

import (
    "context"
    "encoding/json"
    "fmt"
    "os/user"
    "strings"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    tervice "github.com/gravitational/teleport/lib/service"
    tclient "github.com/gravitational/teleport/lib/client"

    pcctx "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/defaults"
    "github.com/stkim1/pc-core/rasker"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/service/ivent"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/utils/dockertool"
    "github.com/stkim1/pc-core/utils/apireq"
)

func InitRoutePathInstallPackage(appLife rasker.RouteTasker, feeder route.ResponseFeeder, sshCfg *tervice.PocketConfig) {
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

            rpProgress = routepath.RpathPackageInstallProgress()
            stopC      = make(chan struct{})
            reportC    = make(chan service.Event)
            client     = apireq.NewClient(apireq.ConnTimeout, false)
            repoList   = []string{}
            pkg        *model.Package = nil
            uRoot      *model.UserMeta = nil
            pkgID      = ""
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

        // TODO 3. check if this package has been installed

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

        // 6. get the live node to install
        err = appLife.BindDiscreteEvent(ivent.IventReportLiveNodesResult, reportC)
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to install package due to invalid node status"))
        }
        // ask node list
        appLife.BroadcastEvent(service.Event{Name:ivent.IventReportLiveNodesRequest})
        nr := <- reportC
        appLife.UntieDiscreteEvent(ivent.IventReportLiveNodesResult)
        nlist, ok := nr.Payload.([]string)
        if !ok {
            return feedError(errors.WithMessage(err, "Unable to install package due to invalid node status"))
        }

        // --- --- --- --- --- download meta first --- --- --- --- ---
        _ = makeMessageFeedBack(feeder, rpProgress, "Downloading package information...")
        metaReq, err := apireq.NewRequest(fmt.Sprintf("%s%s", defaults.PocketClusterAPIHost, pkg.MetaURL), false)
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to access package meta data"))
        }
        metaData, err := apireq.ReadRequest(metaReq, client)
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
        repoReq, err := apireq.NewRequest(fmt.Sprintf("%s/service/v014/package/repo", defaults.PocketClusterAPIHost), false)
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to access repository list"))
        }
        repoData, err := apireq.ReadRequest(repoReq, client)
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
        cSyncReq, err := apireq.NewRequest(fmt.Sprintf("%s%s", defaults.PocketClusterAPIHost, pkg.CoreImageSync), true)
        if err != nil {
            return feedError(errors.WithMessage(err, "unable to sync core image"))
        }
        cSyncData, err := apireq.ReadRequest(cSyncReq, client)
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
        nSyncReq, err := apireq.NewRequest(fmt.Sprintf("%s%s", defaults.PocketClusterAPIHost, pkg.NodeImageSync), true)
        if err != nil {
            return feedError(errors.WithMessage(err, "unable to sync node image"))
        }
        nSyncData, err := apireq.ReadRequest(nSyncReq, client)
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

        for _, tNode := range nlist {
            switch tNode {
                case defaults.PocketClusterCoreName: {
                    // --- --- --- --- --- install image to core --- --- --- --- ---
                    _ = makeMessageFeedBack(feeder, rpProgress, "Installing core image...")
                    ccli, err := dockertool.NewContainerClient(fmt.Sprintf("tcp://%s:%s", defaults.PocketClusterCoreName, defaults.DefaultSecureDockerPort), "1.26")
                    if err != nil {
                        return feedError(errors.WithMessage(err, "unable to make connection to pc-core"))
                    }
                    err = dockertool.InstallImageFromRepository(ccli, pkg.CoreImageName)
                    if err != nil {
                        return feedError(errors.WithMessage(err, "unable to sync image to " + defaults.PocketClusterCoreName))
                    }

                    // --- --- --- --- --- setup core node --- --- --- --- ---
                    // data paths to build
                    cdPath := strings.Split(pkg.CoreDataPath, "|")
                    cdPathCmds := []string{}
                    for _, cdp := range cdPath {
                        cdPathCmds = append(cdPathCmds, fmt.Sprintf("mkdir -p %s", cdp))
                        /*
                         * as in, vboxmount option with the core user dictates how permission is given to a new dir.
                         * all the new dir addition by any1 will have a same ownership + perm.
                         */
                        //cdPathCmds = append(cdPathCmds, fmt.Sprintf("chown -R %s:%s %s", luname, luname, cdp))
                        //cdPathCmds = append(cdPathCmds, fmt.Sprintf("chmod -R 755 %s", cdp))
                    }
                    log.Info("core data path commands %v", cdPathCmds)

                    cssh, err := tclient.MakeNewClient(sshCfg, uRoot.Login, defaults.PocketClusterCoreName)
                    if err != nil {
                        return feedError(errors.WithMessage(err, "unable to setup package to " + defaults.PocketClusterCoreName))
                    }
                    for _, cdpc := range cdPathCmds {
                        err = cssh.APISSH(context.TODO(), []string{cdpc}, uRoot.Password,false)
                        if err != nil {
                            log.Error(cdpc)
                            return feedError(errors.WithMessage(err, "unable to setup package to " + defaults.PocketClusterCoreName))
                        }
                    }
                    err = cssh.Logout()
                    if err != nil {
                        return feedError(errors.WithMessage(err, "unable to setup package to " + defaults.PocketClusterCoreName))
                    }
                }

                default: {
                    // --- --- --- --- --- install image to nodes --- --- --- --- ---
                    _ = makeMessageFeedBack(feeder, rpProgress, fmt.Sprintf("Installing node image to %v", tNode))
                    ncli, err := dockertool.NewContainerClient(fmt.Sprintf("tcp://%s:%s", tNode, defaults.DefaultSecureDockerPort), "1.26")
                    if err != nil {
                        return feedError(errors.WithMessage(err, "unable to make connection to " + tNode))
                    }
                    err = dockertool.InstallImageFromRepository(ncli, pkg.NodeImageName)
                    if err != nil {
                        return feedError(errors.WithMessage(err, "unable to sync image to " + tNode))
                    }

                    // --- --- --- --- --- setup node --- --- --- --- ---
                    // data paths to build
                    ndPath := strings.Split(pkg.NodeDataPath, "|")
                    // ndpath setup commands
                    ndPathCmds := []string{}
                    for _, ndp := range ndPath {
                        ndPathCmds = append(ndPathCmds, fmt.Sprintf("mkdir -p %s", ndp))
                        ndPathCmds = append(ndPathCmds, fmt.Sprintf("chown -R %s:%s %s", luname, luname, ndp))
                        ndPathCmds = append(ndPathCmds, fmt.Sprintf("chmod -R 755 %s", ndp))
                    }
                    log.Infof("node data path commands %v", ndPathCmds)

                    nssh, err := tclient.MakeNewClient(sshCfg, uRoot.Login, tNode)
                    if err != nil {
                        return feedError(errors.WithMessage(err, "unable to setup package to " + tNode))
                    }
                    for _, ndpc := range ndPathCmds {
                        err = nssh.APISSH(context.TODO(), []string{ndpc}, uRoot.Password,false)
                        if err != nil {
                            log.Error(ndpc)
                            return feedError(errors.WithMessage(err, "unable to setup package to " + tNode))
                        }
                    }
                    err = nssh.Logout()
                    if err != nil {
                        return feedError(errors.WithMessage(err, "unable to setup package to " + tNode))
                    }
                }
            }
        }


        // --- --- --- --- --- make installed package record --- --- --- --- ---
        err = model.UpsertRecords([]*model.PkgRecord{
            {
                AppVer: pkg.AppVer,
                PkgID:  pkg.PkgID,
                PkgVer: pkg.PkgVer,
                PkgChksum: pkg.PkgChksum,
            },
        })
        if err != nil {
            return feedError(errors.WithMessage(err, "unable to record package history" + pkg.Name))
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
