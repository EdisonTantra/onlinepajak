package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	Name        string `valid:"required,alphanum"`     // Required, alphanumeric
	Email       string `valid:"email,required"`        // Required, valid email
	Username    string `valid:"alphanum,required"`     // Required, alphanumeric
	Description string `valid:"stringlength(0|20)"`    // Optional, max 20 chars
	Age         int    `valid:"range(18|60),required"` // Required, in range
	Bio         string `valid:"optional"`
	Handle      string `valid:"matches(^[a-zA-Z0-9_-]+$)"`
}

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name: "valid input",
			input: User{
				Name:        "Alice",
				Email:       "alice@example.com",
				Username:    "alice",
				Description: "alice",
				Age:         20,
				Bio:         "alice",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			input: User{
				Name:        "",
				Email:       "alice@example.com",
				Username:    "bob",
				Description: "bob",
				Age:         20,
				Bio:         "bob",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			input: User{
				Name:        "Alice",
				Email:       "not-an-email",
				Username:    "alice",
				Description: "alice",
				Age:         20,
				Bio:         "",
			},
			wantErr: true,
		},
		{
			name: "empty struct",
			input: User{
				Name:        "",
				Email:       "",
				Username:    "",
				Description: "",
				Age:         0,
				Bio:         "",
			},
			wantErr: true,
		},
		{
			name: "valid input",
			input: User{
				Name:        "Alice123",
				Email:       "alice@example.com",
				Username:    "Alice123",
				Description: "This is fine",
				Age:         25,
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			input: User{
				Name:        "Bob123",
				Email:       "not-an-email",
				Username:    "Bob123",
				Description: "Good one",
				Age:         30,
			},
			wantErr: true,
		},
		{
			name: "age too low",
			input: User{
				Name:        "Charlie",
				Email:       "charlie@example.com",
				Username:    "Charlie1",
				Description: "Valid desc",
				Age:         17,
			},
			wantErr: true,
		},
		{
			name: "age too high",
			input: User{
				Name:        "David",
				Email:       "david@example.com",
				Username:    "David1",
				Description: "Valid desc",
				Age:         100,
			},
			wantErr: true,
		},
		{
			name: "description too long",
			input: User{
				Name:        "Eve123",
				Email:       "eve@example.com",
				Username:    "Eve123",
				Description: "This description is way too long for the field",
				Age:         28,
			},
			wantErr: true,
		},
		{
			name: "username with symbol",
			input: User{
				Name:        "Frank",
				Email:       "frank@example.com",
				Username:    "Frank_123", // _ is not alphanum
				Description: "Valid desc",
				Age:         35,
			},
			wantErr: true,
		},
		{
			name: "name with symbol",
			input: User{
				Name:        "John@Doe", // @ is invalid
				Email:       "john@example.com",
				Username:    "John123",
				Description: "Cool guy",
				Age:         22,
			},
			wantErr: true,
		},
		{
			name: "valid with empty optional bio",
			input: User{
				Name:        "Alice123",
				Email:       "alice@example.com",
				Username:    "Alice123",
				Description: "Short desc",
				Age:         25,
				Bio:         "", // optional, should be valid
			},
			wantErr: false,
		},
		{
			name: "valid with filled optional bio",
			input: User{
				Name:        "Bob",
				Email:       "bob@example.com",
				Username:    "Bob321",
				Description: "Some info",
				Age:         30,
				Bio:         "I love coding and music.", // valid even if filled
			},
			wantErr: false,
		},
		{
			name: "invalid: other field error but optional bio is fine",
			input: User{
				Name:        "", // required
				Email:       "bob@example.com",
				Username:    "Bob321",
				Description: "Some info",
				Age:         30,
				Bio:         "", // still okay
			},
			wantErr: true,
		},

		// regex
		{
			name: "valid handle with dash and underscore",
			input: User{
				Name:        "TestUser",
				Email:       "test@example.com",
				Username:    "TestUser1",
				Description: "Short text",
				Age:         25,
				Bio:         "",
				Handle:      "my-handle_name123", // valid
			},
			wantErr: false,
		},
		{
			name: "invalid handle with symbol",
			input: User{
				Name:        "TestUser",
				Email:       "test@example.com",
				Username:    "TestUser1",
				Description: "Short text",
				Age:         25,
				Bio:         "",
				Handle:      "invalid@handle!", // contains @ and !
			},
			wantErr: true,
		},
		{
			name: "invalid handle with space",
			input: User{
				Name:        "TestUser",
				Email:       "test@example.com",
				Username:    "TestUser1",
				Description: "Short text",
				Age:         25,
				Bio:         "",
				Handle:      "handle with space", // contains space
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.input)
			if tt.wantErr {
				assert.Error(t, err, "expected error but got none")
			} else {
				assert.NoError(t, err, "expected no error but got one")
			}
		})
	}
}
