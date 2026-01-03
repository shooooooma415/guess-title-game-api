package user

import "time"

// User represents a user entity
type User struct {
	id        UserID
	name      UserName
	createdAt time.Time
}

// NewUser creates a new User
func NewUser(id UserID, name UserName) *User {
	return &User{
		id:        id,
		name:      name,
		createdAt: time.Now(),
	}
}

// ID returns the user ID
func (u *User) ID() UserID {
	return u.id
}

// Name returns the user name
func (u *User) Name() UserName {
	return u.name
}

// CreatedAt returns the creation timestamp
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// ChangeName changes the user's name
func (u *User) ChangeName(name UserName) {
	u.name = name
}
