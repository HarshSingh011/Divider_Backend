package domain

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userRepo      UserRepository
	tokenProvider TokenProvider
}

func NewAuthService(userRepo UserRepository, tokenProvider TokenProvider) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo:      userRepo,
		tokenProvider: tokenProvider,
	}
}

func (as *AuthServiceImpl) Register(req AuthRequest) (*AuthResponse, error) {
	if req.Email == "" || req.Password == "" || req.Username == "" {
		return nil, errors.New("email, username, and password are required")
	}

	// Check if user already exists
	existing, err := as.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("error checking if user exists")
	}

	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &User{
		ID:        generateID(),
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().In(ISTLocation),
	}

	if err := as.userRepo.Save(user); err != nil {
		return nil, errors.New("failed to save user")
	}

	token, err := as.tokenProvider.GenerateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}, nil
}

func (as *AuthServiceImpl) Login(req AuthRequest) (*AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	// Find user by email
	user, err := as.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("database error while finding user")
	}

	// Check if user exists
	if user == nil {
		return nil, errors.New("account does not exist")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := as.tokenProvider.GenerateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}, nil
}

func (as *AuthServiceImpl) ValidateToken(token string) (*Claims, error) {
	return as.tokenProvider.ValidateToken(token)
}

func generateID() string {
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}
