// internal/usecase/user_usecase.go
package usecase

import (
	"GonPay_Backend/internal/domain"
	"GonPay_Backend/pkg/validator"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthResponse struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

type UserUseCase struct {
	userRepo  domain.UserRepository
	validator validator.ValidatorInterface
	jwtSecret string
	jwtTTL    int64
}

func NewUserUseCase(
	userRepo domain.UserRepository,
	validator validator.ValidatorInterface,
	jwtSecret string,
	jwtTTL int64,
) *UserUseCase {
	return &UserUseCase{
		userRepo:  userRepo,
		validator: validator,
		jwtSecret: jwtSecret,
		jwtTTL:    jwtTTL,
	}
}

func (u *UserUseCase) Register(username, email, phoneNumber, password string) (*AuthResponse, error) {
	if err := u.validator.ValidateUsername(username); err != nil {
		return nil, err
	}
	if err := u.validator.ValidateEmail(email); err != nil {
		return nil, err
	}
	if err := u.validator.ValidatePhone(phoneNumber); err != nil {
		return nil, err
	}
	if err := u.validator.ValidatePassword(password); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		Email:        email,
		PhoneNumber:  phoneNumber,
		PasswordHash: string(hashedPassword),
		Status:       domain.UserStatusActive,
		Preferences:  "{}",
		Role:         domain.RoleUser,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := u.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (u *UserUseCase) Login(email, password string) (*AuthResponse, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if user.Status != domain.UserStatusActive {
		return nil, errors.New("account is inactive")
	}

	token, err := u.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// Tiếp tục của file user_usecase.go

func (u *UserUseCase) GetUserByID(id int64) (*domain.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *UserUseCase) UpdateUser(user *domain.User) error {
	if err := u.validator.ValidateUsername(user.Username); err != nil {
		return err
	}
	if err := u.validator.ValidateEmail(user.Email); err != nil {
		return err
	}
	if err := u.validator.ValidatePhone(user.PhoneNumber); err != nil {
		return err
	}

	user.UpdatedAt = time.Now()
	return u.userRepo.Update(user)
}

func (u *UserUseCase) ChangePassword(userID int64, oldPassword, newPassword string) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return domain.ErrInvalidCredentials
	}

	// Validate new password
	if err := u.validator.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return u.userRepo.Update(user)
}

func (u *UserUseCase) generateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role, // Thêm role vào claims
		"exp":     time.Now().Add(time.Hour * time.Duration(u.jwtTTL)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtSecret))
}
