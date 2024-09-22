package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/kumersun/bnovo/entity"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type GuestValidator struct {
	guest *entity.Guest
}

func NewGuestValidator(guest *entity.Guest) *GuestValidator {
	return &GuestValidator{guest}
}

func (v *GuestValidator) Validate() error {
	if err := validate.Struct(v.guest); err != nil {
		return err
	}

	if v.guest.Email != "" {
		if err := validate.Var(v.guest.Email, "email"); err != nil {
			return errors.New("Email is not valid")
		}
	}

	if v.guest.Country != "" {
		if err := validate.Var(v.guest.Country, "country_code"); err != nil {
			return errors.New("Country code is not valid")
		}
	}

	return nil
}
