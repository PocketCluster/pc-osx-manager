package config

import (
    "fmt"
    "gopkg.in/yaml.v2"
)

func ExampleInitConfigBuild() {
    cfg := buildInitConfig()
    out, err := yaml.Marshal(cfg)
    if err != nil {
        fmt.Print(err.Error())
        return
    }
    fmt.Print(string(out))
    // Output:
    // config-version: 1.0.1
    // binding-status: unbounded
    // master-section:
    //   master-binder-agent: ""
    //   master-ip4-addr: ""
    //   master-timezone: ""
    // slave-section:
    //   slave-mac-addr: ""
    //   slave-node-name: ""
    //   slave-ip4-addr: ""
    //   slave-net-mask: ""
    //   slave-gateway: ""
}