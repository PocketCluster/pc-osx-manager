package main

import (
    "fmt"
    "os/user"
    "io/ioutil"

    "github.com/stkim1/pcrypto"
)

func TestSingingCertificate(saveFile bool) {
    // pk & ca
    _, cakey, cacert, err := pcrypto.GenerateClusterCertificateAuthorityData("cluster-uuid-here", "KR")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    _, nodekey, _, err := pcrypto.GenerateStrongKeyPair()
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    signer, err := pcrypto.NewCertAuthoritySigner(cakey, cacert, "cluster-uuid", "KR")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    nodecert, err := signer.GenerateSignedCertificate("odroid", "192.168.1.152", nodekey)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    _, masterkey, _, err := pcrypto.GenerateStrongKeyPair()
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    mastercert, err := signer.GenerateSignedCertificate("master", "", masterkey)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    cuser, err := user.Current()
    if err == nil && saveFile {
        ioutil.WriteFile(cuser.HomeDir + "/temp/signtest/gen/ca-key.pem",       cakey,      0600)
        ioutil.WriteFile(cuser.HomeDir + "/temp/signtest/gen/ca-cert.pub",      cacert,     0600)
        ioutil.WriteFile(cuser.HomeDir + "/temp/signtest/gen/node-key.pem",     nodekey,    0600)
        ioutil.WriteFile(cuser.HomeDir + "/temp/signtest/gen/node.cert",        nodecert,   0600)
        ioutil.WriteFile(cuser.HomeDir + "/temp/signtest/gen/master-key.pem",   masterkey,  0600)
        ioutil.WriteFile(cuser.HomeDir + "/temp/signtest/gen/master.cert",      mastercert, 0600)
    }
}

func main() {
    TestSingingCertificate(true)
}
