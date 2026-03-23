package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type InviteTokenClaims struct {
	jwt.RegisteredClaims
	RoleIDs []uint8 `json:"roleIds"`
}

type InviteService interface {
	CreateToken(ctx context.Context, roleIDs []uint8) (*string, error)
	GetRoles(ctx context.Context, tokenString string) ([]RoleResponseDTO, error)
	RevokeToken(ctx context.Context, tokenString string) error
}

type inviteService struct {
	jwtRepo  repository.JWTRepository
	roleRepo repository.RoleRepository
	db       *gorm.DB
	config   InviteServiceConfig
	log      logger.Logger
}

func NewInviteService(
	jwtRepo repository.JWTRepository,
	roleRepo repository.RoleRepository,
	db *gorm.DB,
	config InviteServiceConfig,
	log logger.Logger,
) InviteService {
	return &inviteService{
		jwtRepo:  jwtRepo,
		roleRepo: roleRepo,
		db:       db,
		config:   config,
		log:      log,
	}
}

func (s *inviteService) CreateToken(ctx context.Context, roleIDs []uint8) (*string, error) {
	s.log.Info("Starting create invite token")
	token, err := s.generateToken(ctx, roleIDs)
	if err != nil || token == nil {
		s.log.Error("Failed to create invite token", "error", err.Error())
		return nil, fmt.Errorf("failed to generate invite token: %w", err)
	}
	s.log.Info("Invite token created successfully")
	return token, nil
}

func (s *inviteService) generateToken(ctx context.Context, roleIDs []uint8) (*string, error) {
	// Check roles existence
	var count int64
	if err := s.db.WithContext(ctx).
		Model(&model.Role{}).
		Where("id IN (?)", roleIDs).
		Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to check roles existence: %w", err)
	}
	if int(count) != len(roleIDs) {
		return nil, fmt.Errorf("some roles were not found by IDs")
	}
	// Assemble claims
	claims := InviteTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: s.config.TokenIssuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
		RoleIDs: roleIDs,
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.config.JWTSecret)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (s *inviteService) parseToken(tokenString string) (*InviteTokenClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &InviteTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.config.JWTSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if claims, ok := token.Claims.(*InviteTokenClaims); ok && token.Valid {
		// Check token issuer
		if claims.Issuer != s.config.TokenIssuer {
			return nil, fmt.Errorf("invalid invite token issuer")
		}
		return claims, nil
	}
	return nil, fmt.Errorf("invalid invite token (failed to parse token)")
}

func (s *inviteService) GetRoles(ctx context.Context, tokenString string) ([]RoleResponseDTO, error) {
	// Check if token was revoked
	revoked, err := s.jwtRepo.IsRevoked(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to check if invite token was revoked")
	}
	if revoked {
		return nil, fmt.Errorf("invite token was revoked")
	}
	// Parse token
	claims, err := s.parseToken(tokenString)
	if err != nil || claims == nil {
		return nil, fmt.Errorf("failed to parse invite token")
	}
	// Get roleIDs
	roleIDs := claims.RoleIDs
	if len(roleIDs) == 0 {
		return nil, fmt.Errorf("list of the role IDs cannot be empty")
	}
	// Fetch roles
	roles, err := s.roleRepo.FindByIDs(ctx, roleIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles by IDs")
	}
	// Collect DTOs
	dtos := make([]RoleResponseDTO, len(roles))
	for i, role := range roles {
		dtos[i] = *RoleToDTO(&role)
	}
	return dtos, nil
}

func (s *inviteService) RevokeToken(ctx context.Context, tokenString string) error {
	s.log.Info("Starting invite token revoke")
	// Check if token was already revoked
	revoked, err := s.jwtRepo.IsRevoked(ctx, tokenString)
	if err != nil {
		return fmt.Errorf("failed to check if invite token was already revoked")
	}
	if revoked {
		return fmt.Errorf("invite token was already revoked")
	}
	// Parse token
	parsedToken, err := s.parseToken(tokenString)
	if err != nil || parsedToken == nil {
		return fmt.Errorf("failed to parse invite token: %w", err)
	}
	// Check if token already expired
	if parsedToken.ExpiresAt.Time.Before(time.Now()) {
		return fmt.Errorf("invite token already expired")
	}
	// Revoke token
	if err := s.jwtRepo.Revoke(ctx, tokenString, time.Until(parsedToken.RegisteredClaims.ExpiresAt.Time)); err != nil {
		return fmt.Errorf("failed to revoke invite token: %w", err)
	}
	s.log.Info("Invite token revoked successfully")
	return nil
}
