package config

import (
	"io/ioutil"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiteshwadhwani/go-rest/pkg/log"
	env "github.com/qiangxue/go-env"
	"gopkg.in/yaml.v3"
)

const (
	defaultServerPort         = 8080
	defaultJWTExpirationHours = 72
)

type Config struct {
	ServerPort         int    `yaml:"server_port" env:"SERVER_PORT"`
	JWTSecret          string `yaml:"jwt_secret" env:"JWT_SECRET"`
	JWTExpirationHours int    `yaml:"jwt_expiration_hours" env:"JWT_EXPIRATION_HOURS"`
	DSN                string `yaml:"dsn" env:"DSN"`
}

func (c Config) ValidateMyStruct() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.DSN, validation.Required),
		validation.Field(&c.JWTSecret, validation.Required))
}

func Load(fileName string, logger log.Logger) (*Config, error) {
	// default config
	c := Config{
		ServerPort:         defaultServerPort,
		JWTExpirationHours: defaultJWTExpirationHours,
	}

	bytes, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	if err := env.New("APP_", logger.Infof).Load(&c); err != nil {
		return nil, err
	}

	if err := c.ValidateMyStruct(); err != nil {
		return nil, err
	}

	return &c, nil
}
