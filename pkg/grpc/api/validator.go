package api

import (
	"errors"
	"regexp"

	"github.com/douglaszuqueto/go-grpc-user/proto"
)

var (
	usernameValidator = regexp.MustCompile("^[[:alnum:]]+$")
	passwordValidator = regexp.MustCompile("^.{6,}$")
)

var (
	errUserInvalidUsername = errors.New("Username - apenas caracteres [a-zA-Z0-9] são aceitos")
	errUserPasswordLength  = errors.New("Senha - no mínimo seis caracteres")
)

func validateUsername(text string) error {
	if !usernameValidator.MatchString(text) {
		return errUserInvalidUsername
	}
	return nil
}

func validatePassword(text string) error {
	if !passwordValidator.MatchString(text) {
		return errUserPasswordLength
	}
	return nil
}

func storeValidateOrFail(u *proto.User) error {
	if len(u.Username) == 0 {
		return errors.New("Username deve ser preenchido")
	}

	if len(u.Password) == 0 {
		return errors.New("Senha deve ser preenchido")
	}

	if err := validateUsername(u.Username); err != nil {
		return err
	}

	return validatePassword(u.Password)
}
