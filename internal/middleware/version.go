package middleware

import (
	"slices"
	"strings"

	"github.com/blackflagsoftware/tithe-declare/config"
	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	"github.com/labstack/echo/v4"
)

func VersionHandler(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		accept := c.Request().Header.Get("Accept")
		if accept == "" {
			return ae.MissingAcceptHeaderError(config.Srv.AppName)
		}
		versionPlus := strings.Replace(accept, "application/vnd."+config.Srv.AppName+".", "", 1) // remove the content type prefix with application name, should get back v1+json
		versionSplit := strings.Split(versionPlus, "+")
		if len(versionPlus) < 2 {
			// not in the expected format v1+json
			return ae.MissingAcceptHeaderError(config.Srv.AppName)
		}
		knownVersions := config.KnownVersions()
		if !slices.Contains(knownVersions, versionSplit[0]) {
			return ae.MissingAcceptHeaderError(config.Srv.AppName)
		}
		c.Set("version", versionSplit[0])
		if err := h(c); err != nil {
			c.Error(err)
		}
		return
	}
}
