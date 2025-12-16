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
  host: "0.0.0.0"
  port: 3000
  readTimeout: 5s
  writeTimeout: 10s
grpc:
  host: "0.0.0.0"
  port: 50051
  connectTimeout: 5s
db:
  dsn: "postgres://user:pass@localhost:5432/db"
  migrationsDir: "./migrations/"
storage: "memory"
rmq:
  host: "0.0.0.0"
  port: "5672"
  login: "guest"
  password: "guest"
scheduler:
  period: 60s
  retentionPeriod: 60s
  queue: "calendar_events"
`

func TestNewConfig(t *testing.T) {
	r := bytes.NewReader([]byte(yamlData))
	cfg, err := New(r)
	require.NoError(t, err)
	require.Equal(t, "debug", cfg.Logger.Level)
	require.Equal(t, "0.0.0.0", cfg.HTTP.Host)
	require.Equal(t, "3000", cfg.HTTP.Port)
	require.Equal(t, 5*time.Second, cfg.HTTP.ReadTimeout)
	require.Equal(t, 10*time.Second, cfg.HTTP.WriteTimeout)
	require.Equal(t, "0.0.0.0", cfg.GRPC.Host)
	require.Equal(t, "50051", cfg.GRPC.Port)
	require.Equal(t, 5*time.Second, cfg.GRPC.ConnectTimeout)
	require.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DB.Dsn)
	require.Equal(t, "./migrations/", cfg.DB.MigrationsDir)
	require.Equal(t, "memory", cfg.Storage)
	require.Equal(t, "0.0.0.0", cfg.RMQ.Host)
	require.Equal(t, "5672", cfg.RMQ.Port)
	require.Equal(t, "guest", cfg.RMQ.Login)
	require.Equal(t, "guest", cfg.RMQ.Password)
	require.Equal(t, 60*time.Second, cfg.Scheduler.Period)
	require.Equal(t, 60*time.Second, cfg.Scheduler.RetentionPeriod)
	require.Equal(t, "calendar_events", cfg.Scheduler.Queue)
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
