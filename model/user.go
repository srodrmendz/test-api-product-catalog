package model

import "time"

// Represents user structure
type User struct {
	ID            string     `json:"id" bson:"_id"`
	UserName      string     `json:"username" bson:"username"`
	Email         string     `json:"email" bson:"email"`
	EmailVerified bool       `json:"email_verified" bson:"email_verified"`
	Phone         *string    `json:"phone,omitempty" bson:"phone,omitempty"`
	NickName      string     `json:"nickname" bson:"nickname"`
	Picture       *string    `json:"picture,omitempty" bson:"picture,omitempty"`
	Blocked       bool       `json:"blocked" bson:"blocked"`
	Password      string     `json:"password" bson:"password"`
	CreatedAt     time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" bson:"updated_at"`
	LastLogin     *time.Time `json:"last_login,omitempty" bson:"last_login,omitempty"`
}

// Represents user structure used for response purpose
type UserResponse struct {
	ID        string    `json:"id"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone,omitempty"`
	NickName  string    `json:"nickname"`
	Picture   *string   `json:"picture,omitempty"`
	Blocked   bool      `json:"blocked"`
	CreatedAt time.Time `json:"created_at"`
}

// Map User data model into UserResponse
func MapUserToResponse(u User) *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		UserName:  u.UserName,
		Email:     u.Email,
		Phone:     u.Phone,
		NickName:  u.NickName,
		Picture:   u.Picture,
		Blocked:   u.Blocked,
		CreatedAt: u.CreatedAt,
	}
}
