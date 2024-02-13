package models

type SubtaskMessage struct {
	Id     string    `json:"id" bson:"_id,omitempty"`
	Labels []*string `json:"labels" bson:"labels"`
}
