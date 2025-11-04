package repository

import (
	"database/sql"
	"errors"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	Update(user *domain.User) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *domain.User) error {
	query := `INSERT INTO users (email, password_hash, name, totp_secret, totp_enabled, role) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return r.db.QueryRow(query,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.TOTPSecret,
		user.TOTPEnabled,
		user.Role,
	).Scan(&user.ID)
}

func (r *userRepo) FindByEmail(email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, name, totp_secret, totp_enabled, role 
	          FROM users WHERE email=$1`
	var u domain.User
	err := r.db.QueryRow(query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.TOTPSecret, &u.TOTPEnabled, &u.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) FindByID(id uint) (*domain.User, error) {
	query := `SELECT id, email, password_hash, name, totp_secret, totp_enabled, role 
	          FROM users WHERE id=$1`
	var u domain.User
	err := r.db.QueryRow(query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.TOTPSecret, &u.TOTPEnabled, &u.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Update(user *domain.User) error {
	query := `UPDATE users 
	          SET email=$1, password_hash=$2, name=$3, totp_secret=$4, totp_enabled=$5, role=$6 
	          WHERE id=$7`
	res, err := r.db.Exec(query,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.TOTPSecret,
		user.TOTPEnabled,
		user.Role,
		user.ID,
	)
	if err != nil {
		return err
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("user not found")
	}
	return nil
}
