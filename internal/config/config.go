package config

import (
	"github.com/go-funcards/envconfig"
	"github.com/go-funcards/grpc-server"
	"github.com/go-funcards/logger"
	"github.com/go-funcards/mongodb"
	"github.com/go-funcards/validate"
	"sync"
)

type Config struct {
	Debug   bool           `yaml:"debug" env:"DEBUG_MODE" env-default:"false"`
	Log     logger.Config  `yaml:"log" env-prefix:"LOG_"`
	MongoDB mongodb.Config `yaml:"mongodb" env-prefix:"MONGODB_"`
	Server  struct {
		Listen grpcserver.Config `yaml:"listen" env-prefix:"LISTEN_"`
	} `yaml:"server" env-prefix:"SERVER_"`
	Rules validate.TypeRules `yaml:"rules"`
}

var (
	cfg  Config
	once sync.Once
)

func GetConfig(path string) (Config, error) {
	var err error
	once.Do(func() {
		err = envconfig.ReadConfig(path, &cfg)
	})
	return cfg, err
}
