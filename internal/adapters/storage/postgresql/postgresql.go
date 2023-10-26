package postgresql

import (
	"context"
	"database/sql"
	"errors"

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
	// add user or return ErrUserLoginNotUnique
	if _, err := s.db.ExecContext(ctx, "INSERT INTO users (login, password) VALUES ($1, $2)", user.Login, user.Password); err != nil {
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

// orders
func (s *Storage) AddOrder(ctx context.Context, order entity.Order, user entity.User) error {
	var existingOrderUser sql.NullInt64
	err := s.db.QueryRowContext(ctx, "SELECT user_id FROM orders WHERE number = $1", order.Number).Scan(&existingOrderUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if _, err := s.db.ExecContext(ctx, "INSERT INTO orders (number, user_id, status) VALUES ($1, $2, $3)", order.Number, user.ID, "NEW"); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	// an order with this number already exists
	if existingOrderUser.Int64 == int64(user.ID) {
		return entity.ErrOrderAlreadyUploaded
	}
	return entity.ErrOrderAlreadyUploadedAnotherUser
}

func (s *Storage) GetOrder(ctx context.Context, user entity.User) ([]entity.Order, error) {
	var orders []entity.Order
	rows, err := s.db.QueryContext(ctx, "SELECT number,status,accrual,uploaded_at FROM orders WHERE user_id = $1 ORDER BY id DESC", user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, entity.ErrNoContent
	}

	return orders, err
}
