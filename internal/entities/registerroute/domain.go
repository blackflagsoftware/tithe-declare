package registerroute

import (
	"encoding/json"
	"slices"

	"context"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	a "github.com/blackflagsoftware/tithe-declare/internal/audit"
	c "github.com/blackflagsoftware/tithe-declare/internal/contract"
	l "github.com/blackflagsoftware/tithe-declare/internal/middleware/logging"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	"github.com/blackflagsoftware/tithe-declare/internal/util/function"
	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

//go:generate mockgen -source=domain.go -destination=mock.go -package=registerroute
type (
	DataRegisterRouteV1Adapter interface {
		Read(context.Context, *RegisterRoute) error
		ReadAll(context.Context, *[]RegisterRoute, RegisterRouteParam) (int, error)
		Create(context.Context, *RegisterRoute) error
		Update(context.Context, RegisterRoute) error
		Delete(context.Context, *RegisterRoute) error
	}

	DomainRegisterRouteV1 struct {
		dataRegisterRouteV1 DataRegisterRouteV1Adapter
		auditWriter         a.AuditAdapter
	}
)

// these are the routes that will not, by default, set 'admin' for the roles to be checked against in the middlewar
// they are in the format of "HTTP_METHOD/path/to/route" with the path being the raw path found in the "rest.go" file
var InsecureRoutes = []string{
	"POST/login/reset/pwd",
	"GET/login/forgot-password/:email_addr",
	"POST/login/sign-in",
	"POST/login/oauth2/authorize",
	"POST/login/oauth2/verify-consent",
}

func NewDomainRegisterRouteV1(cregV1 DataRegisterRouteV1Adapter) *DomainRegisterRouteV1 {
	aw := a.AuditInit()
	return &DomainRegisterRouteV1{dataRegisterRouteV1: cregV1, auditWriter: aw}
}

func (m *DomainRegisterRouteV1) Get(ctx context.Context, reg *RegisterRoute) error {
	if reg.RawPath == "" {
		return ae.MissingParamError("RawPath")
	}
	return m.dataRegisterRouteV1.Read(ctx, reg)
}

func (m *DomainRegisterRouteV1) Search(ctx context.Context, reg *[]RegisterRoute, param RegisterRouteParam) (int, error) {
	// the second argument (map[string]string) is a list of columns to use for filtering
	// the key matches the json struct tag, the value is the actual table column name (this should change if aliases are used in your query)
	param.Param.CalculateParam("transformed_path", map[string]string{"raw_path": "raw_path", "transformed_path": "transformed_path", "roles": "roles"})
	param.Param.PaginationString = stor.FormatPagination(param.Param.Limit, param.Param.Offset)

	return m.dataRegisterRouteV1.ReadAll(ctx, reg, param)
}

func (m *DomainRegisterRouteV1) Post(ctx context.Context, reg *RegisterRoute) error {
	if !reg.TransformedPath.Valid {
		return ae.MissingParamError("TransformedPath")
	}
	if reg.TransformedPath.Valid && len(reg.TransformedPath.ValueOrZero()) > 255 {
		return ae.StringLengthError("TransformedPath", 255)
	}
	if err := m.dataRegisterRouteV1.Create(ctx, reg); err != nil {
		return err
	}
	go a.AuditCreate(m.auditWriter, *reg, RegisterRouteConst, a.KeysToString("raw_path", reg.RawPath))
	return nil
}

func (m *DomainRegisterRouteV1) Patch(ctx context.Context, regIn RegisterRoute) error {
	reg := &RegisterRoute{RawPath: regIn.RawPath}
	errGet := m.dataRegisterRouteV1.Read(ctx, reg)
	if errGet != nil {
		return errGet
	}
	existingValues := make(map[string]any)
	// TransformedPath
	if regIn.TransformedPath.Valid {
		if regIn.TransformedPath.Valid && len(regIn.TransformedPath.ValueOrZero()) > 255 {
			return ae.StringLengthError("TransformedPath", 255)
		}
		existingValues["transformed_path"] = reg.TransformedPath.String
		reg.TransformedPath = regIn.TransformedPath
	}

	// Roles
	if regIn.Roles != nil {
		if !function.ValidJson(*regIn.Roles) {
			return ae.ParseError("Invalid JSON syntax for Roles")
		}
		existingValues["roles"] = reg.Roles
		reg.Roles = regIn.Roles
	}
	if err := m.dataRegisterRouteV1.Update(ctx, *reg); err != nil {
		return err
	}
	go a.AuditPatch(m.auditWriter, *reg, RegisterRouteConst, a.KeysToString("raw_path", reg.RawPath), existingValues)
	return nil
}

func (m *DomainRegisterRouteV1) Delete(ctx context.Context, reg *RegisterRoute) error {
	if reg.RawPath == "" {
		return ae.MissingParamError("RawPath")
	}
	if err := m.dataRegisterRouteV1.Delete(ctx, reg); err != nil {
		return err
	}
	go a.AuditDelete(m.auditWriter, *reg, RegisterRouteConst, a.KeysToString("raw_path", reg.RawPath))
	return nil
}

func (m *DomainRegisterRouteV1) Create(path, transformPath string) error {
	ctx := context.Background()
	byteRoles := json.RawMessage([]byte("[\"admin\"]"))
	if slices.Contains(InsecureRoutes, path) {
		byteRoles = json.RawMessage([]byte("[]"))
	}
	regRoute := RegisterRoute{RawPath: path, TransformedPath: null.StringFrom(transformPath), Roles: &byteRoles}
	return m.Post(ctx, &regRoute)
}

func (m *DomainRegisterRouteV1) Refresh() map[string]c.RouteRoles {
	regRoutes := make(map[string]c.RouteRoles)
	ctx := context.Background()
	registerRoutes := []RegisterRoute{}
	_, err := m.Search(ctx, &registerRoutes, RegisterRouteParam{})
	if err != nil {
		l.Default.Printf("Unable to register routes: %s", err)
		return regRoutes
	}
	for i, rr := range slices.All(registerRoutes) {
		roles := []string{}
		byteRoles, err := registerRoutes[i].Roles.MarshalJSON()
		if err != nil {
			l.Default.Printf("unable to marshal rawmessage: %s", err)
			continue
		}
		err = json.Unmarshal(byteRoles, &roles)
		if err != nil {
			l.Default.Printf("unable to marshal roles: %s", err)
			continue
		}
		routeRole := c.RouteRoles{TransformedPath: registerRoutes[i].TransformedPath.String, Roles: roles}
		regRoutes[rr.RawPath] = routeRole
	}
	return regRoutes
}

func (m *DomainRegisterRouteV1) FindRoles(rawPath string) []string {
	roles := []string{}
	ctx := context.Background()
	registerRoutes := []RegisterRoute{}
	_, err := m.Search(ctx, &registerRoutes, RegisterRouteParam{Param: h.Param{Search: h.Search{Filters: []h.Filter{{Column: "raw_path", Value: rawPath, Compare: "="}}}}})
	if err != nil {
		l.Default.Printf("Unable to register routes: %s", err)
		return roles
	}
	if len(registerRoutes) > 0 {
		byteRoles, err := registerRoutes[0].Roles.MarshalJSON()
		if err != nil {
			l.Default.Printf("unable to marshal rawmessage: %s", err)
			return roles
		}
		err = json.Unmarshal(byteRoles, &roles)
		if err != nil {
			l.Default.Printf("unable to marshal roles: %s", err)
			return roles
		}
	}
	return roles
}

func (m *DomainRegisterRouteV1) Bulk(ctx context.Context, reg BulkRegisterRoute) error {
	if len(reg.RawPaths) == 0 {
		return ae.MissingParamError("RawPaths")
	}
	for _, rawPath := range reg.RawPaths {
		if rawPath == "" {
			continue
		}
		r := RegisterRoute{RawPath: rawPath}
		if err := m.Get(ctx, &r); err != nil {
			return err
		}
		roles := []string{}
		err := json.Unmarshal(*r.Roles, &roles)
		if err != nil {
			return ae.ParseError("Invalid JSON syntax for Roles")
		}
		if len(reg.AddRoles) > 0 {
			for _, addRole := range reg.AddRoles {
				if !slices.Contains(roles, addRole) {
					roles = append(roles, addRole)
				}
			}
		}
		if len(reg.RemoveRoles) > 0 {
			for _, removeRole := range reg.RemoveRoles {
				if slices.Contains(roles, removeRole) {
					roles = slices.Delete(roles, slices.Index(roles, removeRole), slices.Index(roles, removeRole)+1)
				}
			}
		}
		if len(roles) == 0 {
			r.Roles = nil
		} else {
			byteRoles, err := json.Marshal(roles)
			if err != nil {
				return ae.ParseError("Invalid JSON syntax for Roles")
			}
			jRoles := json.RawMessage(byteRoles)
			r.Roles = &jRoles
		}
		if err := m.dataRegisterRouteV1.Update(ctx, r); err != nil {
			return err
		}
		go a.AuditCreate(m.auditWriter, r, RegisterRouteConst, a.KeysToString("raw_path", r.RawPath))
	}
	return nil
}
