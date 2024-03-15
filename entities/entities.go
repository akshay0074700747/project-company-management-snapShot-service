package entities

type SnapMsg struct {
	Email     string `json:"email"`
	ProjectID string `json:"project_id"`
	Snapshots []Snap `json:"snapshots"`
	Key       string `json:"key"`
	IsStaged  bool   `json:"is_staged"`
}

type Snap struct {
	File        []byte `json:"file"`
	Filename    string `json:"file_name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Progress    int    `json:"progress"`
}

type SnapMessage struct {
	UserID    string     `json:"user_id" bson:"user_id"`
	ProjectID string     `json:"project_id" bson:"project_id"`
	Snapshots []Snapshot `json:"snapshots" bson:"snapshots"`
}

type Snapshot struct {
	Filename    string `json:"file_name" bson:"filename"`
	Key         string `json:"key" bson:"key"`
	Description string `json:"description" bson:"description"`
	Progress    int    `json:"progress" bson:"progress"`
	IsStaged    bool   `json:"is_staged" bson:"is_staged"`
}

type StageRes struct {
	Stages  int             `json:"stages"`
	Details []StagesDetails `json:"details"`
}

type StagesDetails struct {
	Key         string `json:"key" bson:"key"`
	Description string `json:"description" bson:"description"`
	Filename    string `json:"file_name" bson:"filename"`
}

type ListofUserProgress struct {
	UserAndProgress []UserProgress `json:"user_and_progress"`
}

type UserProgress struct {
	UserID string `json:"user_id" bson:"user_id"`
	Stages int    `json:"stages" bson:"stages"`
}

type GetSnapshot struct {
	Snap []byte `json:"snap"`
}
