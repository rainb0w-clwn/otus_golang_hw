package config

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/creasty/defaults"
	gocfg "github.com/dsbasko/go-cfg"
	"gopkg.in/yaml.v3"
)

type key int

const (
	ctxKey key = iota
)

var ErrNoConfigInContext = errors.New("no config found in context")

type Config struct {
	Logger struct {
		Level string `yaml:"level" env:"LOG_LEVEL" default:"debug"`
	} `yaml:"logger"`
	HTTP struct {
		Host         string        `yaml:"host" env:"HTTP_HOST"`
		Port         string        `yaml:"port" env:"HTTP_PORT"`
		ReadTimeout  time.Duration `default:"5s" yaml:"readTimeout" env:"HTTP_READ_TIMEOUT"`
		WriteTimeout time.Duration `default:"5s" yaml:"writeTimeout" env:"HTTP_WRITE_TIMEOUT"`
	} `yaml:"http"`
	GRPC struct {
		Host           string        `yaml:"host" env:"GRPC_HOST"`
		Port           string        `yaml:"port" env:"GRPC_PORT"`
		ConnectTimeout time.Duration `default:"5s" yaml:"connectTimeout" env:"GRPC_CONNECT_TIMEOUT"`
	} `yaml:"grpc"`
	DB struct {
		Dsn           string `yaml:"dsn" env:"DB_DSN"`
		MigrationsDir string `yaml:"migrationsDir" env:"DB_MIGRATIONS_DIR"`
		Migrate       bool   `yaml:"migrate" env:"DB_MIGRATE"`
	} `yaml:"db"`
	Storage   string `yaml:"storage" env:"STORAGE"`
	Scheduler struct {
		Period          time.Duration `default:"3s" yaml:"period" env:"SCHEDULER_PERIOD"`
		Queue           string        `yaml:"queue" env:"SCHEDULER_QUEUE"`
		RetentionPeriod time.Duration `default:"8760h" yaml:"retentionPeriod" env:"SCHEDULER_RETENTION_PERIOD"`
	} `yaml:"scheduler"`
	RMQ struct {
		Host     string `yaml:"host" env:"RMQ_HOST"`
		Port     string `yaml:"port" env:"RMQ_PORT"`
		Login    string `yaml:"login" env:"RMQ_LOGIN"`
		Password string `yaml:"password" env:"RMQ_PASSWORD"`
	} `yaml:"rmq"`
}

func New(r io.Reader) (*Config, error) {
	config := &Config{}
	if err := yaml.NewDecoder(r).Decode(config); err != nil {
		return nil, err
	}
	if err := gocfg.ReadEnv(config); err != nil {
		return nil, err
	}
	if err := defaults.Set(config); err != nil {
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
