package dockertool

import (
    "time"

    "github.com/pkg/errors"
    "golang.org/x/net/context"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"

    log "github.com/Sirupsen/logrus"
    "github.com/docker/docker/api/types/filters"
)

const (
    // 30 second waiting context
    deadline = time.Second * time.Duration(3)
    timeout  = time.Second * time.Duration(5)
)

// all container cleanup
func CleanupContainer(cli *client.Client) error {
    ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(deadline))
    defer cancel()
    containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All:true})
    if err != nil {
        return errors.WithStack(err)
    }
    for _, c := range containers {
        if err := cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force:true}); err != nil {
            log.Errorf("[CONTAINER] container cleanup error %v", err.Error())
        }
    }
    select {
        case <- time.After(timeout): {
            return errors.Errorf("[CONTAINER] container cleanup timeout")
        }
        case <- ctx.Done(): {
            log.Info("[CONTAINER] container cleanup ok")
            return nil
        }
    }
}

// global network cleanup
func CleanupNetwork(cli *client.Client) error {
/*
    // need api 1.25+
    if rpt, err := cli.NetworksPrune(context.TODO(), types.NetworksPruneConfig{}); err != nil {
        log.Error(err.Error())
    } else {
        log.Infof("network prune report %v", rpt)
    }
*/
    ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(deadline))
    defer cancel()
    networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
    if err != nil {
        return errors.WithStack(err)
    }
    for _, n := range networks {
        if n.Scope == "global" {
            if err := cli.NetworkRemove(ctx, n.ID); err != nil {
                log.Errorf("[CONTAINER] network cleanup error %v", err.Error())
            }
        }
    }
    select {
        case <- time.After(timeout): {
            return errors.Errorf("[CONTAINER] network cleanup timeout")
        }
        case <- ctx.Done(): {
            log.Info("[CONTAINER] network cleanup ok")
            return nil
        }
    }
}

func CleanupVolume(cli *client.Client) error {
    ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(deadline))
    defer cancel()
    args, err := filters.ParseFlag("dangling=true" , filters.NewArgs())
    if err != nil {
        return errors.WithStack(err)
    }
    vlist, err := cli.VolumeList(ctx, args)
    if err != nil {
        return errors.WithStack(err)
    }
    for i, _ := range vlist.Volumes {
        v := vlist.Volumes[i]
        cli.VolumeRemove(ctx, v.Name, true)
    }
    select {
        case <- time.After(timeout): {
            return errors.Errorf("[CONTAINER] volume pruning timeout")
        }
        case <- ctx.Done(): {
            log.Info("[CONTAINER] volume cleanup ok")
            return nil
        }
    }
}