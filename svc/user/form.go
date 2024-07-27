package user

import "github.com/google/uuid"

type CreateForm struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateForm struct {
	Id       uuid.UUID `json:"id" binding:"required"`
	Password *string   `json:"password"`
}

type SearchForm struct {
	Name     *string `json:"name"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
}
