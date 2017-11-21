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
                            rptTick.Stop()
                            log.Debugf("[REGISTER] stopped")
                            return nil
                        }
                        case <- stopC: {
                            rptTick.Stop()
                            log.Debugf("[REGISTER] stopped")
                            return nil
                        }
                        case <- candidC: {

                        }
                        case ts := <- rptTick.C: {
                            list := regMan.UnregisteredNodeList(ts)
                            lrr := feedGetMessage(feeder, rpUnregNodeList, "node-unreged", "unreged-list", list)
                            if lrr != nil {
                                log.Debugf("[REGISTER] unregistered node report error : %v", lrr)
                            }
                        }
                        case b := <-beaconC: {
                            bp, ok := b.Payload.(ucast.BeaconPack)
                            if ok {
                                brr := regMan.GuideNodeRegistrationWithBeacon(bp, time.Now())
                                if brr != nil {
                                    log.Debugf("[REGISTER] BEACON-RX Error : %v", brr)
                                }
                            }
                        }
                        case s := <-searchC: {
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
