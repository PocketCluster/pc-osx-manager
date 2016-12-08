package pcconfig

import (
    "github.com/gravitational/teleport/lib/service"
)

type NodeSetup struct {
    // network ip address of current host
    IP4Addr    string
    // key & cert save directory. This is where privatekey and cert will be saved
    KeyCertDir string
    // docker ca pub path
    DockerAuthFile string
    // docker Key file path
    DockerKeyFile string
    // docker cert file path
    DockerCertFile string
}

// Config structure is used to initialize _all_ services PocketCluster & Teleporot can run.
// Some settings are globl (like DataDir) while others are grouped into sections, like AuthConfig
type Config struct {
    // original key and cert
    service.Config
    // Slave node config
    NodeSetup
}
