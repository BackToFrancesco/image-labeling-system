package models

type Subtask struct {
	Id             string         `bson:"_id,omitempty"`
	Labels         []*string       `bson:"labels"`
	Assignee       []string       `bson:"assignee"`
	AssignedLabels map[string]int `bson:"assignedLabels"`
}