package main

import (
    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/udpnet/ucast"
)

const (
    NodeServiceSearch = "node_service_search"
    NodeServiceBeacon = "node_service_beacon"
)

func main() {
    app := NewPocketApplication()

    initSearchService(app)
    initBeaconService(app)

    err := app.Start()
    if err != nil {
        log.Debug(err)
    }
    app.Wait()
}

func initSearchService(app *PocketApplication) error {
    caster, err := mcast.NewSearchCaster()
    if err != nil {
        return err
    }
    eventsC := make(chan Event)
    cancleC := make(chan struct{})

    app.WaitForEvent(NodeServiceSearch, eventsC, cancleC)
    app.RegisterFunc(func() error {
        log.Debugf("[SEARCH] starting master serach service...")

        for {
            select {
                case _ = <- eventsC: {
                    caster.Send(nil)
                }
            }
        }

        return caster.Close()
    })

    return nil
}

func initBeaconService(app *PocketApplication) error {
    beacon, err := ucast.NewBeaconAgent()
    if err != nil {
        return err
    }
    eventsC := make(chan Event)
    cancelC := make(chan struct{})

    app.WaitForEvent(NodeServiceBeacon, eventsC, cancelC)
    app.RegisterFunc(func() error {
        log.Debugf("[BEACON] starting beacon service...")
        go func() {
            for v := range beacon.ChRead {
                log.Debugf("Received message %v", v.Message)
            }
        }()

        for {
            select {
                case _ = <- eventsC: {
                    beacon.Send("", nil)
                }

            }
        }
        //return agent.Close()
        return nil
    })
    return nil
}