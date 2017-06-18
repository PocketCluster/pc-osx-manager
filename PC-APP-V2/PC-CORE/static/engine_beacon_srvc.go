package main

import (
    "time"

    "github.com/gravitational/teleport/embed"
    tervice "github.com/gravitational/teleport/lib/service"
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/pc-core/beacon"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/model"
    swarmemb "github.com/stkim1/pc-core/extlib/swarm"
)

type beaconEventRoute struct {
    service.ServiceSupervisor
    *tervice.PocketConfig
}

func (b *beaconEventRoute) close() error {
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

func initMasterBeaconService(a *mainLife, clusterID string, tcfg *tervice.PocketConfig) error {
    var (
        beaconC = make(chan service.Event)
        searchC = make(chan service.Event)
    )
    a.RegisterServiceWithFuncs(
        operation.ServiceBeaconMaster,
        func() error {
            var (
                timer = time.NewTicker(time.Second)

                beaconRoute *beaconEventRoute = &beaconEventRoute{
                    ServiceSupervisor:a.ServiceSupervisor,
                    PocketConfig: tcfg,
                }

                beaconMan, err = beacon.NewBeaconManagerWithFunc(
                    clusterID,
                    beaconRoute,
                    func(host string, payload []byte) error {
                        log.Debugf("[BEACON-SEND-SLAVE] [%v] Host %v", time.Now(), host)
                        a.BroadcastEvent(service.Event{
                            Name: ucast.EventBeaconCoreWriteLocation,
                            Payload:ucast.BeaconSend{
                                Host:       host,
                                Payload:    payload,
                            },
                        })
                        return nil
                    })
            )
            if err != nil {
                return errors.WithStack(err)
            }

            log.Debugf("[AGENT] starting agent service...")
            a.BroadcastEvent(service.Event{Name:iventBeaconManagerSpawn, Payload:beaconMan})
            for {
                select {
                    case <-a.StopChannel(): {
                        timer.Stop()
                        err = beaconMan.Shutdown()
                        if err != nil {
                            log.Debug(err.Error())
                        }
                        err = beaconRoute.close()
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
                                log.Debugf("[BEACON-TRANSITION] %v", err)
                            }
                        }
                    }
                    case s := <-searchC: {
                        cp, ok := s.Payload.(mcast.CastPack)
                        if ok {
                            err = beaconMan.TransitionWithSearchData(cp, time.Now())
                            if err != nil {
                                log.Debugf("[SEARCH-TRANSITION] %v", err)
                            }
                        }
                    }
                    case <-timer.C: {
                        err = beaconMan.TransitionWithTimestamp(time.Now())
                        if err != nil {
                            log.Debug(err.Error())
                        }
                    }
                }
            }
            return nil
        },
        service.BindEventWithService(ucast.EventBeaconCoreReadLocation, beaconC),
        service.BindEventWithService(mcast.EventBeaconCoreReadSearch, searchC))

    return nil
}

func initSwarmService(a *mainLife) error {
    swarmSrvC := make(chan service.Event)
    a.RegisterServiceWithFuncs(
        operation.ServiceSwarmEmbeddedOperation,
        func() error {
            var (
                swarmsrv *swarmemb.SwarmService = nil
            )
            select {
                case se := <- swarmSrvC: {
                    srv, ok := se.Payload.(*swarmemb.SwarmService)
                    if ok {
                        swarmsrv = srv
                    }
                }
                case <- a.StopChannel(): {
                    if swarmsrv != nil {
                        err := swarmsrv.Close()
                        return errors.WithStack(err)
                    }
                    return errors.Errorf("[ERR] null SWARM instance")
                }
            }
            return nil
        },
        service.BindEventWithService(iventSwarmInstanceSpawn, swarmSrvC))

    beaconManC := make(chan service.Event)
    a.RegisterServiceWithFuncs(
        operation.ServiceSwarmEmbeddedServer,
        func() error {
            be := <- beaconManC
            beaconMan, ok := be.Payload.(beacon.BeaconManger)
            if !ok {
                return errors.Errorf("[ERR] invalid beacon manager type")
            }
            ctx := context.SharedHostContext()
            caCert, err := ctx.CertAuthCertificate()
            if err != nil {
                return errors.WithStack(err)
            }
            hostCrt, err := ctx.MasterHostCertificate()
            if err != nil {
                return errors.WithStack(err)
            }
            hostPrv, err := ctx.MasterHostPrivateKey()
            if err != nil {
                return errors.WithStack(err)
            }
            swarmctx, err := swarmemb.NewContextWithCertAndKey(caCert, hostCrt, hostPrv, beaconMan)
            if err != nil {
                return errors.WithStack(err)
            }
            swarmsrv, err := swarmemb.NewSwarmService(swarmctx)
            if err != nil {
                return errors.WithStack(err)
            }
            a.BroadcastEvent(service.Event{Name:iventSwarmInstanceSpawn, Payload:swarmsrv})
            err = swarmsrv.ListenAndServeSingleHost()
            return errors.WithStack(err)
        },
        service.BindEventWithService(iventBeaconManagerSpawn, beaconManC))

    return nil
}

