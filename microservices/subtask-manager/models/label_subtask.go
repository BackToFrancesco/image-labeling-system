package models

type LabelImagePair struct {
	ImageId string `bson:"ImageId" json:"imageId"`
	Label   string `bson:"label" json:"label"`
}

type LabelSubtask struct {
	UserId             string           `bson:"_id,omitempty" json:"userId"`
	AssignedLabels     []*LabelImagePair `json:"assignedLabels" bson:"assignedLabels"`
}
