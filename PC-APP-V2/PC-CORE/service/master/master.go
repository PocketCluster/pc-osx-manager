package master

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    tervice "github.com/gravitational/teleport/lib/service"

    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/pc-vbox-comm/masterctrl"
    "github.com/stkim1/pc-core/beacon"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/extlib/pcssh/sshadmin"
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
    clt, err := sshadmin.OpenAdminClientWithAuthService(b.PocketConfig)
    if err != nil {
        return errors.WithStack(err)
    }
    defer clt.Close()
    token, err := sshadmin.GenerateNodeInviationWithTTL(clt, sshadmin.MaxInvitationTLL)
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
        nodeC   = make(chan service.Event)
        statC   = make(chan service.Event)
    )
    appLife.RegisterServiceWithFuncs(
        operation.ServiceBeaconMaster,
        func() error {
            var (
                beaconRoute *beaconEventRoute  = &beaconEventRoute{
                    ServiceSupervisor: appLife,
                    PocketConfig:      tcfg,
                }
                readyMarker  = map[string]bool{
                    ivent.IventPcsshProxyInstanceSpawn: false,
                    ivent.IventVboxCtrlInstanceSpawn:   false,
                }
                readyChecker = func(marker map[string]bool) bool {
                    for k := range marker {
                        if !marker[k] {
                            return false
                        }
                    }
                    return true
                }
                failtimout         *time.Ticker = time.NewTicker(time.Minute)
                timer             *time.Ticker = time.NewTicker(time.Second)
                beaconMan  beacon.BeaconManger = nil
                vmctrl     masterctrl.VBoxMasterControl
                err        error               = nil
            )

            // wait pre-requisites to start
            for {
                select {
                    // fail to start service after one minute
                    case <- failtimout.C: {
                        failtimout.Stop()
                        timer.Stop()
                        beaconRoute.terminate()
                        log.Errorf("[AGENT] fail to start agent service")
                        return errors.Errorf("[AGENT] fail to start agent service")
                    }
                    // waiting teleport to start
                    case <- teleC: {
                        readyMarker[ivent.IventPcsshProxyInstanceSpawn] = true
                        log.Infof("[AGENT] pcssh ready")

                        if readyChecker(readyMarker) {
                            goto buildagent
                        }
                    }
                    // wait for vbox control
                    case vc := <- vboxC: {
                        vbc, ok := vc.Payload.(*ivent.VboxCtrlBrcstObj)
                        if vbc != nil && ok {
                            vmctrl, ok = vbc.VBoxMasterControl.(masterctrl.VBoxMasterControl)
                            if !ok {
                                return errors.Errorf("[AGENT] (ERR) invalid VBoxMasterControl type")
                            }
                            readyMarker[ivent.IventVboxCtrlInstanceSpawn] = true
                            log.Infof("[AGENT] vbox core ready")

                            if readyChecker(readyMarker) {
                                goto buildagent
                            }
                        } else {
                            return errors.Errorf("[AGENT] (ERR) invalid VBoxMasterControl type")
                        }
                    }
                }
            }

            buildagent:
            // stop failtimout
            failtimout.Stop()
            // beacon manager
            beaconMan, err = beacon.NewBeaconManagerWithFunc(
                clusterID,
                vmctrl,
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

            appLife.BroadcastEvent(service.Event{
                Name:ivent.IventBeaconManagerSpawn,
                Payload:beaconMan})
            log.Debugf("[AGENT] starting agent service...")

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
                    case <-timer.C: {
                        err = beaconMan.TransitionWithTimestamp(time.Now())
                        if err != nil {
                            log.Debug(err.Error())
                        }
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
                    // network monitor event
                    case <- netC: {
                        // TODO update primary address
                        log.Debugf("[AGENT] Host Address changed")
                    }
                    // package node list service
                    case <- nodeC: {
                        nodeList := beaconMan.RegisteredNodesList()
                        appLife.BroadcastEvent(service.Event{
                            Name:ivent.IventReportNodeListResult,
                            Payload:nodeList})
                    }
                    // node status report service
                    case re := <- statC: {
                        _, ok := re.Payload.(int64)
                        if !ok {
                            appLife.BroadcastEvent(service.Event{
                                Name:    ivent.IventMonitorNodeRespBeacon,
                                Payload: errors.Errorf("inaccurate timestamp")})
                        }
                        // need unregistered node, registered node, bounded node
                        regNodes := beaconMan.RegisteredNodesList()
                        appLife.BroadcastEvent(service.Event{
                            Name:ivent.IventMonitorNodeRespBeacon,
                            Payload:regNodes})
                    }
                }
            }
            return nil
        },
        service.BindEventWithService(ucast.EventBeaconCoreLocationReceive, beaconC),
        service.BindEventWithService(mcast.EventBeaconCoreSearchReceive,   searchC),
        service.BindEventWithService(ivent.IventNetworkAddressChange,      netC),
        service.BindEventWithService(ivent.IventReportNodeListRequest,     nodeC),
        service.BindEventWithService(ivent.IventMonitorNodeReqStatus,      statC),

        // service readiness checker
        service.BindEventWithService(ivent.IventPcsshProxyInstanceSpawn,   teleC),
        service.BindEventWithService(ivent.IventVboxCtrlInstanceSpawn,     vboxC))

    return nil
}
