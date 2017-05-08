package main

import (
    "os"
    "net"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "gopkg.in/vmihailenco/msgpack.v2"

    "github.com/stkim1/pc-node-agent/dhcp"
    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/udpnet/ucast"
)

const (
    NodeServiceSearch = "node_service_search"
    NodeServiceBeacon = "node_service_beacon"
    NodeServiceAgent  = "node_service_agent"
)

func initDhcpListner(app *PocketApplication) error {

    app.RegisterFunc(func () error {
        log.Info("[DHCP] dhcp listner started...")

        buf := make([]byte, 20480)
        dhcpEvent := &dhcp.DhcpEvent{}

        // firstly clear off previous socket
        os.Remove(dhcp.DHCPEventSocketPath)
        listen, err := net.ListenUnix("unix", &net.UnixAddr{dhcp.DHCPEventSocketPath, "unix"})
        if err != nil {
            return errors.WithStack(err)
        }
        defer os.Remove(dhcp.DHCPEventSocketPath)
        defer listen.Close()

        for {
            select {
                case conn, err := listen.AcceptUnix(): {
                    if err != nil {
                        log.Error(errors.WithStack(err))
                        continue
                    }
                    count, err := conn.Read(buf)
                    if err != nil {
                        log.Error(errors.WithStack(err))
                        continue
                    }
                    err = msgpack.Unmarshal(buf[0:count], dhcpEvent)
                    if err != nil {
                        log.Error(errors.WithStack(err))
                        continue
                    }

                    //log.Info(spew.Sdump(dhcpEvent))

                    err = conn.Close()
                    if err != nil {
                        log.Error(errors.WithStack(err))
                        continue
                    }
                }
            }
        }
    })

    return nil
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

func initAgentService(app *PocketApplication) error {
    return nil
}

func main() {
    var (
        err error = nil
        app *PocketApplication
    )
    app = NewPocketApplication()

    // dhcp listner
    err = initDhcpListner(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // search service
    initSearchService(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // beacon service
    initBeaconService(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // application
    err = app.Start()
    if err != nil {
        log.Panic(errors.WithStack(err))
    }
    app.Wait()
}

