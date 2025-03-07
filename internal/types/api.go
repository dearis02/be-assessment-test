package types

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Remark     string      `json:"remark,omitempty"`
	Errors     any         `json:"errors,omitempty"`
}

func (r ApiResponse) MarshalJSON() ([]byte, error) {
	type Alias ApiResponse

	if r.Message == "" {
		r.Message = http.StatusText(r.StatusCode)
	}
	return json.Marshal(&Alias{
		StatusCode: r.StatusCode,
		Message:    r.Message,
		Data:       r.Data,
		// not sure is this necessary?, cause there is already message field
		Remark: r.Remark,
		Errors: r.Errors,
	})
}
