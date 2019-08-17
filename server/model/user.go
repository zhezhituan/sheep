package model

type User struct {
	Uid  int    `json:"Uid"`
	Name string `json:"Name"`
	Pw   string `json:"Pw"`
}
