package regnode

import (
    "encoding/json"
    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-core/route"
)

func feedPostError(feeder route.ResponseFeeder, rpath, fpath string, irr error) error {
    log.Error(irr.Error())
    data, frr := json.Marshal(route.ReponseMessage{
        fpath: {
            "status": false,
            "error" : irr.Error(),
        },
    })
    // this should never happen
    if frr != nil {
        log.Error(frr.Error())
    }
    frr = feeder.FeedResponseForPost(rpath, string(data))
    if frr != nil {
        log.Error(frr.Error())
    }
    return irr
}

func feedGetError(feeder route.ResponseFeeder, rpath, fpath string, irr error) error {
    log.Error(irr.Error())
    data, frr := json.Marshal(route.ReponseMessage{
        fpath: {
            "status": false,
            "error" : irr.Error(),
        },
    })
    // this should never happen
    if frr != nil {
        log.Error(frr.Error())
    }
    frr = feeder.FeedResponseForGet(rpath, string(data))
    if frr != nil {
        log.Error(frr.Error())
    }
    return irr
}

func feedGetMessage(feeder route.ResponseFeeder, rpath, fpath, ppath string, fdata interface{}) error {
    data, frr := json.Marshal(route.ReponseMessage{
        fpath: {
            "status": true,
            ppath: fdata,
        },
    })
    // this should never happen
    if frr != nil {
        log.Error(frr.Error())
    }
    return feeder.FeedResponseForGet(rpath, string(data))
}

func feedGetOkMessage(feeder route.ResponseFeeder, rpath, fpath string) error {
    data, frr := json.Marshal(route.ReponseMessage{
        fpath: {
            "status": true,
        },
    })
    // this should never happen
    if frr != nil {
        log.Error(frr.Error())
    }
    return feeder.FeedResponseForGet(rpath, string(data))
}


func feedPostMessage(feeder route.ResponseFeeder, rpath, fpath, ppath string, fdata interface{}) error {
    data, frr := json.Marshal(route.ReponseMessage{
        fpath: {
            "status": true,
            ppath: fdata,
        },
    })
    // this should never happen
    if frr != nil {
        log.Error(frr.Error())
    }
    return feeder.FeedResponseForPost(rpath, string(data))
}