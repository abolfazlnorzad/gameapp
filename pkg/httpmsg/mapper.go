package httpmsg

import (
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	"net/http"
)

func Error(err error) (int, string) {
	switch err.(type) {
	case richerror.RichError:
		rErr, _ := err.(richerror.RichError)
		msg := rErr.GetMessage()
		code := mapKindToHttpStatusCode(rErr.GetKind())
		if code >= 500 {
			msg = errmsg.SomethingWentWrong
		}
		return code, msg
	default:
		return http.StatusBadRequest, err.Error()
	}
}

func mapKindToHttpStatusCode(k richerror.Kind) int {
	switch k {
	case richerror.KindUnexpected:
		return http.StatusInternalServerError
	case richerror.KindForbidden:
		return http.StatusForbidden
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
