package main

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/coreos/etcd/embed"

    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/pc-core/extlib/registry"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/event/operation"
)

const (
    eventBeaconCoreReadSearch       = "event.beacon.core.read.search"
    eventBeaconCoreReadLocation     = "event.beacon.core.read.location"
    eventBeaconCoreWriteLocation    = "event.beacon.core.write.location"
    eventBeaconManagerSpawn         = "event.beacon.manager.spawn"
    eventSwarmInstanceSpawn         = "event.swarm.instance.spawn"
)

func initSearchCatcher(a *mainLife) error {
    // TODO : use network interface
    catcher, err := mcast.NewSearchCatcher("en0")
    if err != nil {
        return errors.WithStack(err)
    }

    a.RegisterServiceWithFuncs(
        operation.ServiceBeaconCatcher,
        func() error {
            log.Debugf("NewSearchCatcher :: MAIN BEGIN")
            for {
                select {
                    case <-a.StopChannel(): {
                        catcher.Close()
                        log.Debugf("NewSearchCatcher :: MAIN CLOSE")
                        return nil
                    }
                    case cp := <-catcher.ChRead: {
                        a.BroadcastEvent(service.Event{Name:eventBeaconCoreReadSearch, Payload:cp})
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
    a.RegisterServiceWithFuncs(
        operation.ServiceBeaconLocationRead,
        func() error {
            log.Debugf("NewBeaconLocator READ :: MAIN BEGIN")
            for {
                select {
                    case <- a.StopChannel(): {
                        belocat.Close()
                        log.Debugf("NewBeaconLocator READ :: MAIN CLOSE")
                        return nil
                    }
                    case bp := <- belocat.ChRead: {
                        a.BroadcastEvent(service.Event{Name:eventBeaconCoreReadLocation, Payload:bp})
                    }
                }
            }
            return nil
        })

    // beacon locator write
    beaconC := make(chan service.Event)
    a.RegisterServiceWithFuncs(
        operation.ServiceBeaconLocationWrite,
        func() error {
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
        },
        service.BindEventWithService(eventBeaconCoreWriteLocation, beaconC))

    return nil
}

func initStorageServie(a *mainLife, config *embed.PocketConfig) error {
    a.RegisterServiceWithFuncs(
        operation.ServiceStorageProcess,
        func() error {
            etcd, err := embed.StartPocketEtcd(config)
            if err != nil {
                return errors.WithStack(err)
            }
            // startup preps
            select {
                case <-etcd.Server.ReadyNotify(): {
                    log.Printf("[ETCD] server is ready to run")
                }
                case <-time.After(120 * time.Second): {
                    etcd.Server.Stop() // trigger a shutdown
                    return errors.Errorf("[ETCD] Server took too long to start!")
                }
            }
            // until server goes down, errors and stop signal will be constantly checked
            for {
                select {
                    case err = <-etcd.Err(): {
                        log.Printf("[ETCD] error : %v", err)
                    }
                    case <- a.StopChannel(): {
                        etcd.Close()
                        log.Printf("[ETCD] server shuts down")
                        return nil
                    }
                }
            }
            return nil
        })

    return nil
}

func initRegistryService(a *mainLife, config *registry.PocketRegistryConfig) error {
    a.RegisterServiceWithFuncs(
        operation.ServiceContainerRegistry,
        func() error {
            reg, err := registry.NewPocketRegistry(config)
            if err != nil {
                return errors.WithStack(err)
            }
            err = reg.Start()
            if err != nil {
                return errors.WithStack(err)
            }

            // wait for service to stop
            <- a.StopChannel()
            reg.Stop(time.Second)
            return nil
        })
    return nil
}