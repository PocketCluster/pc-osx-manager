// +build darwin
package teleport

import (
    "testing"
    "fmt"

    . "gopkg.in/check.v1"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/stkim1/pc-core/context"
)

func TestConfig(t *testing.T) { TestingT(t) }

type ConfigSuite struct {
    dataDir     string
}

var _ = Suite(&ConfigSuite{})

func (s *ConfigSuite) SetUpSuite(c *C) {
    utils.InitLoggerForTests()
}

func (s *ConfigSuite) TearDownSuite(c *C) {
}

func (s *ConfigSuite) SetUpTest(c *C) {
    context.DebugContextPrepare()

    dataDir, _ := context.SharedHostContext().ApplicationUserDataDirectory()
    s.dataDir = dataDir + "/teleport"
    c.Logf("[INFO] User DataDir %s", dataDir)
}

func (s *ConfigSuite) TearDownTest(c *C) {
    context.DebugContextDestroy()

    s.dataDir = ""
}

func (s *ConfigSuite) TestDefaultConfig(c *C) {
    config := MakePocketTeleportConfig()
    c.Assert(config, NotNil)

    // all 3 services should be enabled by default
    c.Assert(config.Auth.Enabled, Equals, true)
    c.Assert(config.SSH.Enabled, Equals, true)
    c.Assert(config.Proxy.Enabled, Equals, true)

    localAuthAddr := utils.NetAddr{AddrNetwork: "tcp", Addr: "0.0.0.0:3025"}
    localProxyAddr := utils.NetAddr{AddrNetwork: "tcp", Addr: "0.0.0.0:3023"}
    localSSHAddr := utils.NetAddr{AddrNetwork: "tcp", Addr: "0.0.0.0:3022"}

    // data dir, hostname and auth server
    c.Assert(config.DataDir, Equals, s.dataDir)
    if len(config.Hostname) < 2 {
        c.Error("default hostname wasn't properly set")
    }

    // auth section
    auth := config.Auth
    c.Assert(auth.SSHAddr, DeepEquals, localAuthAddr)
    c.Assert(auth.Limiter.MaxConnections, Equals, int64(LimiterMaxConnections))
    c.Assert(auth.Limiter.MaxNumberOfUsers, Equals, LimiterMaxConcurrentUsers)

    c.Assert(auth.KeysBackend.Type, Equals, "SQLite")
    c.Assert(auth.KeysBackend.Params, Equals, fmt.Sprintf(`{"path": "%s/keys.db"}`, s.dataDir))
    c.Assert(auth.EventsBackend.Type, Equals, "SQLite")
    c.Assert(auth.EventsBackend.Params, Equals, fmt.Sprintf(`{"path": "%s/events.db"}`, s.dataDir))
    c.Assert(auth.RecordsBackend.Type, Equals, "SQLite")
    c.Assert(auth.RecordsBackend.Params, Equals, fmt.Sprintf(`{"path": "%s/records.db"}`, s.dataDir))

    // SSH section
    ssh := config.SSH
    c.Assert(ssh.Addr, DeepEquals, localSSHAddr)
    c.Assert(ssh.Limiter.MaxConnections, Equals, int64(LimiterMaxConnections))
    c.Assert(ssh.Limiter.MaxNumberOfUsers, Equals, LimiterMaxConcurrentUsers)

    // proxy section
    proxy := config.Proxy
    c.Assert(proxy.AssetsDir, Equals, s.dataDir)
    c.Assert(proxy.SSHAddr, DeepEquals, localProxyAddr)
    c.Assert(proxy.Limiter.MaxConnections, Equals, int64(LimiterMaxConnections))
    c.Assert(proxy.Limiter.MaxNumberOfUsers, Equals, LimiterMaxConcurrentUsers)
}
