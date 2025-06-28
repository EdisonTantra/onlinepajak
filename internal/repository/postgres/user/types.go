package user

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          int64      `db:"id"`
	ExternalID  uuid.UUID  `db:"external_id" sql:",type:uuid"`
	Email       string     `db:"email"`
	PhoneNumber string     `db:"phone_number"`
	Username    string     `db:"username"`
	Password    string     `db:"password"`
	LoginCount  int64      `db:"login_count"`
	IsPremium   bool       `db:"is_premium"`
	IsVerified  bool       `db:"is_verified"`
	IsActive    bool       `db:"is_active"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

type UserProfile struct {
	ID          big.Int    `db:"id"`
	UserID      big.Int    `db:"user_id"`
	FirstName   string     `db:"first_name"`
	LastName    string     `db:"last_name"`
	Gender      string     `db:"gender"`
	Description string     `db:"description"`
	Age         int        `db:"age"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

type UserPatchArg struct {
	ID          uuid.UUID `db:"id" sql:",type:uuid"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	Email       string    `db:"email"`
	PhoneNumber string    `db:"phone_number"`
	Username    string    `db:"username"`
	Password    string    `db:"password"`
}

type UserLoginArg struct {
	PhoneNumber string `db:"phone_number"`
	Password    string `db:"password"`
}
