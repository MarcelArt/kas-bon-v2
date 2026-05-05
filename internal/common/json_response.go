package common

import "fmt"

type JSONResponse struct {
	Items     any    `json:"items"`
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

func NewJSONResponse(items any, message string) *JSONResponse {
	err, ok := items.(error)
	if ok {
		if message == "" {
			message = err.Error()
		} else {
			message = fmt.Sprintf("%s: %s", message, err.Error())
		}
		return &JSONResponse{
			Items:     nil,
			IsSuccess: false,
			Message:   message,
		}
	}

	return &JSONResponse{
		Items:     items,
		IsSuccess: true,
		Message:   message,
	}
}
