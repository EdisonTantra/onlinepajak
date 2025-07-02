package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EdisonTantra/lemonPajak/pkg/lib/logat"
	"github.com/EdisonTantra/lemonPajak/pkg/lib/tracer"

	"github.com/EdisonTantra/lemonPajak/internal/core/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	DBLemon  *sql.DB
	DBLemonx *sqlx.DB
}

func New(dbx *sqlx.DB) *Store {
	return &Store{
		DBLemon:  dbx.DB,
		DBLemonx: dbx,
	}
}

func (s *Store) CreateUser(data *domain.User) (*domain.User, error) {
	panic("not implemented")
}

func (s *Store) Login(phone string, password string) (*domain.User, error) {
	panic("not implemented")
}

func (s *Store) GetUserByID(id string) (*domain.User, error) {
	panic("not implemented")
}

func (s *Store) PatchUserByID(id string, data *domain.User) (*domain.User, error) {
	panic("not implemented")
}

func (s *Store) GetActiveUserByExternalID(ctx context.Context, externalID string) (*domain.User, error) {
	tr := tracer.StartTrace(ctx, "UserStore-GetActiveUserByExternalID")
	defer tr.Finish()

	q := `
		SELECT 
		    ud.id, 
		    ud.external_id,
		    ud.username,
		    ud.email,
		    ud.phone_number,
		    ud.is_premium,
		    ud.is_verified,
		    ud.is_active,
-- 		    up.first_name,
-- 		    up.last_name,
-- 		    up.age,
-- 		    up.gender,
-- 		    up.description,
		    ud.created_at, 
		    ud.updated_at 
		FROM user_customer ud
-- 		JOIN user_profile up on ud.id = up.user_id
		WHERE external_id = :external_id AND is_active = TRUE;
	`

	validID, err := uuid.Parse(externalID)
	if err != nil {
		logat.GetLogger().Error(ctx, "invalid uuid", "GetActiveUserByExternalID", err)
		return nil, err
	}

	arg := User{ExternalID: validID}
	resQuery := User{}

	nstmt, err := s.DBLemonx.PrepareNamed(q)
	if err != nil {
		logat.GetLogger().Error(ctx, "error prepare named", "GetActiveUserByExternalID", err)
		return nil, err
	}

	rows, err := nstmt.Queryx(&arg)
	if err != nil {
		logat.GetLogger().Error(ctx, "error queryx", "GetActiveUserByExternalID", err)
		return nil, err
	}

	for rows.Next() {
		err = rows.StructScan(&resQuery)
		if err != nil {
			return nil, err
		}
	}

	age := 10
	gender := "MALE"

	return &domain.User{
		ID:          fmt.Sprintf("%d", resQuery.ID),
		ExternalID:  resQuery.ExternalID.String(),
		Username:    resQuery.Username,
		Password:    resQuery.Password,
		Email:       resQuery.Email,
		PhoneNumber: resQuery.PhoneNumber,
		Profile: domain.UserProfile{
			FirstName:   "ed",
			LastName:    "son",
			Age:         &age,
			Gender:      &gender,
			Description: "test desc",
		},
		CreatedAt: &resQuery.CreatedAt,
		UpdatedAt: resQuery.UpdatedAt,
	}, nil
}
