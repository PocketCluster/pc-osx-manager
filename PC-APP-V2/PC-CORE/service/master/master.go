package master

import (
    "time"

    "github.com/gravitational/teleport/embed"
    tervice "github.com/gravitational/teleport/lib/service"
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/pc-core/beacon"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/extlib/pcssh/sshproc"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/service/ivent"
)

type beaconEventRoute struct {
    service.ServiceSupervisor
    *tervice.PocketConfig
}

func (b *beaconEventRoute) terminate() error {
    b.ServiceSupervisor = nil
    b.PocketConfig = nil
    return nil
}

func (b *beaconEventRoute) BeaconEventPrepareJoin(slave *model.SlaveNode) error {
    clt, err := embed.OpenAdminClientWithAuthService(b.PocketConfig)
    if err != nil {
        return errors.WithStack(err)
    }
    defer clt.Close()
    token, err := embed.GenerateNodeInviationWithTTL(clt, embed.MaxInvitationTLL)
    if err != nil {
        return errors.WithStack(err)
    }
    return errors.WithStack(slave.SetAuthToken(token))
}

func (b *beaconEventRoute) BeaconEventResurrect(slaves []model.SlaveNode) error {
    return nil
}

func (b *beaconEventRoute) BeaconEventTranstion(state beacon.MasterBeaconState, slave *model.SlaveNode, ts time.Time, transOk bool) error {
    if transOk {
        log.Debugf("(INFO) [%v | %v] BeaconEventTranstion -> %v | SUCCESS ", ts, slave.AuthToken, state.String())
    } else {
        log.Debugf("(INFO) [%v | %v] BeaconEventTranstion -> %v | FAILED ", ts, slave.AuthToken, state.String())
    }

    return nil
}

func (b *beaconEventRoute) BeaconEventDiscard(slave *model.SlaveNode) error {
    return nil
}

func (b *beaconEventRoute) BeaconEventShutdown() error {
    return nil
}

func InitMasterBeaconService(appLife service.ServiceSupervisor, clusterID string, tcfg *tervice.PocketConfig) error {
    var (
        beaconC = make(chan service.Event)
        searchC = make(chan service.Event)
        vboxC   = make(chan service.Event)
        teleC   = make(chan service.Event)
        netC    = make(chan service.Event)
    )
    appLife.RegisterServiceWithFuncs(
        operation.ServiceBeaconMaster,
        func() error {
            var (
                beaconRoute *beaconEventRoute  = &beaconEventRoute{
                    ServiceSupervisor: appLife,
                    PocketConfig:      tcfg,
                }
                timer                          = time.NewTicker(time.Second)
                beaconMan  beacon.BeaconManger = nil
                err        error               = nil
            )
            // wait for vbox control
            vc := <- vboxC
            vbc, ok := vc.Payload.(vboxCtrlObjBrcst)
            if !ok {
                log.Debugf("[AGENT] (ERR) invalid VBoxMasterControl type")
                return errors.Errorf("[ERR] invalid VBoxMasterControl type")
            }

            // beacon manager
            beaconMan, err = beacon.NewBeaconManagerWithFunc(
                clusterID,
                vbc.VBoxMasterControl,
                beaconRoute,
                func(host string, payload []byte) error {
                    log.Debugf("[AGENT] BEACON-TX [%v] Host %v", time.Now(), host)
                    appLife.BroadcastEvent(
                        service.Event{
                            Name:       ucast.EventBeaconCoreLocationSend,
                            Payload:    ucast.BeaconSend{
                                Host:       host,
                                Payload:    payload,
                            },
                        })
                    return nil
                })
            if err != nil {
                return errors.WithStack(err)
            }

            // waiting teleport to start
            <- teleC

            log.Debugf("[AGENT] starting agent service...")
            appLife.BroadcastEvent(service.Event{Name:ivent.IventBeaconManagerSpawn, Payload:beaconMan})
            for {
                select {
                    case <-appLife.StopChannel(): {
                        timer.Stop()
                        err = beaconMan.Shutdown()
                        if err != nil {
                            log.Debug(err.Error())
                        }
                        err = beaconRoute.terminate()
                        if err != nil {
                            log.Debug(err.Error())
                        }
                        log.Debugf("[AGENT] stopping agent service...")
                        return nil
                    }
                    case b := <-beaconC: {
                        bp, ok := b.Payload.(ucast.BeaconPack)
                        if ok {
                            err = beaconMan.TransitionWithBeaconData(bp, time.Now())
                            if err != nil {
                                log.Debugf("[AGENT] BEACON-RX Error : %v", err)
                            }
                        }
                    }
                    case s := <-searchC: {
                        cp, ok := s.Payload.(mcast.CastPack)
                        if ok {
                            err = beaconMan.TransitionWithSearchData(cp, time.Now())
                            if err != nil {
                                log.Debugf("[AGENT] SEARCH-RX Error : %v", err)
                            }
                        }
                    }
                    case <-timer.C: {
                        err = beaconMan.TransitionWithTimestamp(time.Now())
                        if err != nil {
                            log.Debug(err.Error())
                        }
                        regNodes := beaconMan.RegisteredNodesList()
                        appLife.BroadcastEvent(service.Event{Name:ivent.IventMonitorRegisteredNode, Payload:regNodes})
                    }
                    case <- netC: {
                        // TODO update primary address
                        log.Debugf("[AGENT] Host Address changed")
                    }
                }
            }
            return nil
        },
        service.BindEventWithService(ucast.EventBeaconCoreLocationReceive, beaconC),
        service.BindEventWithService(mcast.EventBeaconCoreSearchReceive,   searchC),
        service.BindEventWithService(ivent.IventVboxCtrlInstanceSpawn,           vboxC),
        service.BindEventWithService(sshproc.EventPCSSHServerProxyStarted, teleC),
        service.BindEventWithService(ivent.IventNetworkAddressChange,            netC))

    return nil
}
