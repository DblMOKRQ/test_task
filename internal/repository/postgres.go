package postgres

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"

	"testtask/internal/domain"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxExecutor interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type WalletRepo struct {
	exec pgxExecutor
	log  *zap.Logger
}

type unitOfWork struct {
	tx pgx.Tx
	WalletRepo
}

// Это заглушка, не вызывать!
func (u *unitOfWork) Do(ctx context.Context, fn func(uow domain.UnitOfWork) error) error {
	return errors.New("Do() should not be called on a transactional unitOfWork")
}

func (u *unitOfWork) Wallets() domain.WalletRepository {
	return &u.WalletRepo
}

type Store struct {
	pool *pgxpool.Pool
	WalletRepo
	log *zap.Logger
}

func NewStore(ctx context.Context, user string, password string, host string, port string, dbname string, sslmode string, log *zap.Logger) (*Store, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)

	log = log.With(zap.String("dbname", dbname),
		zap.String("host:port", fmt.Sprintf("%s:%s", host, port)),
		zap.String("user", user),
	)

	log.Info("Connecting to PostgreSQL")

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Error("Error parsing connection string", zap.Error(err))
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}
	config.MaxConns = 50
	config.HealthCheckPeriod = 30 * time.Second
	config.MinConns = 2

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Error("Failed connecting to PostgreSQL", zap.Error(err))
		return nil, fmt.Errorf("error connecting to PostgreSQL: %w", err)
	}

	log.Info("Testing database connection")
	if err := db.Ping(ctx); err != nil {
		log.Error("failed pinging PostgreSQL", zap.Error(err))
		return nil, fmt.Errorf("failed pinging PostgreSQL: %w", err)
	}

	log.Info("Successfully connected to PostgreSQL")

	log.Info("Starting database migrations")

	if err := runMigrations(connStr); err != nil {
		log.Error("Failed to run migrations", zap.Error(err))
		return nil, fmt.Errorf("failed to run migration: %w", err)
	}

	log.Info("Successfully migrated database")
	return &Store{
		pool:       db,
		WalletRepo: WalletRepo{exec: db, log: log},
		log:        log.Named("repository"),
	}, nil
}

// Он нужен для операций чтения, не требующих транзакции
func (s *Store) Wallets() domain.WalletRepository {
	return &s.WalletRepo
}

// Do - это главная точка входа для выполнения бизнес-логики в транзакции.
func (s *Store) Do(ctx context.Context, fn func(uow domain.UnitOfWork) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.log.Error("Failed to begin transaction", zap.Error(err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	uow := &unitOfWork{
		tx: tx,
		WalletRepo: WalletRepo{
			exec: tx,
			log:  s.log,
		},
	}

	if err := fn(uow); err != nil {
		s.log.Debug("Transaction function returned error, rolling back", zap.Error(err))
		return fmt.Errorf("transaction function returned error: %w", err)
	}

	s.log.Debug("Committing transaction")
	return tx.Commit(ctx)
}

func (r *Store) Close() {
	r.log.Info("Closing database connection")
	r.pool.Close()
}

func runMigrations(connStr string) error {
	migratePath := os.Getenv("MIGRATE_PATH")
	if migratePath == "" {
		migratePath = "./migrations"
	}
	absPath, err := filepath.Abs(migratePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	absPath = filepath.ToSlash(absPath)
	migrateUrl := fmt.Sprintf("file://%s", absPath)
	m, err := migrate.New(migrateUrl, connStr)
	if err != nil {
		return fmt.Errorf("start migrations error %v", err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("migration up error: %v", err)
	}
	return nil
}
