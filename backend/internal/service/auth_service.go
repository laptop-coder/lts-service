package service

import (
	"backend/pkg/apperrors"
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
		if errors.Is(err, apperrors.ErrUserNotFound) {
			s.log.Warn("login failed: invalid credentials: user with this email does not exist")
			return nil, nil, fmt.Errorf("invalid credentials: %w", apperrors.ErrInvalidCredentials)
		}
		return nil, nil, err
	}
	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.log.Warn("login failed: invalid password")
		return nil, nil, fmt.Errorf("invalid credentials: %w", apperrors.ErrInvalidCredentials)
	}
	// Generate tokens
	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}
	s.log.Info("user logged in successfully")
	// Return response
	return tokens, UserToDTO(user), nil
}

//TODO: rename to RefreshTokens?
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	s.log.Info("starting tokens refresh...")
	// Parse (and validate) refresh token
	claims, err := s.ParseToken(refreshToken)
	if err != nil || claims == nil { // TODO: check in the whole code like here
		return nil, fmt.Errorf("invalid refresh token: %s: %w", err.Error(), apperrors.ErrInvalidToken)
	}
	// Check token type
	if claims.TokenType != RefreshToken {
		return nil, fmt.Errorf("invalid token type: expected refresh: %w", apperrors.ErrInvalidToken)
	}
	// Check if token is in black list (i.e. revoked)
	revoked, err := s.jwtRepo.IsRevoked(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to check if token was revoked: %w", err)
	}
	if revoked {
		return nil, fmt.Errorf("token was revoked: %w", apperrors.ErrTokenRevoked)
	}
	// Find user
	user, err := s.userRepo.FindByID(ctx, &claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found by id (%s): %s: %w", claims.UserID.String(), err.Error(), apperrors.ErrInvalidToken)
	}
	// Generate new token pair
	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}
	// Revoke old refresh token
	if err := s.jwtRepo.Revoke(ctx, refreshToken, time.Until(claims.RegisteredClaims.ExpiresAt.Time)); err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}
	// Return response
	s.log.Info("tokens refreshed successfully")
	return tokens, nil
}

func (s *authService) RevokeToken(ctx context.Context, token string) error {
	s.log.Info("starting token revoke...")
	// Check if token was already revoked
	revoked, err := s.jwtRepo.IsRevoked(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to check if token was already revoked: %w", err)
	}
	if revoked {
		return fmt.Errorf("token was already revoked: %w", apperrors.ErrTokenRevoked)
	}
	// Parse token
	parsedToken, err := s.ParseToken(token)
	if err != nil || parsedToken == nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}
	// Check if token already expired
	if parsedToken.ExpiresAt.Time.Before(time.Now()) {
		return fmt.Errorf("token already expired: %w", apperrors.ErrTokenExpired)
	}
	// Revoke
	if err := s.jwtRepo.Revoke(ctx, token, time.Until(parsedToken.RegisteredClaims.ExpiresAt.Time)); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	s.log.Info("token revoked successfully")
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

func (s *authService) ParseToken(tokenString string) (*TokenClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v: %w", token.Header["alg"], apperrors.ErrInvalidToken)
		}
		return s.config.JWTSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		// Check token issuer
		if claims.Issuer != s.config.TokenIssuer {
			return nil, fmt.Errorf("invalid token issuer: %w", apperrors.ErrInvalidToken)
		}
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token (failed to parse token): %w", apperrors.ErrInvalidToken)
}
