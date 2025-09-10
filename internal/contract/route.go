package contract

type (
	RouteRegistrar interface {
		Create(string, string) error
		Refresh() map[string]RouteRoles
		FindRoles(string) []string
	}

	RouteRoles struct {
		TransformedPath string
		Roles           []string
	}
)
