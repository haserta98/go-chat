package dto

type GroupCreateRequest struct {
	Name string `json:"name" binding:"required"`
}
