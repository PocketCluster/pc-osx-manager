package regnode

import (
    "encoding/json"
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

func feedPostError(feeder route.ResponseFeeder, rpath, fpath string, irr error) error {
    log.Error(irr.Error())
    data, frr := json.Marshal(route.ReponseMessage{
        fpath: {
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

func feedGetError(feeder route.ResponseFeeder, rpath, fpath string, irr error) error {
    log.Error(irr.Error())
    data, frr := json.Marshal(route.ReponseMessage{
        fpath: {
            "status": false,
            "error" : irr.Error(),
        },
    })
    // this should never happen
    if frr != nil {
        log.Error(frr.Error())
    }
    frr = feeder.FeedResponseForGet(rpath, string(data))
    if frr != nil {
        log.Error(frr.Error())
    }
    return irr
}

func InitNodeRegisterCycle(appLife rasker.RouteTasker, feeder route.ResponseFeeder) error {
    return appLife.GET(routepath.RpathNodeRegStart(), func(_, rpath, _ string) error {
        var (
            beaconC = make(chan service.Event)
            searchC = make(chan service.Event)
            bManC   = make(chan service.Event)
        )
        appLife.RegisterServiceWithFuncs(
            "",
            func() error {
                var (
                    regMan beacon.RegisterManger = nil
                    beaconMan beacon.BeaconManger = nil
                    rptTick *time.Ticker = nil
                    err error = nil
                )

                // we need beacon master here
                appLife.BroadcastEvent(service.Event{Name:ivent.IventBeaconManagerRequest})
                be := <- bManC
                bm, ok := be.Payload.(beacon.BeaconManger)
                if bm != nil && ok {
                    beaconMan = bm
                } else {
                    return feedGetError(feeder, rpath, "", errors.Errorf("[REG-MAN] invalid beacon manager"))
                }

                regMan, err = beacon.NewNodeRegisterManager(beaconMan)
                if err != nil {
                    return feedGetError(feeder, rpath, "", err)
                }
                rptTick = time.NewTicker(time.Second * 3)
                for {
                    select {
                        case <- appLife.StopChannel(): {
                            rptTick.Stop()
                            return nil
                        }
                        case <- rptTick.C: {

                        }
                        case b := <-beaconC: {
                            bp, ok := b.Payload.(ucast.BeaconPack)
                            if ok {
                                err := regMan.GuideNodeRegistrationWithBeacon(bp, time.Now())
                                if err != nil {
                                    log.Debugf("[AGENT] BEACON-RX Error : %v", err)
                                }
                            }
                        }
                        case s := <-searchC: {
                            cp, ok := s.Payload.(mcast.CastPack)
                            if ok {
                                err := regMan.MonitoringMasterSearchData(cp, time.Now())
                                if err != nil {
                                    log.Debugf("[AGENT] SEARCH-RX Error : %v", err)
                                }
                            }
                        }
                    }
                }
            },
            service.BindEventWithService(ucast.EventBeaconCoreLocationReceive, beaconC),
            service.BindEventWithService(mcast.EventBeaconCoreSearchReceive,   searchC),
            service.BindEventWithService(ivent.IventBeaconManagerSpawn,        bManC))
        return nil
    })
}
