package config

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var yamlData = `
logger:
  level: "debug"
http:
  host: "localhost"
  port: 3000
  readTimeout: 5s
  writeTimeout: 10s
grpc:
  host: "localhost"
  port: 50051
  connectTimeout: 5s
db:
  dsn: "postgres://user:pass@localhost:5432/db"
  migrationsDir: "./migrations/"
storage: "memory"
`

func TestNewConfig(t *testing.T) {
	r := bytes.NewReader([]byte(yamlData))
	cfg, err := New(r)
	require.NoError(t, err)
	require.Equal(t, "debug", cfg.Logger.Level)
	require.Equal(t, "localhost", cfg.HTTP.Host)
	require.Equal(t, "3000", cfg.HTTP.Port)
	require.Equal(t, 5*time.Second, cfg.HTTP.ReadTimeout)
	require.Equal(t, 10*time.Second, cfg.HTTP.WriteTimeout)
	require.Equal(t, "localhost", cfg.GRPC.Host)
	require.Equal(t, "50051", cfg.GRPC.Port)
	require.Equal(t, 5*time.Second, cfg.GRPC.ConnectTimeout)
	require.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DB.Dsn)
	require.Equal(t, "./migrations/", cfg.DB.MigrationsDir)
	require.Equal(t, "memory", cfg.Storage)
}

func TestConfigContext(t *testing.T) {
	r := bytes.NewReader([]byte(yamlData))
	cfg, err := New(r)
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}
	ctx := context.Background()
	ctxWithCfg := cfg.WithContext(ctx)
	cfgFromCtx := GetFromContext(ctxWithCfg)
	require.Same(t, cfg, cfgFromCtx)
}

func TestNewConfig_InvalidYAML(t *testing.T) {
	r := bytes.NewReader([]byte("invalid"))
	_, err := New(r)
	require.Error(t, err)
}
