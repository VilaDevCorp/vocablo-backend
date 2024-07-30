package utils

type Page[T any] struct {
	PageNumber int  `json:"pageNumber"`
	HasNext    bool `json:"hasNext"`
	Content    []T  `json:"content"`
}
