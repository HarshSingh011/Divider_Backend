package domain

// MarketService defines the market-related business logic
type MarketService interface {
	GetCurrentPrices() []MarketTick
	Subscribe() chan []MarketTick
}

// AuthService defines the authentication business logic
type AuthService interface {
	Register(req AuthRequest) (*AuthResponse, error)
	Login(req AuthRequest) (*AuthResponse, error)
	ValidateToken(token string) (*Claims, error)
}

// UserRepository defines user storage operations
type UserRepository interface {
	Save(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	Exists(email string) bool
}

// TokenProvider defines JWT token operations
type TokenProvider interface {
	GenerateToken(user *User) (string, error)
	ValidateToken(token string) (*Claims, error)
}
