package sqlstorage

import (
	"context"
	"fmt"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Storage struct { // TODO
	pool   *pgxpool.Pool
	config *config.Config
}

func New(cnf *config.Config) *Storage {
	return &Storage{config: cnf}
}

func (s *Storage) Connect(ctx context.Context) error { // TODO add config
	host := s.config.Database.Host // change
	port := s.config.Database.Port
	user := s.config.Database.Username
	password := s.config.Database.Password
	nameDB := s.config.Database.Name
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, nameDB)
	// url := "postgres://username:password@localhost:5432/database_name" //TODO change
	conn, err := pgxpool.ParseConfig(url)
	if err != nil {
		err = fmt.Errorf("failed to parse pg config: %w", err)
		return err
	}

	// Config
	conn.MaxConns = int32(5)
	conn.MinConns = int32(1)
	conn.HealthCheckPeriod = 1 * time.Minute
	conn.MaxConnLifetime = 24 * time.Hour
	conn.MaxConnIdleTime = 30 * time.Minute
	conn.ConnConfig.ConnectTimeout = 1 * time.Second

	// Create pool
	_, err = pgxpool.ConnectConfig(ctx, conn) // pool
	if err != nil {
		err = fmt.Errorf("failed to connect config: %w", err)
		return err
	}
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.pool.Close()
	return nil
}
