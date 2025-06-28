package application

import "github.com/EdisonTantra/lemonPajak/internal/core/port"

var _ port.ApplicationService = (*Service)(nil)

type Service struct {
	userStore port.UserStore
	appStore  port.ApplicationStore
}

func New(userStore port.UserStore, appStore port.ApplicationStore) *Service {
	return &Service{
		userStore: userStore,
		appStore:  appStore,
	}
}

func (s *Service) Something() {
	panic("not implemented")
}
