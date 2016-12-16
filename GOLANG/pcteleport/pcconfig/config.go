package pcconfig

import (
    "github.com/gravitational/teleport/lib/service"
    "github.com/stkim1/pcrypto"
)

type NodeProperty struct {
    // network ip address of current host
    IP4Addr    string
    // docker ca pub path
    DockerAuthFile string
    // docker Key file path
    DockerKeyFile string
    // docker cert file path
    DockerCertFile string
}

type CoreProperty struct {
    *pcrypto.CaSigner
}

// Config structure is used to initialize _all_ services PocketCluster & Teleporot can run.
// Some settings are globl (like DataDir) while others are grouped into sections, like AuthConfig
type Config struct {
    // original key and cert
    service.Config
    // Slave node config
    NodeProperty
    // Teleport core config
    CoreProperty
}
