package context

import (
    "github.com/pkg/errors"
)

type HostContextUserEnv interface {
    CocoaHomeDirectory() (string, error)
    PosixHomeDirectory() (string, error)
    FullUserName() (string, error)
    LoginUserName() (string, error)
    UserTemporaryDirectory() (string, error)
}

type hostUserEnv struct {
    cocoaHomePath                string
    posixHomePath                string
    fullUserName                 string
    loginUserName                string
    userTempPath                 string
}

func (ctx *hostUserEnv) CocoaHomeDirectory() (string, error) {
    if len(ctx.cocoaHomePath) == 0 {
        return "", errors.Errorf("[ERR] invalid cocoa home directory")
    }
    return ctx.cocoaHomePath, nil
}

func (ctx *hostUserEnv) PosixHomeDirectory() (string, error) {
    if len(ctx.posixHomePath) == 0 {
        return "", errors.Errorf("[ERR] invalid posix home directory")
    }
    return ctx.posixHomePath, nil
}

func (ctx *hostUserEnv) FullUserName() (string, error) {
    if len(ctx.fullUserName) == 0 {
        return "", errors.Errorf("[ERR] invalid full username")
    }
    return ctx.fullUserName, nil
}

func (ctx *hostUserEnv) LoginUserName() (string, error) {
    if len(ctx.loginUserName) == 0 {
        return "", errors.Errorf("[ERR] invalid login user name")
    }
    return ctx.loginUserName, nil
}

func (ctx *hostUserEnv) UserTemporaryDirectory() (string, error) {
    if len(ctx.userTempPath) == 0 {
        return "", errors.Errorf("[ERR] invalid user temp directory")
    }
    return ctx.userTempPath, nil
}
