package models

type CompletedSubtaskMessage struct {
	Id             string          `json:"id"`
	AssignedLabels *map[string]int `json:"assignedLabels"`
}
