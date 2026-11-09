package model

type User struct {
	ID   int    `json:"id" validate:"required,min=1"`
	Name string `json:"name" validate:"required,min=2"`
}

type UserUpdateRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

type Response struct {
	RespCode    string `json:"respCode"`
	RespMessage string `json:"respMessage"`
	RespData    any    `json:"respData"`
}
