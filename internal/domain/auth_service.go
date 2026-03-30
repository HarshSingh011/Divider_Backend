package domain

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AuthServiceImpl implements AuthService interface
type AuthServiceImpl struct {
	userRepo      UserRepository
	tokenProvider TokenProvider
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo UserRepository, tokenProvider TokenProvider) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo:      userRepo,
		tokenProvider: tokenProvider,
	}
}

// Register creates a new user account
func (as *AuthServiceImpl) Register(req AuthRequest) (*AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" || req.Username == "" {
		return nil, errors.New("email, username, and password are required")
	}

	// Check if user already exists
	if as.userRepo.Exists(req.Email) {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	user := &User{
		ID:        generateID(),
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// Save user
	if err := as.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Generate token
	token, err := as.tokenProvider.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}, nil
}

// Login authenticates a user and returns a token
func (as *AuthServiceImpl) Login(req AuthRequest) (*AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	// Find user
	user, err := as.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := as.tokenProvider.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}, nil
}

// ValidateToken checks if a token is valid
func (as *AuthServiceImpl) ValidateToken(token string) (*Claims, error) {
	return as.tokenProvider.ValidateToken(token)
}

// generateID creates a unique user ID
func generateID() string {
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}
