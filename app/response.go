package app

import (
	"encoding/json"
	"net/http"
)

type H map[string]string

type Response struct {
	http.ResponseWriter
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{w}
}

func (r *Response) JSON(status int, data interface{}) {
	b, _ := json.Marshal(data)

	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(status)
	r.Write(b)
	r.Write([]byte{'\n'})
}
