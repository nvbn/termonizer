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
	_, err := s.db.ExecContext(ctx, `
		create table if not exists Goals (
		    id text primary key,
		    period integer,
		    content text,
		    start timestamp,
		    updated timestamp
		)`)
	return err
}

func (s *SQLite) ReadForPeriod(ctx context.Context, period int) ([]model.Goal, error) {
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
		return nil, fmt.Errorf("failed to query: %w", err)
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
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		goal.Start = utils.IgnoreTZ(goal.Start)
		result = append(result, goal)
	}

	return result, nil
}

func (s *SQLite) CountForPeriod(ctx context.Context, period int) (int, error) {
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

func (s *SQLite) Update(ctx context.Context, goals model.Goal) error {
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

func (s *SQLite) Close() error {
	return s.db.Close()
}
