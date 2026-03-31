package domain

type MarketService interface {
	GetCurrentPrices() []MarketTick
	Subscribe() chan []MarketTick
}

type AuthService interface {
	Register(req AuthRequest) (*AuthResponse, error)
	Login(req AuthRequest) (*AuthResponse, error)
	ValidateToken(token string) (*Claims, error)
}

type UserRepository interface {
	Save(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	Exists(email string) bool
}

type TokenProvider interface {
	GenerateToken(user *User) (string, error)
	ValidateToken(token string) (*Claims, error)
}
