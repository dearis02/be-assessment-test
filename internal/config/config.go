package config

import (
	"fmt"
	"os"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string         `yaml:"environment"`
	Server      Server         `yaml:"server"`
	Database    PostgresConfig `yaml:"database"`
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Server, validation.Required),
		validation.Field(&c.Database, validation.Required),
	)
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (c *Server) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (s Server) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Host, validation.Required),
		validation.Field(&s.Port, validation.Required),
	)
}

type PostgresConfig struct {
	ConString   string `yaml:"connection_string"`
	MaxIdleCons int    `yaml:"max_idle_cons"`
	MaxOpenCons int    `yaml:"max_open_cons"`
}

func (p PostgresConfig) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ConString, validation.Required),
		validation.Field(&p.MaxIdleCons, validation.Required),
		validation.Field(&p.MaxOpenCons, validation.Required),
	)
}

func New() *Config {
	cfg := &Config{}

	cfgData, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read config file")
	}

	err = yaml.Unmarshal(cfgData, &cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal config file")
	}

	err = cfg.Validate()
	if err != nil {
		log.Fatal().Err(err).Msg("config validation failed")
	}

	return cfg
}
