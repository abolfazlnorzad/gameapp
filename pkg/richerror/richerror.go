package richerror

type Kind int

type RichError struct {
	operation    string
	wrappedError error
	message      string
	kind         Kind
	meta         map[string]any
}

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
)

func (r RichError) Error() string {
	return r.message
}

func New(op string) RichError {
	return RichError{
		operation: op,
	}
}

func (r RichError) WithMessage(msg string) RichError {
	r.message = msg
	return r
}

func (r RichError) WithErr(err error) RichError {
	r.wrappedError = err
	return r
}

func (r RichError) WithMeta(m map[string]any) RichError {
	r.meta = m
	return r
}

func (r RichError) WithKind(k Kind) RichError {
	r.kind = k
	return r
}

func (r RichError) GetKind() Kind {
	if r.kind != 0 {
		return r.kind
	} else {
		rErr, ok := r.wrappedError.(RichError)
		if !ok {
			return 0
		}
		return rErr.GetKind()
	}
}

func (r RichError) GetMessage() string {
	if r.message != "" {
		return r.message
	} else {
		rErr, ok := r.wrappedError.(RichError)
		if !ok {
			return r.Error()
		}
		return rErr.GetMessage()
	}
}
