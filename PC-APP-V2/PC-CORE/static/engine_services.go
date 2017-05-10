package main

import (
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"
    "time"
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
                case r := <-catcher.ChRead: {
                    a.BroadcastEvent(Event{Name:coreFeedbackSearch, Payload:r})
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
                case r := <- belocat.ChRead: {
                    a.BroadcastEvent(Event{Name:coreFeedbackBeacon, Payload:r.Message})
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
                    log.Debugf("NewBeaconLocator WRITE %v", bs.Host)
                    belocat.Send(bs.Host, bs.Payload)
                }
            }
            }
        }
        return nil
    })

    return nil
}

func initMasterAgentService(a *mainLife) error {
    var (
        beaconC = make(chan Event)
        searchC = make(chan Event)
    )
    a.WaitForEvent(coreFeedbackBeacon, beaconC, make(chan struct{}))
    a.WaitForEvent(coreFeedbackSearch, searchC, make(chan struct{}))

    a.RegisterServiceFunc(func() error {
        bounded := time.NewTicker(time.Second * 10)
        defer bounded.Stop()

        log.Debugf("[AGENT] starting agent service...")
        for {
            select {
                case <-a.StopChannel(): {
                    log.Debugf("[AGENT] stopping agent service...")
                    return nil
                }
                case b := <-beaconC: {
                    log.Debugf("[AGENT] beacon recieved %v", b.Payload)
                }
                case s := <-searchC: {
                    log.Debugf("[AGENT] search recieved %v", s.Payload)
                }
                case <-bounded.C: {
                    log.Debugf("[AGENT] bounded %v", time.Now())
                    a.BroadcastEvent(Event{Name: coreServiceBeacon, Payload:ucast.BeaconSend{Host:"192.168.1.152", Payload:[]byte{0x55, 0x66, 0x77, 0x88, 0x99}}})
                }
            }
        }
        return nil
    })

    return nil
}