package usersvc

import (
	"context"
	"errors"
	"fmt"
	"github.com/EdisonTantra/lemonPajak/pkg/lib/logat"
	"strings"
	"unicode"

	"github.com/EdisonTantra/lemonPajak/internal/core/cons"
	"github.com/EdisonTantra/lemonPajak/internal/core/domain"
	"github.com/EdisonTantra/lemonPajak/internal/core/port"
	"github.com/EdisonTantra/lemonPajak/pkg/lib/tracer"
)

var _ port.UserService = (*Service)(nil)

type Service struct {
	userStore port.UserStore
}

func New(userStore port.UserStore) *Service {
	return &Service{
		userStore: userStore,
	}
}

func (s *Service) Get(ctx context.Context, param *domain.RequestGetUser) (*domain.User, error) {
	tr := tracer.StartTrace(ctx, "UserService-Get")
	defer tr.Finish()
	ctx = tr.Context()

	if param.ExternalID == "" {
		return nil, errors.New("user ID required")
	}

	data, err := s.userStore.GetActiveUserByExternalID(ctx, param.ExternalID)
	if err != nil {
		logat.GetLogger().Error(ctx, "error GetActiveUserByExternalID", cons.EventLogNameUserDetail, err)
		return nil, err
	}

	return data, nil
}

func (s *Service) GetExternalUser(ctx context.Context, param *domain.RequestGetUser) (*domain.User, error) {
	if param.ExternalID == "" {
		return nil, errors.New("user ID required")
	}

	data, err := s.userStore.GetActiveUserByExternalID(ctx, param.ExternalID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) Patch(id string, data *domain.User) (*domain.User, error) {
	panic("not implemented")
}

func validateUserData(data *domain.User) error {
	var errWrap error
	if data.Profile.FirstName != "" || data.Profile.LastName != "" {
		fullName := fmt.Sprintf("%s %s", data.Profile.FirstName, data.Profile.LastName)
		err := validateFullName(fullName)
		if err != nil {
			errWrap = err
		}
	}

	if data.Password != "" {
		err := validatePassword(data.Password)
		if err != nil {
			if errWrap == nil {
				errWrap = err
			} else {
				errWrap = fmt.Errorf("%w; %w", errWrap, err)
			}
		}
	}

	if data.PhoneNumber != "" {
		err := validatePhoneNumber(data.PhoneNumber)
		if err != nil {
			if errWrap == nil {
				errWrap = err
			} else {
				errWrap = fmt.Errorf("%w; %w", errWrap, err)
			}
		}
	}

	if errWrap != nil {
		return errWrap
	}

	return nil
}

func validateFullName(name string) error {
	if len(name) < cons.MinLengthName || len(name) > cons.MaxLengthName {
		return cons.ErrInvalidNameLength
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < cons.MinLengthPass || len(password) > cons.MaxLengthPass {
		return cons.ErrInvalidPasswordLength
	}

	var number, capital, symbol bool
	for _, r := range password {
		switch {
		case unicode.IsNumber(r):
			number = true
		case unicode.IsUpper(r):
			capital = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			symbol = true
		}
	}

	if !number || !capital || !symbol {
		return cons.ErrInvalidPasswordFormat
	}

	return nil
}

func validatePhoneNumber(phone string) error {
	if !strings.HasPrefix(phone, cons.PrefixPhoneID) {
		return cons.ErrInvalidPhonePrefix
	}

	cleanPhone := strings.Replace(phone, cons.PrefixPhoneID, "0", 1)
	if len(cleanPhone) < cons.MinLengthPhone || len(cleanPhone) > cons.MaxLengthPhone {
		return cons.ErrInvalidPhoneLength
	}

	return nil
}
