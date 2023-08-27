package models

import "mime/multipart"

type TestRequest struct {
	Firstname  string  `json:"fistname"`
	Lastname   string  `json:"lastname"`
	Patronymic *string `json:"patronymic"`
}

type TestResponse struct {
	ID string `json:"id"`
}

type FileRequest struct {
	Name        string                  `mapstructure:"name"`
	Files       []*multipart.FileHeader `mapstructure:"file"`
	OneMoreFile *multipart.FileHeader   `mapstructure:"one_more_file"`
}

type FileResponse struct {
	Name        string                  `mapstructure:"name" json:"name"`
	Files       []*multipart.FileHeader `mapstructure:"file[]" json:"files"`
	OneMoreFile *multipart.FileHeader   `mapstructure:"one_more_file" json:"one_more_file"`
}
