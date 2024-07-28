package word

import "vocablo/schema"

type CreateForm struct {
	Term        string              `json:"term" binding:"required"`
	Definitions []schema.Definition `json:"definitions" binding:"required"`
	Lang        string              `json:"lang" binding:"required"`
}
