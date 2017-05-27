package main

import (

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/pc-core/service"
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
                    a.BroadcastEvent(service.Event{Name:coreFeedbackSearch, Payload:cp})
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
                    a.BroadcastEvent(service.Event{Name:coreFeedbackBeacon, Payload:bp})
                }
            }
        }
        return nil
    })

    // beacon locator write
    beaconC := make(chan service.Event)
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
