package dns

import (
    "net"
    "strings"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/miekg/dns"
    "github.com/pkg/errors"

    "github.com/stkim1/pc-core/beacon"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/service/ivent"
)

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

func locaNodeName(beaconMan beacon.BeaconManger, w dns.ResponseWriter, req *dns.Msg) {
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

    remoteIP := w.RemoteAddr().(*net.UDPAddr).IP
    remoteIP4 := remoteIP.To4()

    m := new(dns.Msg)
    m.Id = req.Id
    m.Question = req.Question
    m.Response = true

    switch qtype {
        case dns.TypeA: {
            nn := strings.Trim(question.Name, " .\t\r\n")
            addr, err := beaconMan.AddressForName(nn)

            if remoteIP4 != nil && err == nil {
                rr := new(dns.A)
                rr.Hdr = dns.RR_Header{
                    Name:      question.Name,
                    Rrtype:    question.Qtype,
                    Class:     dns.ClassINET,
                    Ttl:       10,
                }
                rr.A = net.ParseIP(addr)
                m.Answer = []dns.RR{rr}
                m.Authoritative = true
                w.WriteMsg(m)

                log.Debugf("[NAME-SERVICE] '%s' (%s for %s). Error : %v", nn, addr, question.Name, err)
                return
            }
        }
    }

    // libresolv continues to the next server when it receives
    // an invalid referral response. See golang.org/issue/15434.
    m.Rcode = dns.RcodeSuccess
    m.Authoritative = false
    m.RecursionAvailable = false
    w.WriteMsg(m)
}

func InitPocketNameService(appLife service.ServiceSupervisor, clusterID string) error {
    nameServerC := make(chan service.Event)
    appLife.RegisterServiceWithFuncs(
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
        service.BindEventWithService(ivent.IventNameServerInstanceSpawn, nameServerC))

    beaconManC := make(chan service.Event)
    appLife.RegisterServiceWithFuncs(
        operation.ServiceInternalNodeNameControl,
        func() error {
            var (
                udpServer = &dns.Server {
                    Addr:    "127.0.0.1:10059",
                    Net:     "udp",
                }
                failtimeout   *time.Ticker        = time.NewTicker(time.Minute)
                udpPacketConn *net.UDPConn        = nil
                udpAddr       *net.UDPAddr        = nil
                beaconMan     beacon.BeaconManger = nil
                err           error               = nil
            )

            // monitor beacon agent
            select {
                case <-failtimeout.C: {
                    // cleanup resources
                    failtimeout.Stop()
                    udpServer.Shutdown()
                    return errors.Errorf("[ERR] unable to detect beacon agent")
                }
                // wait for beacon manager to come up
                case be := <- beaconManC: {
                    bm, ok := be.Payload.(beacon.BeaconManger)
                    if bm != nil && ok {
                        beaconMan = bm
                    } else {
                        return errors.Errorf("[ERR] invalid beacon manager type")
                    }
                }
            }

            failtimeout.Stop()
            log.Infof("[NAME-SERVICE] service ready to start...")
            // setup dns handler
            dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
                locaNodeName(beaconMan, w, req)
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
            appLife.BroadcastEvent(service.Event{
                Name:ivent.IventNameServerInstanceSpawn,
                Payload:udpServer})

            // wait for stop event
            <- appLife.StopChannel()
            log.Debugf("[NAME-SERVICE] service shutting down")
            err = udpServer.Shutdown()
            udpServer = nil
            dns.HandleRemove(".")
            udpPacketConn = nil
            udpAddr = nil
            return errors.WithStack(err)
        },
        service.BindEventWithService(ivent.IventBeaconManagerSpawn, beaconManC))

    return nil
}
