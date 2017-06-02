package main

import (
    "fmt"
    "net"

    "github.com/pkg/errors"
    "github.com/miekg/dns"
    "github.com/stkim1/pc-node-agent/service"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pcrypto"
)

const (
    localPocketMasterName string = "pc-master."
    fqdnPocketMasterName  string = localPocketMasterName + pcrypto.FormFQDNClusterID
)

func failWithRcode(w dns.ResponseWriter, r *dns.Msg, rCode int) {
    m := new(dns.Msg)
    m.SetRcode(r, rCode)
    w.WriteMsg(m)
}

func locaNameServe(w dns.ResponseWriter, req *dns.Msg) {
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
    m := new(dns.Msg)
    m.Id = req.Id

    switch qtype {

        case dns.TypeA: {

            remoteIP4 := remoteIP.To4();
            ma, err := slcontext.SharedSlaveContext().GetMasterAgent()

            if err == nil && remoteIP4 != nil {

                fqdn := fmt.Sprintf(fqdnPocketMasterName, ma)
                maddr, err := slcontext.SharedSlaveContext().GetMasterIP4Address()

                if err == nil && (question.Name == localPocketMasterName || question.Name == fqdn) {
                    rr := new(dns.A)
                    rr.Hdr = dns.RR_Header{
                        Name:      question.Name,
                        Rrtype:    question.Qtype,
                        Class:     dns.ClassINET,
                        Ttl:       10,
                    }
                    rr.A = net.ParseIP(maddr)
                    m.Answer = []dns.RR{rr}
                }

            }
        }
    }

    m.Question = req.Question
    m.Response = true
    m.Authoritative = true
    w.WriteMsg(m)
}

func initDNSService(app service.AppSupervisor) error {
    dns.HandleFunc(".", locaNameServe)
    var (
        udpServer = &dns.Server {
            Addr:    "127.0.0.1:53",
            Net:     "udp",
        }
        udpPacketConn *net.UDPConn = nil
        udpAddr *net.UDPAddr = nil
        err error = nil
    )
    udpAddr, err = net.ResolveUDPAddr(udpServer.Net, udpServer.Addr)
    if err != nil {
        return errors.WithStack(err)
    }
    udpPacketConn, err = net.ListenUDP(udpServer.Net, udpAddr)
    if err != nil {
        return errors.WithStack(err)
    }
    udpServer.PacketConn = udpPacketConn

    app.RegisterFunc(func() error {
        return errors.WithStack(udpServer.ActivateAndServe())
    })
    app.OnExit(func(interface{}) {
        udpServer.Shutdown()
        dns.HandleRemove(".")
    })

    return nil
}