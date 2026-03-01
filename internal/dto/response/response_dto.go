package responsedto

type PaginatedResponse[T any] struct {
	TotalCount int `json:"total_count"`
	Data       []T `json:"data"`
}

func NewPaginatedResponse[T any](data []T, totalCount int) PaginatedResponse[T] {
	return PaginatedResponse[T]{
		Data:       data,
		TotalCount: totalCount,
	}
}
