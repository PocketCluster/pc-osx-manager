package regnode

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"

    "github.com/stkim1/pc-core/beacon"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/rasker"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/service/ivent"
)

const (
    iventNodeRegisterCandid string = "ivent.node.register.candidate"
    iventNodeRegisterStop   string = "ivent.node.register.stop"
    raskerNodeRegisterCycle string = "rasker.node.register.cycle"
)

func InitNodeRegisterCycle(appLife rasker.RouteTasker, feeder route.ResponseFeeder) error {
    return appLife.GET(routepath.RpathNodeRegStart(), func(_, rpath, _ string) error {
        var (
            candidC = make(chan service.Event)
            stopC   = make(chan service.Event)
            beaconC = make(chan service.Event)
            searchC = make(chan service.Event)
            bManC   = make(chan service.Event)
        )
        appLife.RegisterServiceWithFuncs(
            raskerNodeRegisterCycle,
            func() error {
                var (
                    rpUnregNodeList = routepath.RpathNodeUnregList()
                    rpNodeRegCnfrm  = routepath.RpathNodeRegConfirm()
                    rpNodeRegCandid = routepath.RpathNodeRegCandiate()
                    regMan beacon.RegisterManger = nil
                    beaconMan beacon.BeaconManger = nil
                    rptTick *time.Ticker = nil
                    err error = nil
                )

                // we need beacon master here

                log.Debugf("[REGISTER] request agent instance")
                appLife.BroadcastEvent(service.Event{Name:ivent.IventLiveBeaconManagerReq})
                be := <- bManC
                bm, ok := be.Payload.(beacon.BeaconManger)
                if bm != nil && ok {
                    beaconMan = bm
                } else {
                    return feedGetError(feeder, rpath, "node-reg-start", errors.Errorf("[REGISTER] invalid beacon manager"))
                }

                regMan, err = beacon.NewNodeRegisterManager(beaconMan)
                if err != nil {
                    return feedGetError(feeder, rpath, "node-reg-start", err)
                }
                rptTick = time.NewTicker(time.Second * 4)
                log.Debugf("[REGISTER] started")
                for {
                    select {
                        case <- appLife.StopChannel(): {
                            log.Debugf("[REGISTER] stopped")
                            rptTick.Stop()
                            return nil
                        }
                        case <- stopC: {
                            log.Debugf("[REGISTER] stopped")
                            rptTick.Stop()
                            return nil
                        }
                        case <- candidC: {
                            if regMan.IsRegistering() {
                                continue
                            }
                            if rerr := regMan.RegisterMonitoredNodes(time.Now()); rerr == nil {
                                log.Debug("[REGISTER] node registration went ok")
                                if frr := feedGetOkMessage(feeder, rpNodeRegCandid, "node-reg-candidate"); frr != nil {
                                    log.Errorf("[REGISTER] node registration success feedback fail %v", frr.Error())
                                }
                            } else {
                                log.Errorf("[REGISTER] node registration failed %v", rerr.Error())
                                if frr := feedGetError(feeder, rpNodeRegCandid, "node-reg-candidate", errors.Errorf("[REGISTER] invalid beacon manager")); frr != nil {
                                    log.Errorf("[REGISTER] node registration failure feedback fail %v", frr.Error())
                                }
                            }
                        }
                        case ts := <- rptTick.C: {
                            if regMan.IsRegistering() {
                                if regMan.IsAllNodeRegistered(ts) {
                                    log.Debugf("[REGISTER] stopped. every node is all registered")
                                    rptTick.Stop()
                                    return feedGetOkMessage(feeder, rpNodeRegCnfrm, "node-reg-confirm")
                                }
                                if regMan.IsRegistrationTimedOut(ts) {
                                    log.Debugf("[REGISTER] stopped. registration timeout")
                                    rptTick.Stop()
                                    return feedGetError(feeder, rpNodeRegCnfrm, "node-reg-confirm", errors.Errorf("[REGISTER] confirmation timeout failure. some node does not report."))
                                }
                            } else {
                                // FIXME : this is an easy fix. we need a serious fix on register manager at beacon
                                list := regMan.UnregisteredNodeList(ts)
                                if lrr := feedGetMessage(feeder, rpUnregNodeList, "node-unreged", "unreged-list", list); lrr != nil {
                                    log.Debugf("[REGISTER] unregistered node report error : %v", lrr.Error())
                                }
                            }
                        }
                        case b := <-beaconC: {
                            if !regMan.IsRegistering() {
                                continue
                            }
                            bp, ok := b.Payload.(ucast.BeaconPack)
                            if ok {
                                brr := regMan.GuideNodeRegistrationWithBeacon(bp, time.Now())
                                if brr != nil {
                                    log.Debugf("[REGISTER] BEACON-RX Error : %v", brr)
                                }
                            }
                        }
                        case s := <-searchC: {
                            if regMan.IsRegistering() {
                                continue
                            }
                            cp, ok := s.Payload.(mcast.CastPack)
                            if ok {
                                srr := regMan.MonitoringMasterSearchData(cp, time.Now())
                                if srr != nil {
                                    log.Debugf("[REGISTER] SEARCH-RX Error : %v", srr)
                                }
                            }
                        }
                    }
                }
            },
            service.BindEventWithService(iventNodeRegisterCandid,              candidC),
            service.BindEventWithService(iventNodeRegisterStop,                stopC),
            service.BindEventWithService(ucast.EventBeaconCoreLocationReceive, beaconC),
            service.BindEventWithService(mcast.EventBeaconCoreSearchReceive,   searchC),
            service.BindEventWithService(ivent.IventLiveBeaconManagerRslt,     bManC))
        return nil
    })
}
