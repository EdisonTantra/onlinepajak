package user

import "time"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID          string `json:"id"`
	AccessToken string `json:"access_token"`
}

type InternalVerifyRequest struct {
	ID          string `json:"user_id"`
	AccessToken string `json:"access_token"`
}

type InternalVerifyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type RegisterRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	UserID  string `json:"user_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type RequestUserDetail struct {
	ExternalID string `json:"external_id"`
}

type ResponseUserDetail struct {
	ID          string     `json:"id"`
	ExternalID  string     `json:"external_id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Gender      *string    `json:"gender"`
	Age         *int       `json:"age"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	//Password    string      `json:"password"`
}
