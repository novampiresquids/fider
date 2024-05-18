package entity

import (
	"github.com/getfider/fider/app/models/enum"
)

// User represents an user inside our application
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Tenant *Tenant `json:"-"`
	BoardRole     enum.Role         `json:"role"`
	Email         string            `json:"-"`
	Providers     []*UserProvider   `json:"-"`
	AvatarBlobKey string            `json:"-"`
	AvatarType    enum.AvatarType   `json:"-"`
	AvatarURL     string            `json:"avatarURL,omitempty"`
	Status        enum.UserStatus   `json:"status"`
	Membership    []*UserMembership `json:"membership"`
}

// HasProvider returns true if current user has registered with given provider
func (u *User) HasProvider(provider string) bool {
	for _, p := range u.Providers {
		if p.Name == provider {
			return true
		}
	}
	return false
}

// IsCollaborator returns true if user has special permissions
func (u *User) Role(board *Tenant) enum.Role {
	if board == nil {
		return enum.RoleAdministrator // FIX TO CHECK FOR SITE ADMIN!! TODO
	}
	for _, m := range u.Membership {
		if m.Board.ID == board.ID {
			return m.Role
		}
	}
	return enum.RoleVisitor
}

// IsCollaborator returns true if user has special permissions
func (u *User) IsCollaborator(board *Tenant) bool {
	return u.Role(board) == enum.RoleCollaborator || u.Role(board) == enum.RoleAdministrator
}

// IsAdministrator returns true if user is administrator
func (u *User) IsAdministrator(board *Tenant) bool {
	return u.Role(board) == enum.RoleAdministrator
}

// UserProvider represents the relationship between an User and an Authentication provide
type UserProvider struct {
	Name string
	UID  string
}

// UserMembership represents the boards that the user is a member of
type UserMembership struct {
	Board *Tenant   `json:"-"`
	Role  enum.Role `json:"role"`
}
