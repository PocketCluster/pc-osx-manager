package main

import (
    "os"
    "os/user"
    "path/filepath"
    "strings"

    "github.com/pkg/errors"
    process "github.com/mitchellh/go-ps"
    "github.com/gravitational/teleport/lib/sshutils/scp"
    "github.com/gravitational/teleport/lib/utils"
)

type StdReadWriter struct {
}

func (rw *StdReadWriter) Read(b []byte) (int, error) {
    return os.Stdin.Read(b)
}

func (rw *StdReadWriter) Write(b []byte) (int, error) {
    return os.Stdout.Write(b)
}

func runScpCommand(cmd *scp.Command) (err error) {
    // pocketd
    sps, err := process.FindProcess(os.Getpid())
    if err != nil {
        return errors.WithStack(err)
    }
    // parent pocketd pid
    rps, err := process.FindProcess(sps.PPid())
    if err != nil {
        return errors.WithStack(err)
    }
    if rps.Executable() != pocketdExecName {
        return errors.Errorf("incorrect parent executable")
    }

    // get user's home dir (it serves as a default destination)
    cmd.User, err = user.Current()
    if err != nil {
        return errors.WithStack(err)
    }
    // see if the target is absolute. if not, use user's homedir to make
    // it absolute (and if the user doesn't have a homedir, use "/")
    slash := string(filepath.Separator)
    withSlash := strings.HasSuffix(cmd.Target, slash)
    if !filepath.IsAbs(cmd.Target) {
        rootDir := cmd.User.HomeDir
        if !utils.IsDir(rootDir) {
            cmd.Target = slash + cmd.Target
        } else {
            cmd.Target = filepath.Join(rootDir, cmd.Target)
            if withSlash {
                cmd.Target = cmd.Target + slash
            }
        }
    }
    if !cmd.Source && !cmd.Sink {
        return errors.Errorf("remote mode is not supported")
    }
    return errors.WithStack(cmd.Execute(&StdReadWriter{}))
}
