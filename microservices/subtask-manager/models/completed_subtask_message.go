package models

type CompletedSubtaskMessage struct {
	Id             string    `json:"id"`
	AssignedLabels []*string `json:"assignedLabels"`
}
