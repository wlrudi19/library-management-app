package response

import (
	"encoding/json"
	"net/http"
	"time"
)

type Meta struct {
	StatusCode int    `json:"statusCode,omitempty"`
	Success    bool   `json:"success,omitempty"`
	Message    string `json:"message,omitempty"`
}

type Response struct {
	AccessTime time.Time   `json:"accessTime,omitempty"`
	Data       interface{} `json:"output,omitempty"`
	Meta
}

type Json struct {
	Response Response `json:"response"`
	HttpCode int
}

func CustomBuilder(httpCode int, success bool, data interface{}, message string) *Json {
	return &Json{
		Response: Response{
			Data:       data,
			AccessTime: time.Now(),
			Meta: Meta{
				StatusCode: httpCode,
				Message:    message,
				Success:    success,
			},
		},
		HttpCode: httpCode,
	}
}

func (j *Json) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(j.HttpCode)
	json.NewEncoder(w).Encode(j.Response)
}
