package model

type Category struct {
	ID        string  `json:"id"`
	ParentID  *string `json:"parentId,omitempty"`
	Name      string  `json:"name" validate:"required,min=2,max=100"`
	Type      string  `json:"type" validate:"required,oneof=expense income"`
	Icon      string  `json:"icon,omitempty" validate:"omitempty,max=50"`
	SortOrder int     `json:"sortOrder"`
}

type CreateCategoryRequest struct {
	ParentID  *string `json:"parentId,omitempty"`
	Name      string  `json:"name" validate:"required,min=2,max=100"`
	Type      string  `json:"type" validate:"required,oneof=expense income"`
	Icon      string  `json:"icon,omitempty" validate:"omitempty,max=50"`
	SortOrder int     `json:"sortOrder"`
}

type UpdateCategoryRequest struct {
	Name      string `json:"name" validate:"omitempty,min=2,max=100"`
	Icon      string `json:"icon,omitempty" validate:"omitempty,max=50"`
	SortOrder *int   `json:"sortOrder,omitempty"`
}
