package main

import (
    "github.com/pkg/errors"
    tervice "github.com/gravitational/teleport/lib/service"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/extlib/pcssh/sshadmin"
    "github.com/stkim1/pc-core/model"
    pervice "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/service/ivent"
    "github.com/stkim1/pc-core/vboxglue"
)

const (
    serviceSetupUsersAndCore string = "service.setup.users.and.core"
)

// setup base users
func setupBaseUsersWithVboxCore(appLife pervice.ServiceSupervisor, pcsshCfg *tervice.PocketConfig) error {
    var (
        teleC = make(chan pervice.Event)
    )
    appLife.RegisterServiceWithFuncs(
        serviceSetupUsersAndCore,
        func() error {

            // wait teleport to be activated
            <- teleC

            cli, err := sshadmin.OpenAdminClientWithAuthService(pcsshCfg)
            if err != nil {
                appLife.BroadcastEvent(pervice.Event{Name:ivent.IventSetupUsersAndVboxCore, Payload:err})
                return errors.WithStack(err)
            }
            roots, err := model.FindUserMetaWithLogin("root")
            if err != nil {
                appLife.BroadcastEvent(pervice.Event{Name:ivent.IventSetupUsersAndVboxCore, Payload:err})
                return errors.WithStack(err)
            }
            err = sshadmin.CreateTeleportUser(cli, "root", roots[0].Password)
            if err != nil {
                appLife.BroadcastEvent(pervice.Event{Name:ivent.IventSetupUsersAndVboxCore, Payload:err})
                return errors.WithStack(err)
            }
            uname, err := context.SharedHostContext().LoginUserName()
            if err != nil {
                appLife.BroadcastEvent(pervice.Event{Name:ivent.IventSetupUsersAndVboxCore, Payload:err})
                return errors.WithStack(err)
            }
            lusers, err := model.FindUserMetaWithLogin(uname)
            if err != nil {
                appLife.BroadcastEvent(pervice.Event{Name:ivent.IventSetupUsersAndVboxCore, Payload:err})
                return errors.WithStack(err)
            }
            err = sshadmin.CreateTeleportUser(cli, uname, lusers[0].Password)
            if err != nil {
                appLife.BroadcastEvent(pervice.Event{Name:ivent.IventSetupUsersAndVboxCore, Payload:err})
                return errors.WithStack(err)
            }

            // --- now setup vbox core ---
            cid, err := context.SharedHostContext().MasterAgentName()
            if err != nil {
                appLife.BroadcastEvent(pervice.Event{Name:ivent.IventSetupUsersAndVboxCore, Payload:err})
                return errors.WithStack(err)
            }
            err = vboxglue.BuildVboxCoreDisk(cid, pcsshCfg)
            if err != nil {
                appLife.BroadcastEvent(pervice.Event{Name:ivent.IventSetupUsersAndVboxCore, Payload:err})
                return errors.WithStack(err)
            }
            vcore, err := vboxglue.NewGOVboxGlue()
            if err != nil {
                appLife.BroadcastEvent(pervice.Event{Name:ivent.IventSetupUsersAndVboxCore, Payload:err})
                return errors.WithStack(err)
            }
            err = vboxglue.CreateNewMachine(vcore)
            vcore.Close()
            appLife.BroadcastEvent(pervice.Event{Name:ivent.IventSetupUsersAndVboxCore, Payload:err})
            return errors.WithStack(err)
        },
        pervice.BindEventWithService(ivent.IventPcsshProxyInstanceSpawn,   teleC),
    )
    return nil
}
