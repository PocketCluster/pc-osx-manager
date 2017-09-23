package main

import (
    "fmt"
    "net"

    "github.com/pkg/errors"
    "github.com/miekg/dns"
    "github.com/stkim1/findgate"
    "github.com/stkim1/pc-node-agent/service"
    "github.com/stkim1/pc-vbox-core/crcontext"
)

const (
    localPocketMasterName string = "pc-master."
    fqdnPocketMasterName  string = "pc-master.%s.cluster.pocketcluster.io."
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
            var (
                remoteIP4 = remoteIP.To4()
                cid       = crcontext.SharedCoreContext().CoreClusterID()
            )

            if remoteIP4 != nil && len(cid) != 0 {

                fqdn := fmt.Sprintf(fqdnPocketMasterName, cid)
                if question.Name == localPocketMasterName || question.Name == fqdn {

                    // external interface
                    maddr, merr := crcontext.SharedCoreContext().GetMasterIP4ExtAddr()
                    if merr == nil {
                        ext := new(dns.A)
                        ext.Hdr = dns.RR_Header{
                            Name:      question.Name,
                            Rrtype:    question.Qtype,
                            Class:     dns.ClassINET,
                            Ttl:       64,
                        }
                        ext.A = net.ParseIP(maddr)
                        m.Answer = append(m.Answer, ext)
                    }

                    // internal interface
                    iaddr, ierr := findgate.FindIPv4GatewayWithInterface("eth0")
                    if ierr == nil {
                        inr := new(dns.A)
                        inr.Hdr = dns.RR_Header{
                            Name:      question.Name,
                            Rrtype:    question.Qtype,
                            Class:     dns.ClassINET,
                            Ttl:       32,
                        }
                        inr.A = net.ParseIP(iaddr[0].Address)
                        m.Answer = append(m.Answer, inr)
                    }
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

    app.RegisterServiceWithFuncs(
        func() error {
            return errors.WithStack(udpServer.ActivateAndServe())
        },
        func(_ func(interface{})) error {
            udpServer.Shutdown()
            dns.HandleRemove(".")
            udpPacketConn = nil
            udpAddr = nil
            return nil
        },
    )

    return nil
}