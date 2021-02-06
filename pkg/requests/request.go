package requests

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"net/http"
)

type SlaveParams struct {
	Port      int    `json:"port" binding:"required"`
	StaticDir string `json:"static_dir" binding:"required"`
}

func (params *SlaveParams) Bind(r *http.Request) error {
	return validation.ValidateStruct(params,
		validation.Field(&params.Port, validation.Required),
		validation.Field(&params.StaticDir, validation.Required))
}
