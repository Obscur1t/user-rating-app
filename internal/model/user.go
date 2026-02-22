package model

type User struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	NickName string  `json:"nickname"`
	Likes    int     `json:"likes"`
	Viewers  int     `json:"viewers"`
	Rating   float64 `json:"rating"`
}

func NewUser(name, nickname string, likes, viewers int) *User {
	return &User{
		Name:     name,
		NickName: nickname,
		Likes:    likes,
		Viewers:  viewers,
	}
}
