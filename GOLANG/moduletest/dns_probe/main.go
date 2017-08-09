package main

import (
    "net"
    "strings"
    "sync"
    "fmt"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/miekg/dns"
)

/*
    // logging
    GODEBUG=netdns=1     go run ./main.go

    // verbose logging + cgo
    GODEBUG=netdns=cgo+2 go run ./main.go

    // go resolver + logging level
    GODEBUG=netdns=go    go run ./main.go
    GODEBUG=netdns=go+1  go run ./main.go
    GODEBUG=netdns=go+2  go run ./main.go

    // check server status with dig
    dig @127.0.0.1 -p 10059 pc-master
 */

const (
    pcmaster string = "pc-master"
    pcnode1  string = "pc-node1"
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
    remoteIP4 := remoteIP.To4()

    m := new(dns.Msg)
    m.Id = req.Id
    m.Question = req.Question
    m.Response = true

    switch qtype {
        case dns.TypeA: {

            if remoteIP4 != nil && (strings.HasPrefix(question.Name, pcmaster) || strings.HasPrefix(question.Name, pcnode1)) {
                rr := new(dns.A)
                rr.Hdr = dns.RR_Header{
                    Name:      question.Name,
                    Rrtype:    question.Qtype,
                    Class:     dns.ClassINET,
                    Ttl:       10,
                }
                rr.A = net.ParseIP("192.168.1.166")

                m.Answer = []dns.RR{rr}
                m.Authoritative = true
                w.WriteMsg(m)
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

func initDNSService(wg *sync.WaitGroup) error {
    dns.HandleFunc(".", locaNameServe)
    var (
        udpServer = &dns.Server {
            Addr:    "127.0.0.1:10059",
            Net:     "udp",
        }
        udpPacketConn *net.UDPConn = nil
        udpAddr *net.UDPAddr = nil
        err error = nil
    )
    defer func () {
        udpServer.Shutdown()
        dns.HandleRemove(".")
        udpPacketConn = nil
        udpAddr = nil
        wg.Done()
    }()
    udpAddr, err = net.ResolveUDPAddr(udpServer.Net, udpServer.Addr)
    if err != nil {
        log.Errorf("%v", err)
        return errors.WithStack(err)
    }
    udpPacketConn, err = net.ListenUDP(udpServer.Net, udpAddr)
    if err != nil {
        log.Errorf("%v", err)
        return errors.WithStack(err)
    }
    udpServer.PacketConn = udpPacketConn
    return errors.WithStack(udpServer.ActivateAndServe())
}

func check_host_file(host string) {
    log.Infof("Looking up %s", host)
    addrs, err := net.LookupHost(host)
    if err != nil {
        log.Errorf("%v", err)
    } else {
        log.Infof("%v", strings.Join(addrs, ", "))
    }
    fmt.Printf("\n")
}

func main() {
    var dnsWG sync.WaitGroup
    dnsWG.Add(1)
    go initDNSService(&dnsWG)

    check_host_file(pcmaster)
    check_host_file(pcnode1)
    check_host_file("www.google.com")

    dnsWG.Wait()
}
