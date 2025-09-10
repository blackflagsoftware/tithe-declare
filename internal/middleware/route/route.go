package route

import (
	"fmt"
	"regexp"
	"strings"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	c "github.com/blackflagsoftware/tithe-declare/internal/contract"
	l "github.com/blackflagsoftware/tithe-declare/internal/middleware/logging"

	"github.com/labstack/echo/v4"
)

type (
	Route struct {
		RouteRegister c.RouteRegistrar
		Map           map[string]c.RouteRoles
	}
)

var internalRoute Route

func InitRoute(rt c.RouteRegistrar) {
	routeMap := make(map[string]c.RouteRoles)
	internalRoute = Route{RouteRegister: rt, Map: routeMap}
	internalRoute.RefreshRouteRoles()
}

func RegisterAndAdd(group *echo.Group, method, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route {
	internalRoute.Register(method, path)
	return group.Add(method, path, handler, middleware...)
}

func GetRolesForRegisterRoute(method, subpath string) ([]string, error) {
	path := normalizePath(method, subpath)
	for _, route := range internalRoute.Map {
		reg, err := regexp.Compile(route.TransformedPath)
		if err != nil {
			return []string{}, ae.GeneralError("General error in checking registered routes", err)
		}
		if reg.Match([]byte(path)) {
			return internalRoute.RouteRegister.FindRoles(path), nil
		}
	}
	return []string{}, fmt.Errorf("missing route")
}

func (r *Route) RefreshRouteRoles() {
	r.Map = r.RouteRegister.Refresh()
}

func (r *Route) Register(method, subpath string) {
	path := normalizePath(method, subpath)
	// is the path already exists in memory for rawpath
	// if not add it
	if _, ok := r.Map[path]; !ok {
		// does not exist, add
		transformPath := transformPathToRegex(path)
		if err := r.RouteRegister.Create(path, transformPath); err != nil {
			l.Default.Printf("Register - unable to create new registered path: %s", path)
		}
		r.RefreshRouteRoles()
	}
}

func transformPathToRegex(path string) string {
	pathLen := len(path) // used to know when to end the search
	for {
		idx := strings.Index(path, ":")
		if idx == -1 {
			break
		}
		slashIdx := pathLen // default it the end of the string
		for i := idx; i < pathLen; i++ {
			if path[i] == '/' {
				slashIdx = i
				break
			}
		}
		path = path[:idx] + ".+" + path[slashIdx:]
		pathLen = len(path) // reset length
	}
	return path
}

func normalizePath(method, path string) string {
	if len(path) == 0 {
		return method + "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	return method + path
}
