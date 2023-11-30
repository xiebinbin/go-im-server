package errors

import (
	"fmt"
)

type TmmSDKError struct {
	Code      int    `json:"code"`
	Message   string `json:"msg"`
	IsNetErr  uint8  `json:"is_net"`
	RequestId string `json:"request_id,omitempty"`
}

// NewTmmSDKError
func NewTmmSDKError(message string, code int) error {
	return &TmmSDKError{
		Code:     code,
		Message:  message,
		IsNetErr: 0,
	}
}

func (e *TmmSDKError) Error() string {
	if e.RequestId == "" {
		return fmt.Sprintf("[TmmSDKError] Code=%s, Message=%s", e.Code, e.Message)
	}
	return fmt.Sprintf("[TmmSDKError] Code=%s, Message=%s, RequestId=%s", e.Code, e.Message, e.RequestId)
}

func (e *TmmSDKError) GetCode() int {
	return e.Code
}

func (e *TmmSDKError) GetMessage() string {
	return e.Message
}

func (e *TmmSDKError) GetRequestId() string {
	return e.RequestId
}
