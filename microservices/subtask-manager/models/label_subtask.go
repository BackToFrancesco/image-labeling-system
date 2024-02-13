package models

type LabelSubtask struct {
	UserId             string           `bson:"_id,omitempty" json:"userId"`
	ImageId            string           `bson:"ImageId" json:"imageId"`
	AssignedLabel      string           `json:"assignedLabel" bson:"assignedLabels"`
}
