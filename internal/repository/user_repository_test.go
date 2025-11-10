package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	user := &domain.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpwd",
		Name:         "Test User",
		TOTPSecret:   "secret",
		TOTPEnabled:  true,
		Role:         "admin",
	}

	// Mock behavior
	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Email, user.PasswordHash, user.Name, user.TOTPSecret, user.TOTPEnabled, user.Role).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = repo.Create(user)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), user.ID)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "name", "totp_secret", "totp_enabled", "role"}).
		AddRow(1, "test@example.com", "hashedpwd", "Test User", "secret", true, "admin")

	mock.ExpectQuery("SELECT id, email, password_hash, name, totp_secret, totp_enabled, role FROM users WHERE email=\\$1").
		WithArgs("test@example.com").
		WillReturnRows(rows)

	user, err := repo.FindByEmail("test@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "Test User", user.Name)
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewUserRepository(db)

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: "hashedpwd",
		Name:         "Test User Updated",
		TOTPSecret:   "secret",
		TOTPEnabled:  true,
		Role:         "admin",
	}

	mock.ExpectExec("UPDATE users").
		WithArgs(user.Email, user.PasswordHash, user.Name, user.TOTPSecret, user.TOTPEnabled, user.Role, user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err := repo.Update(user)
	assert.NoError(t, err)
}
