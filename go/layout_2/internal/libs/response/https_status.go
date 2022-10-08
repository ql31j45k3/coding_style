package response

type HttpStatus int

const (
	HttpStatusOk HttpStatus = 200

	HttpStatusBadRequest HttpStatus = 100400
	HttpStatusNotFound   HttpStatus = 100404

	HttpStatusInternalServerError HttpStatus = 100500
)
