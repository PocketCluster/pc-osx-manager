package install

import (
    "encoding/json"
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/pc-core/route"
)

func makeMessageFeedBack(feeder route.ResponseFeeder, rpPath, message string) error {
    data, err := json.Marshal(route.ReponseMessage{
        "package-progress": {
            "message":      message,
        },
    })
    if err != nil {
        log.Errorf(err.Error())
        return errors.WithStack(err)
    }
    err = feeder.FeedResponseForPost(rpPath, string(data))
    if err != nil {
        log.Errorf(err.Error())
        return errors.WithStack(err)
    }
    return nil
}
