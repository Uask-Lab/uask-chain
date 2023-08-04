package types

const (
	OK  = 100
	Err = 101
)

type Response struct {
	Code    int `json:"code"`
	Payload any `json:"payload"`
}

func Ok(payload any) *Response {
	return &Response{
		Code:    OK,
		Payload: payload,
	}
}

func Error(err error) *Response {
	return &Response{
		Code:    Err,
		Payload: err,
	}
}
