package models

type Subtask struct {
	Id             string         `bson:"_id,omitempty"`
	Labels         []string       `bson:"labels"`
	Assignee       []string       `bson:"assignee"`
	AssignedLabels map[string]int `bson:"assignedLabels"`
}

/*package models

type Subtask struct {
	Id    string `bson:"id"`
	Label string `bson:"label"`
}*/
