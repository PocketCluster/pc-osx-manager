package main

import (
    "time"

    "github.com/gravitational/teleport/embed"
    tservice "github.com/gravitational/teleport/lib/service"
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/pc-core/beacon"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/model"
)

type beaconEventRoute struct {
    service.ServiceSupervisor
    *tservice.PocketConfig
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
    return errors.WithStack(slave.SetSlaveID(token))
}

func (b *beaconEventRoute) BeaconEventResurrect(slaves []model.SlaveNode) error {
    return nil
}

func (b *beaconEventRoute) BeaconEventTranstion(state beacon.MasterBeaconState, slave *model.SlaveNode, ts time.Time, transOk bool) error {
    log.Debugf("(INFO) BeaconEventTranstion [%v] %v | %v", ts, state, slave.SlaveUUID)
    return nil
}

func (b *beaconEventRoute) BeaconEventDiscard(slave *model.SlaveNode) error {
    return nil
}

func (b *beaconEventRoute) BeaconEventShutdown() error {
    return nil
}

func initMasterAgentService(a *mainLife, clusterID string, tcfg *tservice.PocketConfig) error {
    var (
        beaconC = make(chan service.Event)
        searchC = make(chan service.Event)
        beaconRoute *beaconEventRoute = &beaconEventRoute{
            ServiceSupervisor:a.ServiceSupervisor,
            PocketConfig: tcfg,
        }
    )
    a.WaitForEvent(coreFeedbackBeacon, beaconC, make(chan struct{}))
    a.WaitForEvent(coreFeedbackSearch, searchC, make(chan struct{}))

    a.RegisterServiceFunc(func() error {
        var (
            beaconMan beacon.BeaconManger = nil
            err error = nil
            timer = time.NewTicker(time.Second)
        )
        beaconMan, err = beacon.NewBeaconManagerWithFunc(clusterID, beaconRoute, func(host string, payload []byte) error {
            log.Debugf("[BEACON-SEND-SLAVE] Host %v", host)
            a.BroadcastEvent(service.Event{
                Name: coreServiceBeacon,
                Payload:ucast.BeaconSend{
                    Host:       host,
                    Payload:    payload,
                },
            })
            return nil
        })
        if err != nil {
            return errors.WithStack(err)
        }

        log.Debugf("[AGENT] starting agent service...")

        for {
            select {
            case <-a.StopChannel(): {
                err = beaconMan.Shutdown()
                if err != nil {
                    log.Debug(err.Error())
                }
                err = beaconRoute.close()
                if err != nil {
                    log.Debug(err.Error())
                }
                timer.Stop()
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
                continue
            }
            }
        }
        return nil
    })

    return nil
}

