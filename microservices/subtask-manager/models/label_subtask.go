package models

type LabelSubtask struct {
	UserId        string `bson:"_id,omitempty" json:"userId" binding:"required"`
	ImageId       string `bson:"ImageId" json:"imageId"`
	AssignedLabel string `json:"assignedLabel" bson:"assignedLabels" binding:"required"`
}
