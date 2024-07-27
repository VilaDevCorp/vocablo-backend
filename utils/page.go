package utils

type Page struct {
	TotalPages int         `json:"totalPages"`
	TotalRows  int         `json:"totalRows"`
	Content    interface{} `json:"content"`
}
