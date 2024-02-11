package models

type SubtaskMessage struct {
	Id     string    `json:"id"`
	Labels []*string `json:"labels"`
}
