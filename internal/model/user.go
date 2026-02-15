package model

type User struct {
	Id       int64
	Name     string
	NickName string
	Likes    int
	Viewers  int
	Rating   float64
}

func NewUser(name, nickname string, likes, viewers int) *User {
	return &User{
		Name:     name,
		NickName: nickname,
		Likes:    likes,
		Viewers:  viewers,
	}
}
