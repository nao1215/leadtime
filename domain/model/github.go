// Package model is domain model and business logic.
package model

// Repository represents GitHub repository information
type Repository struct {
	// ID is repository id
	ID *int64 `json:"id,omitempty"`
	// Owner is repository owner
	Owner *User `json:"owner,omitempty"`
	// Name is repository name
	Name *string `json:"name,omitempty"`
	// FullName is repository full name
	FullName *string `json:"full_name,omitempty"`
	// Description is repository description
	Description *string `json:"description,omitempty"`
}

// User represents a GitHub user.
type User struct {
	// Name is user name.
	Name *string `json:"name,omitempty"`
}
