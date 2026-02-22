package request

type UserRequestDTO struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Likes    int    `json:"likes"`
	Viewers  int    `json:"viewers"`
}

type UpdateUserDTO struct {
	Name     *string `json:"name"`
	Nickname *string `json:"nickname"`
	Likes    *int    `json:"likes"`
	Viewers  *int    `json:"viewers"`
}
