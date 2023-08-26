package common

import (
	"TinyTik/model"
	"github.com/go-playground/validator"
)

func ValidateUserAuth(userAuth model.UserAuth) error {
	validate := validator.New()
	if err := validate.Struct(userAuth); err != nil {
		return err
	}
	return nil
}
