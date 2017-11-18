package dockertool

import (
    "github.com/pkg/errors"
    "golang.org/x/net/context"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"

    log "github.com/Sirupsen/logrus"
)

func CleanupContainer(cli *client.Client) error {
    containers, err := cli.ContainerList(context.TODO(), types.ContainerListOptions{All:true})
    if err != nil {
        return errors.WithStack(err)
    }

    for _, c := range containers {
        /* don't need to kill container
        if c.Status != "exited" {
            if err := cli.ContainerKill(context.TODO(), c.ID, "SIGINT"); err != nil {
                log.Error(err.Error())
            }
        }*/
        if err := cli.ContainerRemove(context.TODO(), c.ID, types.ContainerRemoveOptions{Force:true}); err != nil {
            log.Errorf("container cleanup error %v", err.Error())
        }
    }

    return nil
}

func CleanupNetwork(cli *client.Client) error {
/*
    // need api 1.25+
    if rpt, err := cli.NetworksPrune(context.TODO(), types.NetworksPruneConfig{}); err != nil {
        log.Error(err.Error())
    } else {
        log.Infof("network prune report %v", rpt)
    }
*/
    networks, err := cli.NetworkList(context.TODO(), types.NetworkListOptions{})
    if err != nil {
        return errors.WithStack(err)
    }
    for _, n := range networks {
        if n.Scope == "global" {
            if err := cli.NetworkRemove(context.TODO(), n.ID); err != nil {
                log.Errorf("network cleanup error %v", err.Error())
            }
        }
    }
    return nil
}