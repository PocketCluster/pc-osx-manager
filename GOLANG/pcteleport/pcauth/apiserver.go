package pcauth

import (
    "net/http"
    "github.com/gravitational/teleport"
    "github.com/gravitational/teleport/lib/httplib"
    "github.com/gravitational/teleport/lib/auth"

    "github.com/gravitational/trace"
    "github.com/julienschmidt/httprouter"
)

// APIServer implements http API server for AuthServer interface
type PocketAPIServer struct {
    httprouter.Router
    ar authWithRoles
}

// NewAPIServer returns a new instance of APIServer HTTP handler
func NewPocketAPIServer(config *auth.APIConfig, role teleport.Role, notFound http.HandlerFunc) PocketAPIServer {
    srv := PocketAPIServer{
        ar: authWithRoles{
            authServer:     config.AuthServer,
            permChecker:    config.PermissionChecker,
            sessions:       config.SessionService,
            role:           role,
            alog:           config.AuditLog,
        },
    }
    srv.Router   = *httprouter.New()
    srv.NotFound = notFound

    srv.POST("/v1/cert/issuesigned", httplib.MakeHandler(srv.issueSignedCertificatewithToken))

    return srv
}

const (
    pcHostToken string      = "token"
    pcHostID string         = "hostid"
    pcHostName string       = "hostname"
    pcHostIp4Addr string    = "ip4addr"
    pcHostRole string       = "role"
)

type signedCertificateReq struct {
    Token    string        `json:"token"`
    HostID   string        `json:"hostid"`
    Hostname string        `json:"hostname"`
    IP4Addr  string        `json:"ip4addr"`
    Role     teleport.Role `json:"role"`
}

func (s *PocketAPIServer) issueSignedCertificatewithToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (interface{}, error) {
    var req *signedCertificateReq
    if err := httplib.ReadJSON(r, &req); err != nil {
        return nil, trace.Wrap(err)
    }
    keys, err := s.ar.issueSignedCertificateWithToken(req)
    if err != nil {
        return nil, trace.Wrap(err)
    }
    return keys, nil
}
