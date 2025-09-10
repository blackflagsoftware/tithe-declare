package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	l "github.com/blackflagsoftware/tithe-declare/internal/middleware/logging"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type (
	Output struct {
		Payload any `json:"data,omitempty"`
		*Error  `json:"error,omitempty"`
		*Meta   `json:"meta,omitempty"`
	}

	Error struct {
		Id     string `json:"Id,omitempty"`
		Title  string `json:"Title,omitempty"`
		Detail string `json:"Detail,omitempty"`
		Status string `json:"Status,omitempty"`
	}

	Meta struct {
		TotalCount int `json:"total_count"`
	}

	Param struct {
		Search           Search `json:"search"`
		Limit            int    // holds the calculated limit
		Offset           int    // holds the offset number
		PaginationString string // holds the limit/offset tring
		Sort             string // holds the calculated sort string
		ColumnMapping    map[string]string
	}

	Search struct {
		Filters    []Filter   `json:"filters"`
		Pagination Pagination `json:"pagination"`
		Sort       string     `json:"sort"` // comma separated string, use a '-' before column name to sort DESC i.e.: id,-name => "SORT BY id ASC, name DESC"
	}

	Filter struct {
		Column  string `json:"column"`
		Compare string `json:"compare"`
		Value   any    `json:"value"`
	}

	Pagination struct {
		PageLimit  int `json:"page_limit"`
		PageNumber int `json:"page_number"`
	}
)

func FormatResponse(c echo.Context, statusCode int, payload any, totalCount *int) error {
	var meta *Meta
	if totalCount != nil {
		meta = &Meta{TotalCount: *totalCount}
	}
	output := Output{
		Payload: payload,
		Error:   nil,
		Meta:    meta,
	}
	return c.JSON(statusCode, output)
}

func FormatResponseWithError(c echo.Context, apiError ae.ApiError) error {
	LogError(c, &apiError)
	err := &Error{Id: apiError.ApiErrorCode, Title: apiError.Title, Detail: apiError.Detail, Status: strconv.Itoa(apiError.StatusCode)}
	output := Output{
		Payload: nil,
		Error:   err,
		Meta:    nil,
	}
	return c.JSON(apiError.StatusCode, output)
}

func LogError(c echo.Context, apiError *ae.ApiError) {
	url := c.Request().URL.String()

	l.Default.WithFields(
		logrus.Fields{
			"method":      c.Request().Method,
			"status_code": apiError.StatusCode,
			"status_text": http.StatusText(c.Response().Status),
			"request_url": url,
			"referer":     c.Request().Referer(),
			"user_agent":  c.Request().UserAgent(),
			"remote":      c.Request().RemoteAddr,
			"detail":      apiError.Detail,
		},
	).Errorln("error")
}

func (p *Param) CalculateParam(primarySort string, availableSort map[string]string) (err error) {
	p.ColumnMapping = availableSort
	// calculate the limit
	if p.Search.Pagination.PageLimit > 0 {
		if p.Search.Pagination.PageNumber == 0 {
			// should not be empty, default to first page
			p.Search.Pagination.PageNumber = 1
		}
		p.Limit = p.Search.Pagination.PageLimit
		p.Offset = p.Search.Pagination.PageNumber - 1
		p.Offset *= p.Search.Pagination.PageLimit
	}
	// calculate the sort
	if primarySort == "" {
		return
	}
	if p.Search.Sort == "" {
		p.Search.Sort = primarySort
	}
	sorted := []string{}
	sortParts := strings.Split(p.Search.Sort, ",")
	for _, s := range sortParts {
		direction := "ASC"
		name := s
		if string(name[0]) == "-" {
			direction = "DESC"
			name = string(name[1:])
		}
		if _, ok := availableSort[name]; !ok {
			// if the name is not in the available sort list, you could return and error here
			continue
		}
		sorted = append(sorted, fmt.Sprintf("%s %s", availableSort[name], direction))
	}
	p.Sort = strings.Join(sorted, ", ")
	return
}
