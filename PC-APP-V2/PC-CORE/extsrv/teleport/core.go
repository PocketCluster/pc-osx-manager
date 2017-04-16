package teleport

import (
    "time"

    "github.com/gravitational/teleport/lib/config"
    "github.com/gravitational/teleport/lib/process"
    "github.com/gravitational/teleport/lib/service"
    "github.com/gravitational/teleport/lib/services"

    "github.com/pkg/errors"
)

func NewTeleportCore(cfg *service.PocketConfig) (*process.PocketCoreProcess, error) {
    // add static tokens
    for _, token := range []config.StaticToken{"node:d52527f9-b260-41d0-bb5a-e23b0cfe0f8f", "node:c9s93fd9-3333-91d3-9999-c9s93fd98f43"} {
        roles, tokenValue, err := token.Parse()
        if err != nil {
            return nil, errors.WithStack(err)
        }
        cfg.Auth.StaticTokens = append(cfg.Auth.StaticTokens, services.ProvisionToken{Token: tokenValue, Roles: roles, Expires: time.Unix(0, 0)})
    }
    // new process
    return process.NewCoreProcess(cfg)
}
