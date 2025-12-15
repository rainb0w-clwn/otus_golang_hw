package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // driver import
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/entity"
)

type PgStorage struct {
	db  *sqlx.DB
	ctx context.Context
}

type sqlEvent struct {
	ID             string         `db:"id"`
	UserID         int            `db:"user_id"`
	Title          string         `db:"title"`
	DateTime       time.Time      `db:"datetime"`
	Description    sql.NullString `db:"description"`
	Duration       sql.NullString `db:"duration"`
	RemindTime     sql.NullTime   `db:"remind_time"`
	RemindSentTime sql.NullTime   `db:"remind_sent_time"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}

var ErrConnectFailed = errors.New("error connecting to db")

func (s *PgStorage) Create(event entity.Event) (string, error) {
	query := `
		INSERT INTO event (
			user_id, title, description, datetime, duration, remind_time
		) VALUES (
			:user_id, :title, :description, :datetime, :duration, :remind_time
		)
		RETURNING id
	`

	params := map[string]any{
		"user_id":     event.UserID,
		"title":       event.Title,
		"description": event.Description,
		"datetime":    event.DateTime,
		"duration":    event.Duration,
		"remind_time": event.RemindTime,
	}

	var id string
	stmt, err := s.db.PrepareNamedContext(s.ctx, query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	if err = stmt.GetContext(s.ctx, &id, params); err != nil {
		return "", err
	}

	return id, nil
}

func (s *PgStorage) GetByID(id string) (*entity.Event, error) {
	query := `
		SELECT *
		FROM event
		WHERE id = :id
	`

	stmt, err := s.db.PrepareNamedContext(s.ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var se sqlEvent
	err = stmt.GetContext(
		s.ctx,
		&se,
		map[string]any{"id": id},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrEventNotFound
		}
		return nil, err
	}

	return s.sqlEventToEvent(&se), nil
}

func (s *PgStorage) GetAll() (*entity.Events, error) {
	query := `SELECT * FROM event`

	var rows []sqlEvent
	if err := s.db.SelectContext(s.ctx, &rows, query); err != nil {
		return nil, err
	}

	events := make(entity.Events, 0, len(rows))
	for _, r := range rows {
		events = append(events, s.sqlEventToEvent(&r))
	}

	return &events, nil
}

func (s *PgStorage) Update(event entity.Event) error {
	query := `
		UPDATE event SET
			user_id     = :user_id,
			title       = :title,
			description = :description,
			datetime    = :datetime,
			duration    = :duration,
			remind_time = :remind_time,
			updated_at  = now()
		WHERE id = :id
	`

	params := map[string]any{
		"id":          event.ID,
		"user_id":     event.UserID,
		"title":       event.Title,
		"description": event.Description,
		"datetime":    event.DateTime,
		"duration":    event.Duration,
		"remind_time": event.RemindTime,
	}

	_, err := s.db.NamedExecContext(s.ctx, query, params)
	return err
}

func (s *PgStorage) Delete(id string) error {
	query := `
		DELETE FROM event
		WHERE id = :id
	`

	_, err := s.db.NamedExecContext(
		s.ctx,
		query,
		map[string]any{"id": id},
	)
	return err
}

func (s *PgStorage) GetForPeriod(start time.Time, end time.Time) (*entity.Events, error) {
	query := `
		SELECT *
		FROM event
		WHERE datetime BETWEEN :start AND :end
	`

	stmt, err := s.db.PrepareNamedContext(s.ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows []sqlEvent
	err = stmt.SelectContext(
		s.ctx,
		&rows,
		map[string]any{
			"start": start,
			"end":   end,
		},
	)
	if err != nil {
		return nil, err
	}

	events := make(entity.Events, 0, len(rows))
	for _, r := range rows {
		events = append(events, s.sqlEventToEvent(&r))
	}

	return &events, nil
}

func (s *PgStorage) GetForTime(t time.Time) (*entity.Event, error) {
	query := `
		SELECT *
		FROM event
		WHERE datetime = :datetime
	`

	stmt, err := s.db.PrepareNamedContext(s.ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var se sqlEvent
	err = stmt.GetContext(
		s.ctx,
		&se,
		map[string]any{
			"datetime": t,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrEventNotFound
		}
		return nil, err
	}

	return s.sqlEventToEvent(&se), nil
}

func (s *PgStorage) GetForRemind() (*entity.Events, error) {
	query := `
		SELECT *
		FROM event
		WHERE remind_sent_time IS NULL AND now() >= remind_time
	`

	var rows []sqlEvent
	err := s.db.SelectContext(
		s.ctx,
		&rows,
		query,
	)
	if err != nil {
		return nil, err
	}

	events := make(entity.Events, 0, len(rows))
	for _, r := range rows {
		events = append(events, s.sqlEventToEvent(&r))
	}

	return &events, nil
}

func (s *PgStorage) MarkAsReminded(id string) error {
	query := `
		UPDATE event SET
			remind_sent_time = now(),
			updated_at  = now()
		WHERE id = :id
	`

	params := map[string]any{
		"id": id,
	}

	_, err := s.db.NamedExecContext(s.ctx, query, params)
	return err
}

func New() *PgStorage {
	return &PgStorage{}
}

func (s *PgStorage) DeleteOlderThan(t time.Time) error {
	query := `
		DELETE event
		WHERE datetime < :time
	`

	params := map[string]any{
		"time": t,
	}

	_, err := s.db.NamedExecContext(s.ctx, query, params)
	return err
}

func (s *PgStorage) Connect(ctx context.Context) error {
	cfg := config.GetFromContext(ctx)
	if cfg == nil {
		return config.ErrNoConfigInContext
	}

	db, err := sqlx.Open("pgx", cfg.DB.Dsn)
	if err != nil {
		return fmt.Errorf(ErrConnectFailed.Error()+":%w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf(ErrConnectFailed.Error()+":%w", err)
	}

	s.db = db
	s.ctx = ctx
	if cfg.DB.Migrate {
		return s.migrate(cfg.DB.MigrationsDir)
	}
	return nil
}

func (s *PgStorage) Close(_ context.Context) error {
	if s.db == nil {
		return nil
	}

	err := s.db.Close()
	s.db = nil
	s.ctx = nil

	return err
}

func (s *PgStorage) sqlEventToEvent(se *sqlEvent) *entity.Event {
	e := &entity.Event{
		ID:        se.ID,
		UserID:    se.UserID,
		Title:     se.Title,
		DateTime:  se.DateTime,
		CreatedAt: se.CreatedAt,
		UpdatedAt: se.UpdatedAt,
	}

	if se.Description.Valid {
		e.Description = se.Description.String
	}
	if se.Duration.Valid {
		e.Duration = se.Duration.String
	}
	if se.RemindTime.Valid {
		e.RemindTime = se.RemindTime.Time
	}
	if se.RemindTime.Valid {
		e.RemindSentTime = se.RemindSentTime.Time
	}

	return e
}

func (s *PgStorage) migrate(migrationDir string) error {
	if s.db == nil {
		return fmt.Errorf("database connection is not established")
	}

	if err := goose.SetDialect("pgx"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(s.db.DB, migrationDir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
