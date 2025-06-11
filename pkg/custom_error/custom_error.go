package custom_error

import (
	"errors"

	Constants "mini-accounting/constants"
	Library "mini-accounting/library"
)

type CustomError struct {
	display     error
	plain       error
	path        string
	library     Library.Library
	code        *string
	description *string
	status      *string
	reference   *string
}

func New(
	display error,
	plain error,
	path string,
	library Library.Library,
) error {
	return &CustomError{
		display: display,
		plain:   plain,
		path:    path,
		library: library,
	}
}

func (e *CustomError) Error() string {
	message := map[string]interface{}{
		"display": e.display.Error(),
		"plain":   e.plain.Error(),
		"path":    e.path,
		"code":    e.code,
	}

	result, err := e.library.JsonMarshal(message)
	if err != nil {
		return err.Error()
	}

	return string(result)
}

func (e *CustomError) GetDisplay() error {
	return e.display
}

func (e *CustomError) GetPlain() error {
	return e.plain
}

func (e *CustomError) GetPath() string {
	return e.path
}

func (e *CustomError) GetCode() string {
	return *e.code
}

func (e *CustomError) GetDescription() *string {
	return e.description
}

func (e *CustomError) GetStatus() *string {
	return e.status
}

func (e *CustomError) GetReference() *string {
	return e.reference
}

func (e *CustomError) UnshiftPath(path string) error {
	e.path = path + " > " + e.path
	return e
}

func (e *CustomError) FromListMap(errs []map[string]interface{}) error {
	result, err := e.library.JsonMarshal(errs)
	if err != nil {
		return New(
			Constants.ErrFailedJSONMarshal,
			err,
			"CustomError:FromListMap",
			e.library,
		)
	}

	e.plain = errors.New(string(result))
	return e
}

func (e *CustomError) SetCode(code string) {
	e.code = &code
}

func (e *CustomError) SetDescription(description string) {
	e.description = &description
}

func (e *CustomError) SetStatus(status string) {
	e.status = &status
}

func (e *CustomError) SetReference(reference string) {
	e.reference = &reference
}
