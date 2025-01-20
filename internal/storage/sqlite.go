package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/nvbn/termonizer/internal/utils"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLite(ctx context.Context, path string) (*SQLite, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	s := &SQLite{
		db: db,
	}

	// more explicit?
	if err := s.initSchema(ctx); err != nil {
		return nil, fmt.Errorf("failed to init schema: %w", err)
	}

	return s, nil
}

func (s *SQLite) initSchema(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, `
		create table if not exists Goals (
		    id text primary key,
		    period integer,
		    content text,
		    start timestamp,
		    updated timestamp
		)`); err != nil {
		return fmt.Errorf("failed to create Goals table: %w", err)
	}

	if _, err := s.db.ExecContext(ctx, `
		create table if not exists Settings (
		    id text primary key,
		    value string,
		    updated timestamp
		)`); err != nil {
		return fmt.Errorf("failed to create Settings table: %w", err)
	}

	return nil
}

func (s *SQLite) ReadGoalsForPeriod(ctx context.Context, period int) ([]model.Goal, error) {
	rows, err := s.db.QueryContext(ctx, `
		select
		    id,
		    period,
		    content,
		    start,
		    updated
		from Goals
		where
		    period = ?
		 	and content != ""
		order by start desc
	`, period)
	if err != nil {
		return nil, fmt.Errorf("failed to query goals: %w", err)
	}
	defer rows.Close()

	result := make([]model.Goal, 0)
	for rows.Next() {
		goal := model.Goal{}
		if err := rows.Scan(
			&goal.ID,
			&goal.Period,
			&goal.Content,
			&goal.Start,
			&goal.Updated,
		); err != nil {
			return nil, fmt.Errorf("failed to scan goals: %w", err)
		}
		goal.Start = utils.IgnoreTZ(goal.Start)
		result = append(result, goal)
	}

	return result, nil
}

func (s *SQLite) CountGoalsForPeriod(ctx context.Context, period int) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `
		select
		    count(*)
		from Goals
			where
			    period = ?
			  and content != ""
	`, period).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to query: %w", err)
	}

	return count, nil
}

func (s *SQLite) UpdateGoal(ctx context.Context, goals model.Goal) error {
	_, err := s.db.ExecContext(
		ctx,
		`
			insert or replace into Goals (
				id,
				period,
				content,
				start,
				updated
			) values (?, ?, ?, ?, ?)
		`,
		goals.ID,
		goals.Period,
		goals.Content,
		goals.Start,
		goals.Updated,
	)
	return err
}

func (s *SQLite) ReadSettings(ctx context.Context) ([]model.Setting, error) {
	rows, err := s.db.QueryContext(ctx, `select id, value, updated from Settings`)
	if err != nil {
		return nil, fmt.Errorf("failed to query settings: %w", err)
	}
	defer rows.Close()

	result := make([]model.Setting, 0)
	for rows.Next() {
		setting := model.Setting{}
		if err := rows.Scan(&setting.ID, &setting.Value, &setting.Updated); err != nil {
			return nil, fmt.Errorf("failed to scan settings: %w", err)
		}
		result = append(result, setting)
	}

	return result, nil
}

func (s *SQLite) UpdateSetting(ctx context.Context, settings model.Setting) error {
	_, err := s.db.ExecContext(
		ctx,
		`insert or replace into Settings (id, value, updated) values (?, ?, ?)`,
		settings.ID, settings.Value, settings.Updated,
	)
	return err
}

func (s *SQLite) Cleanup(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		delete from Goals
		where content = ""
	`)

	return err
}

func (s *SQLite) Close() error {
	return s.db.Close()
}
