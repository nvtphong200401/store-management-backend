package models

type PaginationModel[T any] struct {
	Data       []T
	TotalItems int64
}
