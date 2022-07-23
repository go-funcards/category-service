package config

import (
	"github.com/go-funcards/validate"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"sync"
)

type Config struct {
	MongoDB struct {
		URI string `yaml:"uri" env:"URI" env-required:"true"`
	} `yaml:"mongodb" env-prefix:"MONGODB_"`
	GRPC struct {
		Addr string `yaml:"address" env:"ADDR" env-default:":80"`
	} `yaml:"grpc" env-prefix:"GRPC_"`
	Validation struct {
		Rules validate.TypeRules `yaml:"rules" env:"RULES"`
	} `yaml:"validation" env-prefix:"VALIDATION_"`
}

var (
	cfg  Config
	once sync.Once
)

func GetConfig(path string, log logrus.FieldLogger) Config {
	once.Do(func() {
		log.Debugf("read config from path %s", path)

		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			log.Fatal(err)
		}
	})
	return cfg
}
