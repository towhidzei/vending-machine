package httputil

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
}

func MakeResponse(code int, message string, content interface{}) BaseResponse {
	return BaseResponse{
		Code:    code,
		Message: message,
		Content: content,
	}
}
