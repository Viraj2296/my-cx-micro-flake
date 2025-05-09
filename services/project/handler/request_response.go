package handler

type CreateProjectCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AvatarUrl   string `json:"avatarUrl"`
	IsDefault   bool   `json:"isDefault"`
}

type UpdateProjectCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AvatarUrl   string `json:"profileUrl"`
	IsDefault   bool   `json:"isDefault"`
}
