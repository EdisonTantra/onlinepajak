package port

import (
	"context"
	"github.com/EdisonTantra/lemonPajak/internal/core/domain"
)

type UserService interface {
	Get(ctx context.Context, param *domain.RequestGetUser) (*domain.User, error)
	GetExternalUser(ctx context.Context, param *domain.RequestGetUser) (*domain.User, error)
	Patch(id string, data *domain.User) (*domain.User, error)
}

type ApplicationService interface {
	EFakturValidation(ctx context.Context, req *domain.EFakturValidationRequest) (*domain.EFakturValidationResponse, error)
}
