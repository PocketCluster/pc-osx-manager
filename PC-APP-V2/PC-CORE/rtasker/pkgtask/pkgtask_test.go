package pkgtask

import (
    "bytes"
    "io/ioutil"
    "path/filepath"
    "runtime"
    "testing"

    log "github.com/Sirupsen/logrus"
)

func TestTemplateGeneration(t *testing.T) {
    var (
        nodeList = []map[string]string{
            {"nodename":"pc-node1"},
            {"nodename":"pc-node3"},
            {"nodename":"pc-core"},
            {"nodename":"pc-node4"},
            {"nodename":"pc-node6"}}

        etmpl = []byte(`container_name: pc-core
ipv4_address: 172.16.128.1

datanode1:
container_name: pc-node1
ipv4_address: 172.16.128.2

datanode3:
container_name: pc-node3
ipv4_address: 172.16.128.4

datanode4:
container_name: pc-node4
ipv4_address: 172.16.128.5

datanode6:
container_name: pc-node6
ipv4_address: 172.16.128.7
`)
    )
    log.SetLevel(log.DebugLevel)

    // load template file
    _, testfile, _, _ := runtime.Caller(0)
    data, err := ioutil.ReadFile(filepath.Join(filepath.Dir(testfile), "test.taml"))
    if err != nil {
        t.Fatal(err.Error())
    }

    // generate template
    gtmpl, err := buildComposeTemplateWithNodeList(data, nodeList)
    if err != nil {
        t.Fatal(err.Error())
    }
    log.Info(string(gtmpl))

    if bytes.Compare(etmpl, gtmpl) != 0 {
        t.Fatal("invalid template generation")
    }
}