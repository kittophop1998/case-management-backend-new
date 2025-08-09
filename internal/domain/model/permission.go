package model

type UpdatePermissionRequest struct {
	Permission string   `json:"permission" binding:"required"`
	Roles      []string `json:"roles" binding:"required"`
}
