package dockertool

import (
    "net/http"

    "github.com/pkg/errors"
    "golang.org/x/net/context"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"

    pcctx "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/utils/tlscfg"
)

func NewContainerClient(host, version string) (*client.Client, error) {
    caCert, err := pcctx.SharedHostContext().CertAuthCertificate()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    hostCrt, err := pcctx.SharedHostContext().MasterHostCertificate()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    hostKey, err := pcctx.SharedHostContext().MasterHostPrivateKey()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    tlsc, err := tlscfg.BuildTLSConfigWithCAcert(caCert, hostCrt, hostKey, true)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    httpcli := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: tlsc,
        },
    }
    if len(version) == 0 {
        version = client.DefaultVersion
    }

    // empty version enables client to automatically override version
    cli, err := client.NewClient(host, version, httpcli, nil)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return cli, nil
}

func InstallImageFromRepository(cli *client.Client, imageRef string) error {
    if len(imageRef) == 0 {
        return errors.Errorf("cannot install invalid image path")
    }
    _, err := cli.ImagePull(context.TODO(), imageRef, types.ImagePullOptions{})
    return errors.WithStack(err)
}