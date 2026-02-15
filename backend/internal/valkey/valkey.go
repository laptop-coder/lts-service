package valkey

import (
	"context"
	"fmt"
	"github.com/valkey-io/valkey-go"
)

type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func Connect(config Config) (valkey.Client, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{fmt.Sprintf("%s:%d", config.Host, config.Port)},
		Password:    config.Password,
		SelectDB:    config.DB,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Valkey client: %w", err)
	}
	return client, nil
}

func Close(client valkey.Client) error {
	if client == nil {
		return nil
	}
	client.Close()
	return nil
}

func HealthCheck(client valkey.Client) error {
	if client == nil {
		return fmt.Errorf("Valkey client is nil")
	}
	ctx := context.Background()
	return client.Do(ctx, client.B().Ping().Build()).Error()
}
