package models

type RequestSubtasks struct {
	UserId     string    `json:"userId"`
	NumberOfSubtasks int `json:"numberOfSubtasks"`
}
