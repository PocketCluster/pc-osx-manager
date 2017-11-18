package container

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/coreos/etcd/embed"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/extlib/registry"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/service/ivent"
)

func InitDiscoveryService(appLife service.ServiceSupervisor, config *embed.PocketConfig) error {
    appLife.RegisterServiceWithFuncs(
        operation.ServiceDiscoveryServer,
        func() error {
            etcd, err := embed.StartPocketEtcd(config)
            if err != nil {
                return errors.WithStack(err)
            }
            // startup preps
            select {
                case <-etcd.Server.ReadyNotify(): {
                    log.Debugf("[DSCVRY] server is ready to run")
                }
                case <-time.After(120 * time.Second): {
                    etcd.Server.Stop() // trigger a shutdown

                    appLife.BroadcastEvent(service.Event{
                        Name:    ivent.IventInternalSpawnError,
                        Payload: errors.Errorf("[DSCVRY] Server took too long to start!"),
                    })
                    return errors.Errorf("[DSCVRY] Server took too long to start!")
                }
            }

            // report successful start up
            appLife.BroadcastEvent(service.Event{Name:ivent.IventDiscoveryInstanceSpwan})

            // until server goes down, errors and stop signal will be constantly checked
            for {
                select {
                    case err = <-etcd.Err(): {
                        log.Debugf("[DSCVRY] error : %v", err)
                    }
                    case <- appLife.StopChannel(): {
                        etcd.Close()
                        log.Debugf("[DSCVRY] server shuts down")
                        return nil
                    }
                }
            }
            return nil
        })

    return nil
}

func InitRegistryService(appLife service.ServiceSupervisor, config *registry.PocketRegistryConfig) error {
    appLife.RegisterServiceWithFuncs(
        operation.ServiceOrchstRegistry,
        func() error {
            reg, err := registry.NewPocketRegistry(config)
            if err != nil {
                appLife.BroadcastEvent(service.Event{
                    Name:    ivent.IventInternalSpawnError,
                    Payload: errors.WithStack(err),
                })
                return errors.WithStack(err)
            }
            err = reg.Start()
            if err != nil {
                appLife.BroadcastEvent(service.Event{
                    Name:    ivent.IventInternalSpawnError,
                    Payload: errors.WithStack(err),
                })
                return errors.WithStack(err)
            }

            // report successful start up
            appLife.BroadcastEvent(service.Event{Name:ivent.IventRegistryInstanceSpawn})
            log.Debugf("[REGISTRY] server start successful")

            // wait for service to stop
            <- appLife.StopChannel()
            err = reg.Stop(time.Second)
            log.Debugf("[REGISTRY] server exit. Error : %v", err)
            return errors.WithStack(err)
        })
    return nil
}
