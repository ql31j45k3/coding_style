package response

type Status struct {
	code     HttpStatus
	messages string
}

func (s *Status) Code() HttpStatus {
	return s.code
}

func (s *Status) Messages() string {
	if s == nil {
		return ""
	}

	return s.messages
}

func (s *Status) Err() error {
	if s.code == HttpStatusOk {
		return nil
	}

	return s
}

func (s *Status) Error() string {
	return s.Messages()
}

func New(code HttpStatus, messages string) *Status {
	return &Status{
		code:     code,
		messages: messages,
	}
}

func Err(code HttpStatus, messages string) error {
	return New(code, messages).Err()
}

func FromError(err error) (*Status, bool) {
	if err == nil {
		return nil, false
	}

	if s, ok := err.(*Status); ok {
		return s, true
	}

	return nil, false
}

type HttpStatus int

const (
	HttpStatusOk HttpStatus = 200

	HttpStatusBadRequest HttpStatus = 400
	HttpStatusNotFound   HttpStatus = 404

	HttpStatusInternalServerError HttpStatus = 500
)

const (
	Code100401 HttpStatus = 100401
)
