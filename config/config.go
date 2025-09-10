package config

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/kardianos/osext"
)

type (
	Service struct {
		AppName       string
		AppVersion    string
		RestPort      string
		GrpcPort      string
		Env           string
		StorageType   string
		PidPath       string
		ExecDir       string
		RootDir       string
		LogPath       string
		LogType       *os.File
		EnableMetrics bool
	}

	Migration struct {
		Enable   bool
		Dir      string
		SkipInit bool
	}

	Auditing struct {
		Enable   bool
		Storage  string
		FilePath string
	}

	BasicAuth struct {
		// BasicAuthUser        string
		// BasicAuthPwd         string
	}

	Database struct {
		Engine     string
		SqlitePath string
	}

	Email struct {
		Host       string
		Port       string
		Pwd        string
		From       string
		ResetUrl   string
		AdminEmail string
	}
	Auth struct {
		PwdCost              string
		ResetDuration        string
		ExpiresAtDuration    string
		AuthAlg              string
		AuthSecret           string
		AuthPublic           string
		EmailHost            string
		EmailPort            string
		EmailPwd             string
		EmailFrom            string
		EmailResetUrl        string
		Email                string
		BasicAuthUser        string
		BasicAuthPwd         string
		AuthorizationExpires string
		RefreshTokenExpires  string
	}
)

var (
	Srv Service
	Mig Migration
	Aud Auditing
	BA  BasicAuth
	DB  Database
	E   Email
	A   Auth
)

func init() {
	A = Auth{}
	E = Email{}
	DB = Database{}
	Srv = Service{AppName: "tithe-declare", LogType: os.Stdout}
	Srv.ExecDir, _ = osext.ExecutableFolder()
	Mig = Migration{}
	Aud = Auditing{}
	BA = BasicAuth{}
	loadEnvFiles()
	loadEnvVars()
}

func loadEnvVars() {
	Srv.RootDir = GetEnvOrDefault("TITHE_DECLARE_ROOT_DIR", ".")
	Srv.AppVersion = GetEnvOrDefault("TITHE_DECLARE_APP_VERSION", "1.0.0")
	Srv.RestPort = GetEnvOrDefault("TITHE_DECLARE_REST_PORT", "12580")
	Srv.GrpcPort = GetEnvOrDefault("TITHE_DECLARE_GRPC_PORT", "12581")
	Srv.PidPath = GetEnvOrDefault("TITHE_DECLARE_PID_PATH", fmt.Sprintf("/tmp/%s.pid", Srv.AppName))
	Srv.Env = GetEnvOrDefault("TITHE_DECLARE_ENV", "dev")
	Srv.LogPath = GetEnvOrDefault("TITHE_DECLARE_LOG_PATH", fmt.Sprintf("/tmp/%s.out", Srv.AppName))
	Srv.EnableMetrics = GetEnvOrDefaultBool("TITHE_DECLARE_ENABLE_METRICS", true)
	Mig.Enable = GetEnvOrDefaultBool("TITHE_DECLARE_MIGRATION_ENABLED", false)
	Mig.Dir = GetEnvOrDefault("TITHE_DECLARE_MIGRATION_PATH", "")
	Mig.SkipInit = GetEnvOrDefaultBool("TITHE_DECLARE_MIGRATION_SKIP_INIT", false)
	Aud.Enable = GetEnvOrDefaultBool("TITHE_DECLARE_ENABLE_AUDITING", false)
	Aud.Storage = GetEnvOrDefault("TITHE_DECLARE_AUDIT_STORAGE", "file") // file or sql
	Aud.FilePath = GetEnvOrDefault("TITHE_DECLARE_AUDIT_FILE_PATH", "./audit")
	// BA.BasicAuthUser = GetEnvOrDefault("TITHE_DECLARE_BASIC_AUTH_USER", "test")
	// BA.BasicAuthPwd = GetEnvOrDefault("TITHE_DECLARE_BASIC_AUTH_PWD", "test")
	DB.Engine = GetEnvOrDefault("TITHE_DECLARE_SQLITE_DB_ENGINE", "sqlite")
	DB.SqlitePath = GetEnvOrDefault("TITHE_DECLARE_SQLITE_PATH", "")
	E.Host = GetEnvOrDefault("TITHE_DECLARE_EMAIL_HOST", "")
	E.Port = GetEnvOrDefault("TITHE_DECLARE_EMAIL_PORT", "587")
	E.Pwd = GetEnvOrDefault("TITHE_DECLARE_EMAIL_PWD", "")
	E.From = GetEnvOrDefault("TITHE_DECLARE_EMAIL_FROM", "")
	E.ResetUrl = GetEnvOrDefault("TITHE_DECLARE_EMAIL_RESET_URL", "")
	E.AdminEmail = GetEnvOrDefault("TITHE_DECLARE_ADMIN_EMAIL", "")
	A.PwdCost = GetEnvOrDefault("TITHE_DECLARE_PWD_COST", "10")                       // algorithm cost
	A.ResetDuration = GetEnvOrDefault("TITHE_DECLARE_RESET_DURATION", "7")            // in days
	A.ExpiresAtDuration = GetEnvOrDefault("TITHE_DECLARE_EXPIRES_AT_DURATION", "168") // in hours (7 days)
	A.AuthAlg = GetEnvOrDefault("TITHE_DECLARE_AUTH_ALG", "HMAC")                     // HMAC, RSA, ECDSA or EdDSA (only use the 512 size)
	A.AuthSecret = GetEnvOrDefault("TITHE_DECLARE_AUTH_SECRET", "")                   // base64 format: used by all 3, HMAC is insecure use only for dev
	A.AuthPublic = GetEnvOrDefault("TITHE_DECLARE_AUTH_PUBLIC", "")                   // base64 format: only used by RSA or ECDSA
	A.BasicAuthUser = GetEnvOrDefault("TITHE_DECLARE_BASIC_AUTH_USER", "")
	A.BasicAuthPwd = GetEnvOrDefault("TITHE_DECLARE_BASIC_AUTH_PASS", "")
	A.AuthorizationExpires = GetEnvOrDefault("TITHE_DECLARE_AUTHORIZATION_EXPIRES", "60")   // in seconds
	A.RefreshTokenExpires = GetEnvOrDefault("TITHE_DECLARE_REFRESH_TOKEN_EXPIRES", "86400") // in seconds, set -1 to never expire; 0 - to always refresh; >0 in seconds to expire at
}

