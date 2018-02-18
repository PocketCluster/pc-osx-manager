package config

import (
    "fmt"
    "strings"
    "path/filepath"
    "io/ioutil"

    "github.com/pkg/errors"
)

const (
    hostaddr_file          string = "/etc/hosts"
)

func (c *PocketSlaveConfig) UpdateHostsFile() error {
    var (
        hostsFile = filepath.Join(c.rootPath, hostaddr_file)
        rawContent, err = ioutil.ReadFile(hostsFile)
    )
    if err != nil {
        return errors.WithStack(err)
    }
    var (
        hostsData []string = nil
    )
    // filter out old host entry
    for _, hl := range strings.Split(string(rawContent),"\n") {
        if !strings.HasPrefix(hl, "127.0.1.1 ") {
            hostsData = append(hostsData, hl)
        }
    }
    hostsData = append(hostsData, fmt.Sprintf("127.0.1.1    %s %s", c.SlaveSection.SlaveNodeName, c.SlaveSection.SlaveNodeName))
    return ioutil.WriteFile(hostsFile, []byte(strings.Join(hostsData, "\n")), 0644)
}