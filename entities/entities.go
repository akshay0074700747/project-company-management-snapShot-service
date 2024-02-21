package entities

type SnapMsg struct {
	Email     string `json:"email"`
	ProjectID string `json:"project_id"`
	Snapshots []Snap `json:"snapshots"`
}

type Snap struct {
	File        []byte `json:"file"`
	Filename    string `json:"file_name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Progress    int    `json:"progress"`
}

type SnapMessage struct {
	Email     string     `json:"email"`
	ProjectID string     `json:"project_id"`
	Snapshots []Snapshot `json:"snapshots"`
}

type Snapshot struct {
	Filename    string `json:"file_name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Progress    int    `json:"progress"`
}
