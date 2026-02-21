package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID    uuid.UUID `json:"userID"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
	TokenType TokenType `json:"tokenType"`
}

type AuthService interface {
	Login(ctx context.Context, email string, password string) (*TokenResponse, *UserResponseDTO, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	ParseToken(tokenString string) (*TokenClaims, error)
	RevokeToken(ctx context.Context, tokenString string) error
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type authService struct {
	userRepo repository.UserRepository
	jwtRepo  repository.JWTRepository
	db       *gorm.DB
	config   AuthServiceConfig
	log      logger.Logger
}

func NewAuthService(
	userRepo repository.UserRepository,
	jwtRepo repository.JWTRepository,
	db *gorm.DB,
	config AuthServiceConfig,
	log logger.Logger,
) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtRepo:  jwtRepo,
		db:       db,
		config:   config,
		log:      log,
	}
}

func (s *authService) Login(ctx context.Context, email string, password string) (*TokenResponse, *UserResponseDTO, error) {
	s.log.Info("Starting user login")
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, &email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("invalid credentials (user with this email does not exist)")
		}
		return nil, nil, fmt.Errorf("failed to found user by email: %w", err)
	}
	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.log.Warn("Login failed: invalid password")
		return nil, nil, fmt.Errorf("invalid credentials (passwords do not match)")
	}
	// Generate tokens
	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}
	s.log.Info("User logged in successfully")
	// Return response
	return tokens, UserToDTO(user), nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	s.log.Info("Refreshing access token")
	// Validate refresh token
	claims, err := s.validateToken(refreshToken, RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	// Check if token is in black list (i.e. revoked)
	revoked, err := s.jwtRepo.IsRevoked(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to check if token was revoked")
	}
	if revoked {
		return nil, fmt.Errorf("token was revoked")
	}
	// Find user
	user, err := s.userRepo.FindByID(ctx, &claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found by id (%s): %w", &claims.UserID, err)
	}
	// Generate new token pair
	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}
	// Revoke old refresh token
	parsedToken, err := s.ParseToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse old refresh token: %w", err)
	}
	if err := s.jwtRepo.Revoke(ctx, refreshToken, time.Until(parsedToken.RegisteredClaims.ExpiresAt.Time)); err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}
	// Return response
	s.log.Info("Token refreshed successfully")
	return tokens, nil
}

func (s *authService) RevokeToken(ctx context.Context, token string) error {
	s.log.Info("Revoking token")
	// Parse
	parsedToken, err := s.ParseToken(token)
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}
	// Revoke
	if err := s.jwtRepo.Revoke(ctx, token, time.Until(parsedToken.RegisteredClaims.ExpiresAt.Time)); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	// Return response
	s.log.Info("Token revoked successfully")
	return nil
}

// TODO: add versions to JWT (for revoking all tokens after changing password or logging out from all devices, e.g.)

func (s *authService) generateTokenPair(ctx context.Context, user *model.User) (*TokenResponse, error) {
	// Get user roles
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}
	// Generate access token
	accessToken, err := s.createToken(user.ID, user.Email, roles, AccessToken, s.config.AccessTokenExpiry)
	if err != nil {
		return nil, err
	}
	// Generate refresh token
	refreshToken, err := s.createToken(user.ID, user.Email, roles, RefreshToken, s.config.RefreshTokenExpiry)
	if err != nil {
		return nil, err
	}
	// Return response
	return &TokenResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

func (s *authService) createToken(userID uuid.UUID, email string, roles []string, tokenType TokenType, expiry time.Duration) (*string, error) {
	// Calculate token expiration time
	expiresAt := time.Now().Add(expiry)
	// Token claims
	claims := &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.TokenIssuer,
			Subject:   userID.String(),
			ID:        uuid.New().String(),
		},
		UserID:    userID,
		Email:     email,
		Roles:     roles,
		TokenType: tokenType,
	}
	// Create token and sign it
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}
	// Return response
	return &tokenString, nil
}

func (s *authService) validateToken(tokenString string, expectedType TokenType) (*TokenClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims format")
	}
	// Check token type
	if claims.TokenType != expectedType {
		return nil, fmt.Errorf("invalid token type: expected %s, got %s", expectedType, claims.TokenType)
	}
	// Return response
	return claims, nil
}

func (s *authService) ParseToken(tokenString string) (*TokenClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.config.JWTSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token (failed to parse token)")
}
