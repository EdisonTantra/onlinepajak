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
	//TODO create user with profile, preferences, swipe quota

	q := `
		INSERT INTO user_customer (external_id, username, password, email, phone_number) 
		VALUES (:external_id, :username, crypt(:password, gen_salt('bf')), :email, :phone_number)
		RETURNING id, phone_number;
	`

	newUUID, err := uuid.NewV6()
	if err != nil {
		return nil, err
	}

	arg := User{
		ExternalID:  newUUID,
		Username:    data.Username,
		Password:    data.Password,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
	}

	nstmt, err := s.DBLemonx.PrepareNamed(q)
	if err != nil {
		return nil, err
	}

	rows, err := nstmt.Queryx(&arg)
	if err != nil {
		return nil, err
	}

	u := User{}
	for rows.Next() {
		err = rows.StructScan(&u)
		if err != nil {
			return nil, err
		}
	}

	qProfile := `
		INSERT INTO profile (user_id, first_name, last_name) 
		VALUES (:user_id, :first_name, :last_name)
		RETURNING id, first_name, last_name;
	`

	argProfile := UserProfile{
		FirstName:   data.Profile.FirstName,
		LastName:    data.Profile.LastName,
		Gender:      *data.Profile.Gender,
		Age:         *data.Profile.Age,
		Description: data.Profile.Description,
	}

	nstmt2, err := s.DBLemonx.PrepareNamed(qProfile)
	if err != nil {
		return nil, err
	}

	rowsProfile, err := nstmt2.Queryx(&argProfile)
	if err != nil {
		return nil, err
	}

	p := UserProfile{}
	for rowsProfile.Next() {
		err = rowsProfile.StructScan(&p)
		if err != nil {
			return nil, err
		}
	}

	return &domain.User{
		ID:          fmt.Sprintf("%d", u.ID),
		ExternalID:  u.ExternalID.String(),
		Email:       u.Email,
		Username:    u.Username,
		Password:    u.Password,
		PhoneNumber: u.PhoneNumber,
		Profile: domain.UserProfile{
			FirstName: p.FirstName,
			LastName:  p.LastName,
		},
	}, nil
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

// TODO fix implementation
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
		//TODO edison ini ngga bisa pake buat join
		err = rows.StructScan(&resQuery)
		if err != nil {
			return nil, err
		}
	}

	//TODO edison ambil dari resp
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
