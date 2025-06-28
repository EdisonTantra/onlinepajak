package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EdisonTantra/lemonPajak/internal/core/port"
	storeUserPsql "github.com/EdisonTantra/lemonPajak/internal/repository/postgres/user"
	"github.com/jmoiron/sqlx"
)

const (
	driverPostgres = "postgres"

	sslModeDisabled = "disable"
	sslModeEnabled  = "enable"
)

var _ port.AppsRepository = (*Repository)(nil)

type Repository struct {
	lemonDB  *sql.DB
	lemonDBx *sqlx.DB
	Store    *Store
}

type Store struct {
	User port.UserStore
	App  port.ApplicationStore
}

type NewRepoOptions struct {
	Username    string
	Password    string
	Host        string
	Port        int
	Database    string
	SSLMode     bool
	MaxIdleConn int
	MaxOpenConn int
	MaxIdleTime time.Duration
}

func New(ctx context.Context, opts NewRepoOptions) (*Repository, error) {
	sslModeStr := sslModeDisabled
	if opts.SSLMode {
		sslModeStr = sslModeEnabled
	}

	addr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Port,
		opts.Database,
		sslModeStr,
	)

	var pg, err = sqlx.ConnectContext(ctx, driverPostgres, addr)
	if err != nil {
		return nil, err
	}
	pg.SetMaxIdleConns(opts.MaxIdleConn)
	pg.SetMaxOpenConns(opts.MaxOpenConn)
	pg.SetConnMaxIdleTime(opts.MaxIdleTime)

	return &Repository{
		lemonDB:  pg.DB,
		lemonDBx: pg,
	}, nil
}

func (r *Repository) RegisterStore() {
	r.Store = &Store{
		User: storeUserPsql.New(r.lemonDBx),
		App:  storeUserPsql.New(r.lemonDBx),
	}
}

func (r *Repository) Close() error {
	return r.lemonDBx.Close()
}

func (r *Repository) GetUserStore() port.UserStore {
	return r.Store.User
}

func (r *Repository) GetAppStore() port.ApplicationStore {
	return r.Store.App
}
