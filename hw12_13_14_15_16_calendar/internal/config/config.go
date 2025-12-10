package config

import (
	"context"
	"errors"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

type key int

const (
	ctxKey key = iota
)

var ErrNoConfigInContext = errors.New("no config found in context")

type Config struct {
	Logger struct {
		Level string `yaml:"level"`
	} `yaml:"logger"`
	HTTP struct {
		Host         string        `yaml:"host"`
		Port         string        `yaml:"port"`
		ReadTimeout  time.Duration `yaml:"readTimeout"`
		WriteTimeout time.Duration `yaml:"writeTimeout"`
	} `yaml:"http"`
	GRPC struct {
		Host           string        `yaml:"host"`
		Port           string        `yaml:"port"`
		ConnectTimeout time.Duration `yaml:"connectTimeout"`
	} `yaml:"grpc"`
	DB struct {
		Dsn           string `yaml:"dsn"`
		MigrationsDir string `yaml:"migrationsDir"`
	} `yaml:"db"`
	Storage string `yaml:"storage"`
}

func New(r io.Reader) (*Config, error) {
	config := &Config{}
	if err := yaml.NewDecoder(r).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, c)
}

func GetFromContext(ctx context.Context) *Config {
	return ctx.Value(ctxKey).(*Config)
}
