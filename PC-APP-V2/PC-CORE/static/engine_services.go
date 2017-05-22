package main

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/pc-core/beacon"
)

const (
    coreFeedbackSearch = "feedback_search"
    coreFeedbackBeacon = "feedback_beacon"
    coreServiceBeacon  = "service_beacon"
)

func initSearchCatcher(a *mainLife) error {
    // TODO : use network interface
    catcher, err := mcast.NewSearchCatcher("en0")
    if err != nil {
        return errors.WithStack(err)
    }

    a.RegisterServiceFunc(func() error {
        log.Debugf("NewSearchCatcher :: MAIN BEGIN")
        for {
            select {
                case <-a.StopChannel(): {
                    catcher.Close()
                    log.Debugf("NewSearchCatcher :: MAIN CLOSE")
                    return nil
                }
                case cp := <-catcher.ChRead: {
                    a.BroadcastEvent(Event{Name:coreFeedbackSearch, Payload:cp})
                }
            }
        }
        return nil
    })

    return nil
}

func initBeaconLoator(a *mainLife) error {
    belocat, err := ucast.NewBeaconLocator()
    if err != nil {
        return errors.WithStack(err)
    }

    // beacon locator read
    a.RegisterServiceFunc(func() error {
        log.Debugf("NewBeaconLocator READ :: MAIN BEGIN")
        for {
            select {
                case <- a.StopChannel(): {
                    belocat.Close()
                    log.Debugf("NewBeaconLocator READ :: MAIN CLOSE")
                    return nil
                }
                case bp := <- belocat.ChRead: {
                    a.BroadcastEvent(Event{Name:coreFeedbackBeacon, Payload:bp})
                }
            }
        }
        return nil
    })

    // beacon locator write
    beaconC := make(chan Event)
    a.WaitForEvent(coreServiceBeacon, beaconC, make(chan struct{}))
    a.RegisterServiceFunc(func() error {
        log.Debugf("NewBeaconLocator WRITE :: MAIN BEGIN")
        for {
            select {
                case <- a.StopChannel(): {
                    log.Debugf("NewBeaconLocator WRITE :: MAIN CLOSE")
                    return nil
                }
                case b := <- beaconC: {
                    bs, ok := b.Payload.(ucast.BeaconSend)
                    if ok {
                        belocat.Send(bs.Host, bs.Payload)
                    }
                }
            }
        }
        return nil
    })

    return nil
}

func initMasterAgentService(clusterID string, a *mainLife) error {
    var (
        beaconC = make(chan Event)
        searchC = make(chan Event)
    )
    a.WaitForEvent(coreFeedbackBeacon, beaconC, make(chan struct{}))
    a.WaitForEvent(coreFeedbackSearch, searchC, make(chan struct{}))

    a.RegisterServiceFunc(func() error {
        var (
            beaconMan beacon.BeaconManger = nil
            err error = nil
            timer = time.NewTicker(time.Second)
        )
        beaconMan, err = beacon.NewBeaconManagerWithFunc(clusterID, func(host string, payload []byte) error {
            log.Debugf("[BEACON-SLAVE] Host %v", host)
            a.BroadcastEvent(Event{
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
        defer timer.Stop()

        log.Debugf("[AGENT] starting agent service...")

        for {
            select {
                case <-a.StopChannel(): {
                    err = beaconMan.Shutdown()
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
                            log.Debug(err.Error())
                        }
                    }
                }
                case s := <-searchC: {
                    cp, ok := s.Payload.(mcast.CastPack)
                    if ok {
                        err = beaconMan.TransitionWithSearchData(cp, time.Now())
                        if err != nil {
                            log.Debug(err.Error())
                        }
                    }
                }
                case t := <-timer.C: {
                    err = beaconMan.TransitionWithTimestamp(t)
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