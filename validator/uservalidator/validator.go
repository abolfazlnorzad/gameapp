package uservalidator

import (
	"errors"
	"gameapp/dto"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
}

type Validator struct {
	repo Repository
}

func New(r Repository) Validator {
	return Validator{
		repo: r,
	}
}

func (v Validator) ValidateRegisterRequest(req dto.RegisterRequest) (map[string]string, error) {
	const op = "uservalidator"
	if vErr := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 77)),
		validation.Field(&req.PhoneNumber, validation.Required, validation.By(v.checkPhoneNumberIsUnique)),
		validation.Field(&req.Password, validation.Required),
	); vErr != nil {
		fieldErrors := make(map[string]string)

		errV, ok := vErr.(validation.Errors)
		if ok {
			for key, value := range errV {
				if value != nil {
					fieldErrors[key] = value.Error()
				}
			}
		}

		return fieldErrors, richerror.New(op).WithMessage(errmsg.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"req": req}).WithErr(errV)
	}
	return nil, nil
}

func (v Validator) checkPhoneNumberIsUnique(value interface{}) error {
	phoneNumber := value.(string)
	// check uniqueness phone number
	if isUnique, err := v.repo.IsPhoneNumberUnique(phoneNumber); !isUnique || err != nil {
		if err != nil {
			return err
		}
		if !isUnique {
			return errors.New("phone number is not unique .")
		}
	}
	return nil
}
