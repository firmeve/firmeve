package parser

import (
	"mime/multipart"
	"net/http"
)

type (
	MultipartFormData = *multipart.Form

	MultipartForm struct {
		original MultipartFormData
	}
)

func (j *MultipartForm) Bind(v interface{}) error {
	return FormDecoder.Decode(v, j.original.Value)
}

func (j *MultipartForm) Has(key string) bool {
	_, ok := j.original.Value[key]
	return ok
}

func (j *MultipartForm) Get(key string) interface{} {
	if j.Has(key) {
		return j.original.Value[key]
	}

	return nil
}

func (j *MultipartForm) File(key string) (multipart.File, *multipart.FileHeader, error) {
	if j.original.File != nil {
		if fhs := j.original.File[key]; len(fhs) > 0 {
			f, err := fhs[0].Open()
			return f, fhs[0], err
		}
	}
	return nil, nil, http.ErrMissingFile
}

func NewMultipartForm(data MultipartFormData) *MultipartForm {
	return &MultipartForm{
		original: data,
	}
}
