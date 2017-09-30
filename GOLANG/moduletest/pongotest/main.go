package main

import (
    "io/ioutil"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/flosch/pongo2"
)

func main() {
    templ, err := ioutil.ReadFile("test.taml")
    if err != nil {
        log.Error(errors.WithStack(err).Error())
    }
    tpl, err := pongo2.FromString(string(templ))
    if err != nil {
        log.Error(errors.WithStack(err).Error())
    }

    data := pongo2.Context{}

    data.Update(pongo2.Context{"slavenodes": []pongo2.Context{
        {"index":"1", "address":"2", "name":"pc-node1"},
        {"index":"2", "address":"3", "name":"pc-node2"},
        {"index":"3", "address":"5", "name":"pc-node4"},
        {"index":"4", "address":"6", "name":"pc-node5"},
    }})

    out, err := tpl.Execute(data)
    if err != nil {
        if err != nil {
            log.Error(err.Error())
        }
    }
    log.Info(out)
}
