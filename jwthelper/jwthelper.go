package jwthelper

// JWTService defines the structure for the JWT service
type JWTService struct {
	secretKey []byte
}

// NewJWTService creates a new instance of the JWTService with a given secret key
func NewJWTService(secretKey []byte) *JWTService {
	return &JWTService{
		secretKey: secretKey,
	}
}
