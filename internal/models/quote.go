package models

type Quote struct {
	Id     string `json:"id"`
	Text   string `json:"text"`
	Author string `json:"author"`
}
