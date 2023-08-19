package matchingvalidator

import (
	"fmt"
	"gameapp/dto"
	"gameapp/entity"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) ValidateAddToWaitingListRequest(req dto.AddToWaitingListRequest) (map[string]string, error) {
	const op = "matchingvalidator.AddToWaitingListRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Category, validation.Required, validation.By(isCategoryValid)),
	); err != nil {
		fieldErrors := make(map[string]string)

		errV, ok := err.(validation.Errors)
		if ok {
			for key, value := range errV {
				if value != nil {
					fieldErrors[key] = value.Error()
				}
			}
		}

		return fieldErrors, richerror.New(op).WithMessage(errmsg.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"req": req}).WithErr(err)
	}
	return nil, nil
}

func isCategoryValid(value interface{}) error {
	category := value.(entity.Category)
	if !category.IsValid() {
		return fmt.Errorf(errmsg.ErrorMsgInvalidInput)
	}

	return nil
}
