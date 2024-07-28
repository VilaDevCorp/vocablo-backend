package utils

type Page[T any] struct {
	PageNumber int
	HasNext    bool
	Content    []*T
}
