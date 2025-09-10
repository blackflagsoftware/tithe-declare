package src

import (
	"fmt"

	"github.com/blackflagsoftware/tithe-declare/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type (
	Postgres struct{}
)

func (p *Postgres) ConnectDB(c Connection, rootDB bool) (*sqlx.DB, error) {
	var db *sqlx.DB
	dbName := c.DB
	user := c.User
	pwd := c.Pwd
	if rootDB {
		dbName = "postgres"
		if c.AdminUser != "" {
			user = c.AdminUser
		}
		if c.AdminPwd != "" {
			pwd = c.AdminPwd
		}
	}
	conn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", user, pwd, dbName, c.Host)
	if pwd == "" {
		conn = fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable", user, dbName, c.Host)
	}
	db, errOpen := sqlx.Open("postgres", conn)
	if errOpen != nil {
		return db, fmt.Errorf("ConnectDB[postgres]: unable to open DB %s****; %s", conn[:6], errOpen)
	}
	return db, nil
}

func (p *Postgres) CheckUser(db *sqlx.DB, c Connection) error {
	checkSql := "SELECT EXISTS(SELECT rolname FROM pg_roles WHERE rolname = lower($1))"
	exists := false
	err := db.Get(&exists, checkSql, c.User)
	if err != nil {
		return fmt.Errorf("CheckUser[postgres]: unable to check for existing user; %s", err)
	}
	if !exists {
		createSql := fmt.Sprintf("CREATE role %s WITH LOGIN PASSWORD '%s'", c.User, c.Pwd)
		if _, err := db.Exec(createSql); err != nil {
			return fmt.Errorf("CheckUser[postgres]: unable to create role; %s", err)
		}
		role := []string{"pg_read_all_data", "pg_write_all_data"}
		for _, r := range role {
			grantSql := fmt.Sprintf("GRANT %s TO %s", r, c.User)
			if _, err := db.Exec(grantSql); err != nil {
				return fmt.Errorf("CheckUser[postgres]: unable to grant %s role; %s", r, err)
			}
		}
		// give postgres access to new role for DB creation
		grantSql := fmt.Sprintf("GRANT %s TO %s", c.User, c.AdminUser)
		if _, err := db.Exec(grantSql); err != nil {
			return fmt.Errorf("CheckUser[postgres]: unable to grant %s to role %s; %s", c.AdminUser, c.User, err)
		}
	}
	return nil
}

func (p *Postgres) CheckDB(db *sqlx.DB, c Connection) error {
	checkSql := "SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower($1))"
	exists := false
	err := db.Get(&exists, checkSql, c.DB)
	if err != nil {
		return fmt.Errorf("CheckDB[postgres]: unable to check for existing database; %s", err)
	}
	if !exists {
		createSql := fmt.Sprintf("CREATE DATABASE %s OWNER %s", c.DB, c.User)
		if _, err := db.Exec(createSql); err != nil {
			return fmt.Errorf("CheckDB[postgres]: unable to create database; %s", err)
		}
	}
	return nil
}

func (p *Postgres) CheckTable(db *sqlx.DB, c Connection) error {
	checkSql := "SELECT EXISTS(SELECT table_name FROM information_schema.tables WHERE table_name = 'migration')"
	exists := false
	err := db.Get(&exists, checkSql)
	if err != nil {
		return fmt.Errorf("CheckTable[postgres]: unable to check for existing table; %s", err)
	}
	if !exists {
		createSql := "CREATE TABLE migration (id serial, file_name varchar(100) NOT NULL)"
		if _, err := db.Exec(createSql); err != nil {
			return fmt.Errorf("CheckTable[postgres]: unable to create table; %s", err)
		}
	}
	return nil
}

func (p *Postgres) LockTable(db *sqlx.DB) bool {
	number := config.Srv.GetUniqueNumberForLock()
	succeed := false
	lockSql := "SELECT pg_try_advisory_lock($1)"
	if err := db.Get(&succeed, lockSql, number); err != nil {
		fmt.Printf("LockTable[postgres]: unable to lock resource; %s", err)
		return false
	}
	return succeed
}

func (p *Postgres) UnlockTable(db *sqlx.DB) error {
	number := config.Srv.GetUniqueNumberForLock()
	succeed := false
	lockSql := "SELECT pg_advisory_unlock($1)"
	if err := db.Get(&succeed, lockSql, number); err != nil {
		return fmt.Errorf("UnlockTable[postgres]: unable to unlock; %s", err)
	}
	if !succeed {
		return fmt.Errorf("UnlockTable[postgres]: unable to unlock with no error")
	}
	return nil
}
