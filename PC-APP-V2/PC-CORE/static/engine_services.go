package main

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/coreos/etcd/embed"

    "github.com/stkim1/pc-core/extlib/registry"
    "github.com/stkim1/pc-core/event/operation"
)

const (
    iventBeaconManagerSpawn string  = "ivent.beacon.manager.spawn"
    iventSwarmInstanceSpawn string  = "ivent.swarm.instance.spawn"
)

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
                    log.Debugf("[ETCD] server is ready to run")
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
                        log.Debugf("[ETCD] error : %v", err)
                    }
                    case <- a.StopChannel(): {
                        etcd.Close()
                        log.Debugf("[ETCD] server shuts down")
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
            log.Debugf("[REGISTRY] server start successful")

            // wait for service to stop
            <- a.StopChannel()
            reg.Stop(time.Second)
            log.Debugf("[REGISTRY] server exit")
            return nil
        })
    return nil
}