func loadEnvFiles() {
	// load any and all .env.* files if present at the root level of the binary
	// the order of precedence goes from local to just plain .env, the later WILL override earlier
	// if the env var is already declared at the console/terminal level, this will override it too
	rootDir := path.Join(Srv.ExecDir, "..", "..")
	rootDir = path.Clean(rootDir)
	envFiles := []string{".env.local", ".env.dev", ".env.qa", ".env.prod", ".env"} // if you want another name, just add it to this list in the order of precedence
	for _, envFile := range envFiles {
		envFileFull := path.Join(rootDir, envFile)
		if _, err := os.Stat(envFileFull); !os.IsNotExist(err) {
			// found the file
			loadEnvVarsFromFile(envFileFull)
		}
	}
}

func loadEnvVarsFromFile(fileName string) {
	// line format: name=value
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("env file: %s could not be read\n", fileName)
		return
	}
	for line := range bytes.SplitSeq(fileContent, []byte("\n")) {
		if len(line) > 0 {
			if line[0] != byte('#') {
				// not a comment
				split := bytes.Split(line, []byte("="))
				if len(split) == 2 {
					os.Setenv(string(split[0]), string(split[1]))
				}
			}
		}
	}
}

func GetEnvOrDefault(envVar string, defEnvVar string) (newEnvVar string) {
	if newEnvVar = os.Getenv(envVar); len(newEnvVar) == 0 {
		return defEnvVar
	} else {
		return newEnvVar
	}
}

func GetEnvOrDefaultBool(envVar string, defEnvVar bool) (newEnvVar bool) {
	newEnvVarStr := os.Getenv(envVar)
	if len(newEnvVarStr) == 0 {
		return defEnvVar
	}
	return newEnvVarStr == "true"
}

func KnownVersions() []string {
	return []string{"v1"}
}

func ConvertEnvVarStringToInt(strVal, VarName string, defaultInt int) int {
	convertedInt, err := strconv.Atoi(strVal)
	if err != nil {
		// TODO: if unable to print to default log, might want to send error to another feedback loop
		fmt.Printf("%s: unable to parse env var: %s", VarName, err)
		return defaultInt
	}
	return convertedInt
}

func (s Service) GetUniqueNumberForLock() (number int) {
	for i := range s.AppName {
		number += int(s.AppName[i])
	}
	return
}

func (e Email) GetEmailPort() int {
	emailPort := ConvertEnvVarStringToInt(e.Port, "EmailPort", 5432) // postgres
	return emailPort
}

func (a Auth) GetPwdCost() int {
	cost := ConvertEnvVarStringToInt(a.PwdCost, "PwdCost", 10)
	return cost
}

func (a Auth) GetResetDuration() int {
	durationInDays := ConvertEnvVarStringToInt(a.ResetDuration, "ResetDuration", 7)
	return durationInDays
}

func (a Auth) GetExpiresAtDuration() int {
	durationInHours := ConvertEnvVarStringToInt(a.ExpiresAtDuration, "ExpiresAtDuration", 7)
	return durationInHours
}

func (a Auth) GetAuthorizationExpires() int {
	authExpires := ConvertEnvVarStringToInt(a.AuthorizationExpires, "AuthorizationExpires", 3600)
	return authExpires * 60 * 60
}

func (a Auth) GetRefreshTokenExpires() int {
	refreshToken := ConvertEnvVarStringToInt(a.RefreshTokenExpires, "RefreshTokenExpires", 86400)
	return refreshToken
}
