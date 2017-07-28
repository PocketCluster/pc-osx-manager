package main

import (
    "net"
    "os"

    "gopkg.in/vmihailenco/msgpack.v2"
    sysd "github.com/coreos/go-systemd/dbus"
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/service"
    "github.com/stkim1/pc-node-agent/pcssh/sshcfg"
    "github.com/stkim1/pc-node-agent/pcssh/sshproc"
    "github.com/stkim1/pc-node-agent/utils/dhcp"
)

const (
    iventNodeDHCPFeedback string    = "ivent.node.dhcp.feedback"
    systemdDockerServiceUnit string = "docker.service"

    servicePcsshInit string         = "service.pcssh.init"
    servicePcsshStart string        = "service.pcssh.start"
)

func initDhcpListner(app service.AppSupervisor) error {
    // firstly clear off previous socket
    os.Remove(dhcp.PocketDHCPEventSocketPath)
    dhcpListener, err := net.ListenUnix("unix", &net.UnixAddr{Name:dhcp.PocketDHCPEventSocketPath, Net:"unix"})
    if err != nil {
        return errors.WithStack(err)
    }

    app.RegisterServiceWithFuncs(
        func () error {
            var (
                buf []byte = make([]byte, 20480)
                dhcpEvent = &dhcp.PocketDhcpEvent{}
            )

            log.Debugf("[DHCP] starting dhcp listner...")

            // TODO : how do we stop this?
            for {
                conn, err := dhcpListener.AcceptUnix()
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

                app.BroadcastEvent(service.Event{Name: iventNodeDHCPFeedback, Payload: dhcpEvent})

                err = conn.Close()
                if err != nil {
                    log.Error(errors.WithStack(err))
                    continue
                }
            }

            return nil
        },
        func(_ func(interface{})) error {
            dhcpListener.Close()
            os.Remove(dhcp.PocketDHCPEventSocketPath)
            log.Debugf("[DHCP] close dhcp listner...")
            return nil
        },
    )

    return nil
}

func initTeleportNodeService(app service.AppSupervisor) error {
    app.RegisterNamedServiceWithFuncs(
        servicePcsshInit,
        func() error{
            var (
                pcsshNode *sshproc.EmbeddedNodeProcess = nil
                err error = nil
            )
            // restart teleport
            cfg, err := sshcfg.MakeNodeConfig(slcontext.SharedSlaveContext(), true)
            if err != nil {
                return errors.WithStack(err)
            }
            pcsshNode, err = sshproc.NewEmbeddedNodeProcess(app, cfg)
            if err != nil {
                log.Errorf(err.Error())
                return errors.WithStack(err)
            }

            // execute docker engine cert acquisition before SSH node start
            // TODO : create a waitforevent channel and restart docker engine accordingly
            err = pcsshNode.AcquireEngineCertificate(slcontext.DockerEnvironemtPostProcess)
            if err != nil {
                return errors.WithStack(err)
            }

            err = pcsshNode.StartNodeSSH()
            if err != nil {
                return errors.WithStack(err)
            }
            log.Debugf("\n\n(INFO) teleport node started success!\n")

            return nil

            // restart docker engine
            // TODO : FIX /opt/gopkg/src/github.com/godbus/dbus/conn.go:345 send on closed channel
            conn, err := sysd.NewSystemdConnection()
            if err != nil {
                log.Errorf(err.Error())
            } else {
                did, err := conn.RestartUnit(systemdDockerServiceUnit, "replace", nil)
                if err != nil {
                    log.Errorf(err.Error())
                } else {
                    conn.Close()
                    log.Debugf("\n\n(INFO) docker engin restart success! ID %d\n", did)
                }
            }

            return nil
        },
        func(_ func(interface{})) error {
            return nil
        })

    app.RegisterNamedServiceWithFuncs(
        servicePcsshStart,
        func() error{
            var (
                pcsshNode *sshproc.EmbeddedNodeProcess = nil
                err error = nil
            )
            // restart teleport
            cfg, err := sshcfg.MakeNodeConfig(slcontext.SharedSlaveContext(), true)
            if err != nil {
                return errors.WithStack(err)
            }
            pcsshNode, err = sshproc.NewEmbeddedNodeProcess(app, cfg)
            if err != nil {
                return errors.WithStack(err)
            }

            err = pcsshNode.StartNodeSSH()
            if err != nil {
                return errors.WithStack(err)
            }
            log.Debugf("\n\n(INFO) teleport node started success!\n")

            return nil

            // restart docker engine
            // TODO : FIX /opt/gopkg/src/github.com/godbus/dbus/conn.go:345 send on closed channel
            conn, err := sysd.NewSystemdConnection()
            if err != nil {
                log.Errorf(err.Error())
            } else {
                did, err := conn.RestartUnit(systemdDockerServiceUnit, "replace", nil)
                if err != nil {
                    return errors.WithStack(err)
                } else {
                    conn.Close()
                    log.Debugf("\n\n(INFO) docker engin restart success! ID %d\n", did)
                }
            }

            return nil
        },
        func(_ func(interface{})) error {
            return nil
        })

    return nil
}