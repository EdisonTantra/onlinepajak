package port

import (
	"context"
	"github.com/EdisonTantra/lemonPajak/internal/core/domain"
)

type UserStore interface {
	CreateUser(data *domain.User) (*domain.User, error)
	Login(phone string, password string) (*domain.User, error)
	GetUserByID(ID string) (*domain.User, error)
	PatchUserByID(id string, data *domain.User) (*domain.User, error)
	GetActiveUserByExternalID(ctx context.Context, externalID string) (*domain.User, error)
}

type ApplicationStore interface{}
