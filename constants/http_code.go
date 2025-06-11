package constants

const (
	// 200
	StatusOK       = 200
	StatusCreated  = 201
	StatusAccepted = 202

	// 400
	StatusBadRequest            = 400
	StatusAuthenticationFailuer = 401
	StatusForbidden             = 403
	StatusNotFound              = 404
	StatusMethodNotAllowed      = 405
	StatusConflict              = 409
	StatusPreconditionFailed    = 412
	StatusRequestEntityTooLarge = 413

	// 500
	StatusInternalServerError = 500
	StatusNotImplemented      = 501
	StatusServiceUnavailable  = 503
)
