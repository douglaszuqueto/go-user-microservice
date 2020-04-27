package api

import (
	"errors"
	"regexp"

	"github.com/douglaszuqueto/go-grpc-user/proto"
)

var (
	usernameValidator = regexp.MustCompile("^[[:alnum:]]+$")
)

var (
	errUserInvalidUsername = errors.New("Username - apenas caracteres [a-zA-Z0-9] são aceitos")
)

func validateUsername(text string) error {
	if !usernameValidator.MatchString(text) {
		return errUserInvalidUsername
	}
	return nil
}

func storeValidateOrFail(u *proto.User) error {
	if len(u.Username) == 0 {
		return errors.New("Username deve ser preenchido")
	}

	// if err := validateUsername(u.Username); err != nil {
	// 	return err
	// }

	return validateUsername(u.Username)
}
