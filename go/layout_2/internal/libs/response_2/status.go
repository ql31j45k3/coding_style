package response_2

type Status struct {
	Code     int    `json:"code"`
	Messages string `json:"messages"`
}

func (s *Status) WithMsg(msg string) Status {
	return Status{
		Code:     s.Code,
		Messages: s.Messages + msg,
	}
}

func new(code int, messages string) Status {
	return Status{
		Code:     code,
		Messages: messages,
	}
}

var (
	codeOk = new(200, "success")

	CodeInternalError = new(10400, "internal error")
)
