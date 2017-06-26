package main

import (
    "fmt"
    "net"
    "strings"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/coreos/etcd/embed"
    "github.com/miekg/dns"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-core/beacon"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/extlib/registry"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/service"
    swarmemb "github.com/stkim1/pc-core/extlib/swarm"
)

const (
    iventBeaconManagerSpawn string  = "ivent.beacon.manager.spawn"
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
            err = reg.Stop(time.Second)
            log.Debugf("[REGISTRY] server exit. Error : %v", err)
            return errors.WithStack(err)
        })
    return nil
}

func initSwarmService(a *mainLife) error {
    const (
        iventSwarmInstanceSpawn string  = "ivent.swarm.instance.spawn"
    )
    var (
        swarmSrvC = make(chan service.Event)
    )
    a.RegisterServiceWithFuncs(
        operation.ServiceSwarmEmbeddedOperation,
        func() error {
            var (
                swarmsrv *swarmemb.SwarmService = nil
            )
            for {
                select {
                    case se := <- swarmSrvC: {
                        srv, ok := se.Payload.(*swarmemb.SwarmService)
                        if ok {
                            swarmsrv = srv
                        }
                    }
                    case <- a.StopChannel(): {
                        if swarmsrv != nil {
                            err := swarmsrv.Close()
                            return errors.WithStack(err)
                        }
                        return errors.Errorf("[ERR] null SWARM instance")
                    }
                }
            }
            return nil
        },
        service.BindEventWithService(iventSwarmInstanceSpawn, swarmSrvC))

    beaconManC := make(chan service.Event)
    a.RegisterServiceWithFuncs(
        operation.ServiceSwarmEmbeddedServer,
        func() error {
            be := <- beaconManC
            beaconMan, ok := be.Payload.(beacon.BeaconManger)
            if !ok {
                return errors.Errorf("[ERR] invalid beacon manager type")
            }
            ctx := context.SharedHostContext()
            caCert, err := ctx.CertAuthCertificate()
            if err != nil {
                return errors.WithStack(err)
            }
            hostCrt, err := ctx.MasterHostCertificate()
            if err != nil {
                return errors.WithStack(err)
            }
            hostPrv, err := ctx.MasterHostPrivateKey()
            if err != nil {
                return errors.WithStack(err)
            }
            swarmctx, err := swarmemb.NewContextWithCertAndKey(caCert, hostCrt, hostPrv, beaconMan)
            if err != nil {
                return errors.WithStack(err)
            }
            swarmsrv, err := swarmemb.NewSwarmService(swarmctx)
            if err != nil {
                return errors.WithStack(err)
            }
            a.BroadcastEvent(service.Event{Name:iventSwarmInstanceSpawn, Payload:swarmsrv})
            err = swarmsrv.ListenAndServeSingleHost()
            return errors.WithStack(err)
        },
        service.BindEventWithService(iventBeaconManagerSpawn, beaconManC))

    return nil
}

// --- methods for name service --- //
func failWithRcode(w dns.ResponseWriter, r *dns.Msg, rCode int) {
    m := new(dns.Msg)
    m.SetRcode(r, rCode)
    w.WriteMsg(m)
    m = nil
}

func clearNodeName(name, cfqdn string) string {
    var nn = name
    if strings.Contains(name, cfqdn) {
        nn = strings.Trim(name, cfqdn)
    }
    return strings.Trim(nn, " .\t\r\n")
}

func locaNodeName(beaconMan beacon.BeaconManger, clusterID string, w dns.ResponseWriter, req *dns.Msg) {
    if len(req.Question) != 1 {
        failWithRcode(w, req, dns.RcodeRefused)
        return
    }

    question := req.Question[0]
    qtype := question.Qtype
    if question.Qclass != dns.ClassINET {
        failWithRcode(w, req, dns.RcodeRefused)
        return
    }

    cfqdn := fmt.Sprintf(pcrypto.FormFQDNClusterID, clusterID)
    remoteIP := w.RemoteAddr().(*net.UDPAddr).IP
    m := new(dns.Msg)
    m.Id = req.Id

    switch qtype {
        case dns.TypeA: {

            if remoteIP4 := remoteIP.To4(); remoteIP4 != nil {
                nn := clearNodeName(question.Name, cfqdn)
                addr, err := beaconMan.AddressForName(nn)
                if err == nil {
                    rr := new(dns.A)
                    rr.Hdr = dns.RR_Header{
                        Name:      question.Name,
                        Rrtype:    question.Qtype,
                        Class:     dns.ClassINET,
                        Ttl:       10,
                    }
                    rr.A = net.ParseIP(addr)
                    m.Answer = []dns.RR{rr}

                    log.Debugf("[NAME-SERVICE] %s for %s", addr, nn)
                } else {
                    log.Errorf("[NAME-SERVICE] %s ", err.Error())
                }
            }

        }
    }

    m.Question = req.Question
    m.Response = true
    m.Authoritative = true
    w.WriteMsg(m)
    m = nil
}

func initPocketNameService(a *mainLife, clusterID string) error {
    const (
        iventNameServerInstanceSpawn string = "ivent.name.server.instance.spawn"
    )

    nameServerC := make(chan service.Event)
    a.RegisterServiceWithFuncs(
        operation.ServiceInternalNodeNameServer,
        func() error {
            ne := <- nameServerC
            udpServer, ok := ne.Payload.(*dns.Server)
            if !ok {
                return errors.Errorf("[ERR] invalid name service instance type")
            }
            log.Debugf("[NAME-SERVICE] start service...")
            return errors.WithStack(udpServer.ActivateAndServe())
        },
        service.BindEventWithService(iventNameServerInstanceSpawn, nameServerC))

    beaconManC := make(chan service.Event)
    a.RegisterServiceWithFuncs(
        operation.ServiceInternalNodeNameOperation,
        func() error {
            var (
                udpServer = &dns.Server {
                    Addr:    "127.0.0.1:10059",
                    Net:     "udp",
                }
                udpPacketConn *net.UDPConn = nil
                udpAddr *net.UDPAddr = nil
                err error = nil
            )
            // wait for beacon manager to come up
            be := <- beaconManC
            beaconMan, ok := be.Payload.(beacon.BeaconManger)
            if !ok {
                return errors.Errorf("[ERR] invalid beacon manager type")
            }
            dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
                locaNodeName(beaconMan, clusterID, w, req)
            })

            // spawn name server
            udpAddr, err = net.ResolveUDPAddr(udpServer.Net, udpServer.Addr)
            if err != nil {
                return errors.WithStack(err)
            }
            udpPacketConn, err = net.ListenUDP(udpServer.Net, udpAddr)
            if err != nil {
                return errors.WithStack(err)
            }
            udpServer.PacketConn = udpPacketConn

            // send udp server to operation
            a.BroadcastEvent(service.Event{Name:iventNameServerInstanceSpawn, Payload: udpServer})

            // wait for stop event
            <- a.StopChannel()
            log.Debugf("[NAME-SERVICE] service shutting down")
            err = udpServer.Shutdown()
            udpServer = nil
            dns.HandleRemove(".")
            udpPacketConn = nil
            udpAddr = nil
            return errors.WithStack(err)
        },
        service.BindEventWithService(iventBeaconManagerSpawn, beaconManC))

    return nil
}