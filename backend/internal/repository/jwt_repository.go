package repository

import (
	"backend/pkg/logger"
	"context"
	"fmt"
	"github.com/valkey-io/valkey-go"
	"time"
)

type JWTRepository interface {
	Revoke(ctx context.Context, token string, ttl time.Duration) error
	IsRevoked(ctx context.Context, token string) (bool, error)
}

type jwtRepository struct {
	client valkey.Client
	log    logger.Logger
}

func NewJWTRepository(jwtClient valkey.Client, log logger.Logger) JWTRepository {
	if jwtClient == nil {
		log.Error("Valkey client is nil")
		panic("Valkey client is nil")
	}
	return &jwtRepository{client: jwtClient, log: log}
}

func (r *jwtRepository) Revoke(ctx context.Context, tokenString string, ttl time.Duration) error {
	err := r.client.Do(ctx, r.client.B().
		Set().
		Key(tokenString).
		Value("1").
		ExSeconds(int64(ttl.Seconds())).
		Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to revoke JWT: %w", err)
	}
	return nil
}

func (r *jwtRepository) IsRevoked(ctx context.Context, token string) (bool, error) {
	key := ""
	n, err := r.client.Do(ctx, r.client.B().
		Exists().
		Key(key).
		Build(),
	).AsInt64()
	if err != nil {
		return false, fmt.Errorf("failed to check token revoke status")
	}
	if n > 1 {
		return false, fmt.Errorf("tokens are not unique (there are %d similar tokens)", n)
	}
	return n == 1, nil
}
