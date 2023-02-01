package models

type Response struct {
	Error     bool        `json:"error" default:"false"`
	ErrorText string      `json:"error_text"`
	Data      interface{} `json:"data"`
	Code      int         `json:"code"`
}
