package registerroute

import (
	"encoding/json"

	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

type (
	RegisterRoute struct {
		RawPath         string           `db:"raw_path" json:"raw_path"`
		TransformedPath null.String      `db:"transformed_path" json:"transformed_path"`
		Roles           *json.RawMessage `db:"roles" json:"roles"`
	}

	BulkRegisterRoute struct {
		RawPaths    []string `db:"raw_paths" json:"raw_paths"`
		AddRoles    []string `db:"add_roles" json:"add_roles"`
		RemoveRoles []string `db:"remove_roles" json:"remove_roles"`
	}

	RegisterRouteParam struct {
		// TODO: add any other custom params here
		h.Param
	}
)

const RegisterRouteConst = "register_route"

func InitStorageV1() DataRegisterRouteV1Adapter {
	return InitSQLV1()
}
