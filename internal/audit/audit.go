package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/blackflagsoftware/tithe-declare/config"
	"github.com/jmoiron/sqlx"
)

type (
	AuditAdapter interface {
		WriteAudit(Audit)
	}

	AuditFile struct {
		FilePath string
	}

	AuditSQL struct {
		DB        *sqlx.DB
		CreatedAt time.Time       `db:"created_at"`
		Entity    string          `db:"entity"`
		EntityID  string          `db:"entity_id"`
		Changed   json.RawMessage `db:"changed"`
		UserID    int             `db:"user_id"`
		UserUID   string          `db:"user_uid"`
	}

	Audit struct {
		CreatedAt time.Time              `json:"created_at"`
		Created   map[string]any         `json:"created,omitempty"`
		Updated   map[string]AuditUpdate `json:"updated,omitempty"`
		Delete    map[string]any         `json:"delete,omitempty"`
		UserID    int                    `json:"user_id,omitempty"`
		UserUID   string                 `json:"user_uid,omitempty"`
		Entity    string                 `json:"entity,omitempty"`
		EntityID  string                 `json:"entity_id,omitempty"`
	}

	AuditUpdate struct {
		From any `json:"from"`
		To   any `json:"to"`
	}

	AuditColumns struct {
		Created map[string]any         `json:"created,omitempty"`
		Updated map[string]AuditUpdate `json:"updated,omitempty"`
		Delete  map[string]any         `json:"delete,omitempty"`
	}
)

/* if you are using this, you will need to create a table
CREATE TABLE IF NOT EXISTS audit (
	id INT AUTO_INCREMENT, -- or SERIAL
	user_id int null,
	user_uid varchar(50) null,
	entity VARCHAR(50) NOT NULL,
	entity_id VARCHAR(50) NOT NULL,
	changed JSON NOT NULL, -- may not work with sqlite
	created_at DATETIME NOT NULL, -- or TIMESTAMP
	PRIMARY KEY(id)
);

*/
//go:generate mockgen -source=audit.go -destination=mock.go -package=audit
func AuditInit() AuditAdapter {
	// TODO: dynamically add this if using sql
	// if config.Aud.Storage == "sql" {
	// 	as := &AuditSQL{
	// 		DB: stor.InitStorage(),
	// 	}
	// 	return as
	// }
	return &AuditFile{FilePath: config.Aud.FilePath}
}

func AuditCreate(a AuditAdapter, entity any, entityName, entityId string) {
	if config.Aud.Enable {
		if a != nil {
			entityMap := GroupStructToMap(entity, "db")
			audit := Audit{Entity: entityName, EntityID: entityId, CreatedAt: time.Now().UTC(), Created: entityMap}
			a.WriteAudit(audit)
		}
	}
}

func AuditPatch(a AuditAdapter, entity any, entityName, entityId string, existingValues map[string]any) {
	if config.Aud.Enable {
		if a != nil {
			entityMap := GroupStructToMapUpdated(entity, "db", existingValues)
			audit := Audit{Entity: entityName, EntityID: entityId, CreatedAt: time.Now().UTC(), Updated: entityMap}
			a.WriteAudit(audit)
		}
	}
}

func AuditDelete(a AuditAdapter, entity any, entityName, entityId string) {
	if config.Aud.Enable {
		if a != nil {
			entityMap := GroupStructToMap(entity, "db")
			audit := Audit{Entity: entityName, EntityID: entityId, CreatedAt: time.Now().UTC(), Delete: entityMap}
			a.WriteAudit(audit)
		}
	}
}

func (h AuditFile) WriteAudit(audit Audit) {
	bAudit, err := json.Marshal(audit)
	if err != nil {
		fmt.Println("WriteAudit: unable to marshal object:", err)
		return
	}
	bAudit = append(bAudit, []byte(",\n")...)
	file, err := os.OpenFile(h.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer file.Close()
	_, err = file.Write(bAudit)
	if err != nil {
		fmt.Println("WriteAudit: unable to write to file:", err)
	}
}

func (h AuditSQL) WriteAudit(audit Audit) {
	if h.DB == nil {
		fmt.Println("WriteAudit: DB not set")
		return
	}
	auditColumn := AuditColumns{Created: audit.Created, Updated: audit.Updated, Delete: audit.Delete}
	bAuditColumn, err := json.Marshal(auditColumn)
	if err != nil {
		fmt.Println("WriteAudit: unable to marshal columns")
		return
	}
	h.CreatedAt = time.Now().UTC()
	h.Changed = bAuditColumn
	h.UserID = audit.UserID
	h.UserUID = audit.UserUID
	h.Entity = audit.Entity
	h.EntityID = audit.EntityID
	insertSql := `INSERT INTO audit (created_at, changed, user_id, user_uid, entity, entity_id) VALUES (:created_at, :changed, :user_id, :user_uid, :entity, :entity_id)`
	if _, err := h.DB.NamedExec(insertSql, h); err != nil {
		fmt.Println("WriteAudit: error insert record", err)
	}
}

func GroupStructToMap(entity any, tagName string) map[string]any {
	if tagName == "" {
		tagName = "field"
	}
	m := make(map[string]any)
	v := reflect.ValueOf(entity)
	t := reflect.TypeOf(entity)
	for i := 0; i < v.NumField(); i++ {
		name := t.Field(i).Name
		if tagName != "field" {
			structTag := t.Field(i).Tag
			name = structTag.Get(tagName)
			if name == "" {
				name = t.Field(i).Name
			}
		}
		if name != "-" {
			// if for any reason we are skipping the tag's transformation, skip it
			m[name] = v.Field(i).Interface()
		}
	}
	return m
}

func GroupStructToMapUpdated(entity any, tagName string, fields map[string]any) map[string]AuditUpdate {
	if tagName == "" {
		tagName = "field"
	}
	m := make(map[string]AuditUpdate)
	v := reflect.ValueOf(entity)
	t := reflect.TypeOf(entity)
	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		name := fieldName
		if tagName != "field" {
			structTag := t.Field(i).Tag
			name = structTag.Get(tagName)
			if name == "" {
				name = fieldName
			}
		}
		if val, ok := fields[name]; ok {
			if name != "-" {
				// if for any reason we are skipping the tag's transformation, skip it
				m[name] = AuditUpdate{To: v.Field(i).Interface(), From: val}
			}
		}
	}
	return m
}

func KeysToString(keys ...any) string {
	// this assumes keys will be (string, any, string, any)
	output := []string{}
	field := ""
	for i, k := range keys {
		if i%2 == 0 {
			field = k.(string)
			continue
		}
		output = append(output, fmt.Sprintf("%s: %s", field, k))
		field = ""
	}
	return strings.Join(output, ", ")
}
