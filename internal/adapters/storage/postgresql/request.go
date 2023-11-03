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
)

const (
	UniqueViolation = "duplicate key value violates unique constraint"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) (*Storage, error) {
	storage := &Storage{
		db,
	}

	return storage, nil
}

func (s *Storage) UserRegister(ctx context.Context, user entity.User) (int64, error) {
	// add user or return ErrUserLoginNotUnique
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	var userID int64
	err = tx.QueryRowContext(ctx, "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id", user.Login, user.Password).Scan(&userID)
	if err != nil {
		tx.Rollback()

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, entity.ErrUserLoginNotUnique
		}
		return 0, err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO balances (user_id) VALUES ($1)", userID)
	if err != nil {

		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return userID, nil
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
	var userFromStorage entity.User
	if err := s.db.QueryRowContext(ctx, "SELECT id FROM users WHERE login = $1 FOR UPDATE", userFromReq.Login).Scan(&userFromStorage.ID); err != nil {
		return userFromStorage, err
	}
	return userFromStorage, nil
}

// orders
func (s *Storage) AddOrder(ctx context.Context, order entity.Order, user entity.User) error {
	var existingOrderUser sql.NullInt64
	err := s.db.QueryRowContext(ctx, "SELECT user_id FROM orders WHERE number = $1 FOR UPDATE", order.Number).Scan(&existingOrderUser)
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
func (s *Storage) GetAllOrders(ctx context.Context, user entity.User) ([]entity.Order, error) {
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
func (s *Storage) GetAllNotProcessedOrders(ctx context.Context) ([]entity.Order, error) {
	var orders []entity.Order
	rows, err := s.db.QueryContext(ctx, "SELECT number,status,accrual,uploaded_at FROM orders WHERE status NOT IN ($1,$2) FOR UPDATE", entity.StatusInvalid, entity.StatusProcessed)
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
func (s *Storage) SetOrderStatusAndAccrual(ctx context.Context, order entity.Order) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var userID int64
	err = tx.QueryRowContext(ctx, "UPDATE orders SET status = $1, accrual = $2 WHERE number = $3 RETURNING user_id", order.Status, order.Accrual, order.Number).Scan(&userID)
	if err != nil {
		tx.Rollback()

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return entity.ErrUserLoginNotUnique
		}
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE balances SET current = current + $1 WHERE user_id = $2", order.Accrual, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// balance
func (s *Storage) GetBalance(ctx context.Context, user entity.User) (entity.Balance, error) {
	var balance entity.Balance
	rows, err := s.db.QueryContext(ctx, "SELECT current,withdrawn FROM balances WHERE user_id = $1 FOR UPDATE", user.ID)
	if err != nil {
		return balance, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&balance.Current, &balance.Withdrawn)
		if err != nil {
			return balance, err
		}
	}

	if err := rows.Err(); err != nil {
		return balance, err
	}

	return balance, nil
}
func (s *Storage) WithdrawBalance(ctx context.Context, balance entity.BalanceUpdate, user entity.User) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentBalance float64
	if err := tx.QueryRowContext(ctx, "SELECT current FROM balances WHERE id = $1 FOR UPDATE", user.ID).Scan(&currentBalance); err != nil {
		return err
	}

	if currentBalance < balance.Sum {
		return entity.ErrInsufficientBalance
	}

	if _, err := tx.ExecContext(ctx, "UPDATE balances SET current = current - $1, withdrawn = withdrawn + $1 WHERE id = $2", balance.Sum, user.ID); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO orders (number, sum, user_id, status) VALUES ($1, $2, $3, $4)", balance.Order, balance.Sum, user.ID, "PROCESSED"); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Withdrawals
func (s *Storage) Withdrawals(ctx context.Context, user entity.User) ([]entity.BalanceUpdate, error) {
	var orders []entity.Order
	rows, err := s.db.QueryContext(ctx, "SELECT number,sum,uploaded_at FROM orders WHERE sum > $1 AND status = $2 AND user_id = $3 ORDER BY id DESC FOR UPDATE", 0, entity.StatusProcessed, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		err := rows.Scan(&order.Number, &order.Sum, &order.UploadedAt)
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

	var balances []entity.BalanceUpdate
	for _, order := range orders {
		balance := entity.BalanceUpdate{
			Order:      order.Number,
			Sum:        order.Sum,
			UploadedAt: order.UploadedAt,
		}
		balances = append(balances, balance)
	}

	return balances, nil
}
