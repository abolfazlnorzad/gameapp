package uservalidator

import (
	"fmt"
	"gameapp/dto"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) ValidateLoginRequest(req dto.LoginRequest) (map[string]string, error) {
	const op = "uservalidator.login"
	if vErr := validation.ValidateStruct(&req,
		validation.Field(&req.PhoneNumber, validation.Required, validation.By(v.doesPhoneNumberExist)),
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

func (v Validator) doesPhoneNumberExist(value interface{}) error {
	const op = "uservalidation.login"
	phoneNumber := value.(string)
	_, err := v.repo.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		r, ok := err.(richerror.RichError)
		if ok {
			return r
		}
		return fmt.Errorf(errmsg.ErrorMsgNotFound)
	}

	return nil
}
