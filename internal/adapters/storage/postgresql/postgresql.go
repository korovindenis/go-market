package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/korovindenis/go-market/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

const (
	UniqueViolation = "duplicate key value violates unique constraint"
)

type Storage struct {
	db *sql.DB
	config
}

type config interface {
	GetStorageConnectionString() string
	GetStorageSalt() string
}

func New(config config) (*Storage, error) {
	db, err := sql.Open("pgx", config.GetStorageConnectionString())
	if err != nil {
		return nil, err
	}

	storage := &Storage{
		db,
		config,
	}

	if err := storage.runMigrations(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) runMigrations() error {
	return goose.Run("up", s.db, "deployments/db/migrations")
}

func (s *Storage) UserRegister(ctx context.Context, user entity.User) error {
	// add salt to password
	password, err := hashPassword(user.Password, s.config.GetStorageSalt())
	if err != nil {
		return err
	}

	// add user or return ErrUserLoginNotUnique
	if _, err := s.db.ExecContext(ctx, "INSERT INTO users (login, password) VALUES ($1, $2)", user.Login, password); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return entity.ErrUserLoginNotUnique
		}
		return err
	}

	return nil
}

// auth user or return ErrUserLoginUnauthorized
func (s *Storage) UserLogin(ctx context.Context, user entity.User) error {
	var userPassword string
	if err := s.db.QueryRowContext(ctx, "SELECT password FROM users WHERE login = $1", user.Login).Scan(&userPassword); err != nil {
		return entity.ErrUserLoginUnauthorized
	}

	// add salt to password
	user.Password = fmt.Sprintf("%s%s", user.Password, s.config.GetStorageSalt())

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(user.Password)); err != nil {
		return entity.ErrUserLoginUnauthorized
	}

	return nil
}

func (s *Storage) GetUser(ctx context.Context, userFromReq entity.User) (entity.User, error) {
	var userFromStorageID uint64
	if err := s.db.QueryRowContext(ctx, "SELECT id FROM users WHERE login = $1", userFromReq.Login).Scan(&userFromStorageID); err != nil {
		return entity.User{}, err
	}
	userFromStorage := entity.User{
		ID: userFromStorageID,
	}

	return userFromStorage, nil
}

func hashPassword(password, salt string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// orders
func (s *Storage) AddOrder(ctx context.Context, order entity.Order, user entity.User) error {
	var existingOrderUser sql.NullInt64
	err := s.db.QueryRowContext(ctx, "SELECT user_id FROM orders WHERE number = $1", order.Number).Scan(&existingOrderUser)
	if errors.Is(err, sql.ErrNoRows) {
		if _, err := s.db.ExecContext(ctx, "INSERT INTO orders (number, user_id, status) VALUES ($1, $2, 'NEW')", order.Number, user.ID); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		// An order with this number already exists
		if existingOrderUser.Int64 == int64(user.ID) {
			return entity.ErrOrderAlreadyUploaded
		} else {
			return entity.ErrOrderAlreadyUploadedAnotherUser
		}
	}

	return nil
}
