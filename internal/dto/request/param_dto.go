package request

type PaginationQuery struct {
	Limit  int
	Offset int
	Sort   string
}

func NewPaginationQuery(limit int, offset int, sort string) PaginationQuery {
	return PaginationQuery{
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
}
