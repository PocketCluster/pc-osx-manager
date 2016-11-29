srv.POST("/v1/users", httplib.MakeHandler(srv.upsertUser))

// User represents teleport or external user
type User interface {
	// GetName returns user name
	GetName() string
	// GetAllowedLogins returns user's allowed linux logins
	GetAllowedLogins() []string
	// GetIdentities returns a list of connected OIDCIdentities
	GetIdentities() []OIDCIdentity
	// String returns user
	String() string
	// Check checks if all parameters are correct
	Check() error
	// Equals checks if user equals to another
	Equals(other User) bool
}

	// TeleportUser is an optional user entry in the database
	type TeleportUser struct {
		// Name is a user name
		Name string `json:"name"`

		// AllowedLogins represents a list of OS users this teleport
		// user is allowed to login as
		AllowedLogins []string `json:"allowed_logins"`

		// OIDCIdentities lists associated OpenID Connect identities
		// that let user log in using externally verified identity
		OIDCIdentities []OIDCIdentity `json:"oidc_identities"`
	}

func (s *APIServer) upsertUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
	var req *upsertUserReqRaw
	if err := httplib.ReadJSON(r, &req); err != nil {
		return nil, trace.Wrap(err)
	}
	user, err := services.GetUserUnmarshaler()(req.User)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	err = s.a.UpsertUser(user)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return message(fmt.Sprintf("'%v' user upserted", user.GetName())), nil
}

	func (a *AuthWithRoles) UpsertUser(u services.User) error {
		if err := a.permChecker.HasPermission(a.role, ActionUpsertUser); err != nil {
			return trace.Wrap(err)
		}
		return a.authServer.UpsertUser(u)
	}

		// UpsertUser updates parameters about user
		func (s *IdentityService) UpsertUser(user services.User) error {
			if !cstrings.IsValidUnixUser(user.GetName()) {
				return trace.BadParameter("'%v is not a valid unix username'", user.GetName())
			}

			for _, l := range user.GetAllowedLogins() {
				if !cstrings.IsValidUnixUser(l) {
					return trace.BadParameter("'%v is not a valid unix username'", l)
				}
			}
			for _, i := range user.GetIdentities() {
				if err := i.Check(); err != nil {
					return trace.Wrap(err)
				}
			}
			data, err := json.Marshal(user)
			if err != nil {
				return trace.Wrap(err)
			}

			err = s.backend.UpsertVal([]string{"web", "users", user.GetName()}, "params", []byte(data), backend.Forever)
			if err != nil {
				return trace.Wrap(err)
			}
			return nil
		}
