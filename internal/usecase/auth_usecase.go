package usecase

import (
	"errors"
	"os"
	"time"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	"github.com/ifs21014-itdel/concurrent-order-processor/pkg/jwt"
	"github.com/ifs21014-itdel/concurrent-order-processor/pkg/totp"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	Update(user *domain.User) error
}

type AuthUsecase struct {
	repo UserRepo
}

func NewAuthUsecase(r UserRepo) *AuthUsecase {
	return &AuthUsecase{repo: r}
}

func (a *AuthUsecase) Register(email, password, name string) (*domain.User, error) {
	existing, _ := a.repo.FindByEmail(email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := &domain.User{
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
	}
	if err := a.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (a *AuthUsecase) Login(email, password, totpCode string) (string, *domain.User, error) {
	u, err := a.repo.FindByEmail(email)
	if err != nil || u == nil {
		return "", nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	// If user enabled TOTP, validate code
	if u.TOTPEnabled {
		ok, err := totp.ValidateCode(totpCode, u.TOTPSecret)
		if err != nil || !ok {
			return "", nil, errors.New("invalid TOTP code")
		}
	}
	// generate JWT token
	secret := os.Getenv("JWT_SECRET")
	hours := 24
	token, err := jwt.GenerateToken(u.ID, secret, time.Duration(hours)*time.Hour)
	if err != nil {
		return "", nil, err
	}
	return token, u, nil
}

// Generate TOTP secret & provisioning URI
func (a *AuthUsecase) GenerateTOTPForUser(userID uint, issuer string) (string, string, error) {
	u, err := a.repo.FindByID(userID)
	if err != nil || u == nil {
		return "", "", errors.New("user not found")
	}
	key, uri, err := totp.GenerateKeyString(issuer, u.Email)
	if err != nil {
		return "", "", err
	}
	u.TOTPSecret = key
	// don't enable yet; require verification step
	if err := a.repo.Update(u); err != nil {
		return "", "", err
	}
	return key, uri, nil
}

func (a *AuthUsecase) VerifyAndEnableTOTP(userID uint, code string) (bool, error) {
	u, err := a.repo.FindByID(userID)
	if err != nil || u == nil {
		return false, errors.New("user not found")
	}
	ok, err := totp.ValidateCode(code, u.TOTPSecret)
	if err != nil {
		return false, err
	}
	if ok {
		u.TOTPEnabled = true
		_ = a.repo.Update(u)
	}
	return ok, nil
}
