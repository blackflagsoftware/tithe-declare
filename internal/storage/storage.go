package storage

import (
	"fmt"

	"github.com/blackflagsoftware/tithe-declare/config"
	l "github.com/blackflagsoftware/tithe-declare/internal/middleware/logging"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var sqliteDB *sqlx.DB

func InitStorage() *sqlx.DB {
	if sqliteDB == nil {
		var err error
		connStr := GetConnectionString()
		sqliteDB, err = sqlx.Open("sqlite3", connStr)
		if err != nil {
			l.Default.Panicf("Could not connect to the DB host: %s*****; %s", connStr[:6], err)
		}
		sqliteDB.SetMaxOpenConns(1)
	}
	return sqliteDB
}

func GetConnectionString() string {
	return fmt.Sprintf("%s?cache=shared&mode=wrc", config.DB.SqlitePath)
}

func FormatPagination(limit, offset int) string {
	if limit == 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d, %d", offset, limit)
}
