package port

import (
	"context"
	"github.com/EdisonTantra/lemonPajak/internal/core/domain"
)

type DJPClient interface {
	EFakturValidation(ctx context.Context, approvalCode string) (*domain.EFakturDJPResponse, error)
}
