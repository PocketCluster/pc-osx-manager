package install

import (
    "net/http"

    "github.com/pkg/errors"
    "golang.org/x/net/context"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"

    pccctx "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/utils/tlscfg"
)

func NewContainerClient(host string) (*client.Client, error) {
    ctx := pccctx.SharedHostContext()
    caCert, err := ctx.CertAuthCertificate()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    hostCrt, err := ctx.MasterHostCertificate()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    hostPrv, err := ctx.MasterHostPrivateKey()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    tlsc, err := tlscfg.BuildTLSConfigWithCAcert(caCert, hostCrt, hostPrv, true)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    httpcli := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: tlsc,
        },
    }

    // empty version enables client to automatically override version
    cli, err := client.NewClient(host, client.DefaultVersion, httpcli, nil)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return cli, nil
}

func InstallImageFromRepository(cli *client.Client, imageRef string) error {
    if len(imageRef) == 0 {
        return errors.Errorf("cannot install invalid image path")
    }
    r, err := cli.ImagePull(context.Background(), imageRef, types.ImagePullOptions{})
    defer r.Close()
    return errors.WithStack(err)
}