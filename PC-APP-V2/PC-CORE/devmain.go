package main

import (
    "crypto/tls"
    "github.com/docker/swarm/api"
)
func main() {
    var tlsConfig *tls.Config;

    server := api.NewServer([]string{"0.0.0.0:2375"},tlsConfig)
    server.ListenAndServe()
}