package domain

// Category represents a spending category (global or user-defined)
// user_id nil = global category; non-nil = user-defined
type Category struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	UserID *string `json:"user_id,omitempty"`
}

// CreateCategoryInput is the input for creating a category
type CreateCategoryInput struct {
	Name   string  `json:"name"`
	UserID *string `json:"user_id,omitempty"`
}

// UpdateCategoryInput is the input for updating a category (e.g. name)
type UpdateCategoryInput struct {
	Name *string `json:"name,omitempty"`
}
