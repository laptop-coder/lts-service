package database

import (
	. "backend/config"
	. "backend/logger"
	glide "github.com/valkey-io/valkey-glide/go/v2"
	"github.com/valkey-io/valkey-glide/go/v2/config"
)

func initValkey(host string, port int) *glide.Client {
	config := config.NewClientConfiguration().
		WithAddress(&config.NodeAddress{Host: host, Port: port})

	client, err := glide.NewClient(config)
	if err != nil {
		Logger.Error("Error. Can't create Valkey client: " + err.Error())
		return nil
	} else {
		Logger.Info("Valkey client created")
	}

	return client
}

var Valkey = initValkey(Cfg.Valkey.Host, Cfg.Valkey.Port)
