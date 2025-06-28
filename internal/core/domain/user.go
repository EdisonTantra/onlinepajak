package domain

import "time"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username    string     `json:"username"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number"`
	Password    string     `json:"password"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
}

type User struct {
	ID          string      `json:"id"`
	ExternalID  string      `json:"external_id"`
	Username    string      `json:"username"`
	Email       string      `json:"email"`
	PhoneNumber string      `json:"phone_number"`
	Password    string      `json:"password"`
	Profile     UserProfile `json:"profile"`
	CreatedAt   *time.Time  `json:"created_at"`
	UpdatedAt   *time.Time  `json:"updated_at"`
}

type UserProfile struct {
	UserID      int        `json:"userId"`
	ID          int        `json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Gender      *string    `json:"gender"`
	Age         *int       `json:"age"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type RequestGetUser struct {
	ExternalID string `json:"external_id"`
}
