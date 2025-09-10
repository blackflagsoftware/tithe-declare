package util

import (
	"strings"
	"unicode"

	"github.com/blackflagsoftware/tithe-declare/config"
	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	"github.com/blackflagsoftware/tithe-declare/internal/middleware/logging"
	"golang.org/x/crypto/bcrypt"
)

func PasswordValidator(pwd, confirm string) error {
	if pwd != confirm {
		return ae.PasswordValiationError("passwords do not match")
	}
	errors := []string{}
	hasNumber, hasUpper, hasSpecial := false, false, false
	for _, r := range pwd {
		switch {
		case unicode.IsNumber(r):
			hasNumber = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}
	// character length
	if len(pwd) < 8 {
		errors = append(errors, "must be > 8 characters long")
	}
	// character one upper case
	if !hasUpper {
		errors = append(errors, "must have one uppercase letter")
	}
	// character one digit
	if !hasNumber {
		errors = append(errors, "must have one number")
	}
	// character one special
	if !hasSpecial {
		errors = append(errors, "must have one special character [@$!%*?]")
	}
	if len(errors) > 0 {
		return ae.PasswordValiationError(strings.Join(errors, ", "))
	}
	return nil
}

func EncryptPassword(pwd string) (string, error) {
	pwdByte := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(pwdByte, config.A.GetPwdCost())
	if err != nil {
		logging.Default.Println("EncryptPassword: error on generating")
		return "", ae.GeneralError("encrypt password failed", err)
	}
	return string(hash), nil
}

func CheckPassword(pwd, hash string) error {
	pwdByte := []byte(pwd)
	hashByte := []byte(hash)
	if err := bcrypt.CompareHashAndPassword(hashByte, pwdByte); err != nil {
		return ae.EmailPasswordComboError()
	}
	return nil
}
