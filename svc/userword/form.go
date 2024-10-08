package userword

import "vocablo/schema"

type CreateForm struct {
	Term        string              `json:"term" binding:"required"`
	Definitions []schema.Definition `json:"definitions" binding:"required"`
	Lang        string              `json:"lang" binding:"required"`
}

type UpdateForm struct {
	ID          string               `json:"id" binding:"required"`
	Term        *string              `json:"term" `
	Definitions *[]schema.Definition `json:"definitions" `
}

type SearchForm struct {
	Term     *string `json:"term"`
	Lang     *string `json:"lang"`
	Learned  *bool   `json:"learned"`
	OrderBy  *string `json:"orderBy"`
	OrderDir *string `json:"orderDir"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
	Count    bool    `json:"count"`
}
