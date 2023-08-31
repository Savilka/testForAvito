package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	api "testForAvito/internal/api/v1/models"
	"time"
)

type Storage struct {
	db *pgx.Conn
}

func Connect(storageUrl string) (*Storage, error) {
	const op = "storage.postgres.Connect"

	db, err := pgx.Connect(context.Background(), storageUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) AddNewUser() (string, error) {
	const op = "storage.postgres.AddNewUser"

	var id string

	err := s.db.QueryRow(context.Background(), "INSERT INTO \"user\" DEFAULT VALUES RETURNING id").Scan(&id)
	if err != nil {
		return "", fmt.Errorf("%s: %v", op, err)
	}

	return id, err
}

func (s *Storage) AddSegment(slug string, percent float32) (string, error) {
	const op = "storage.postgres.AddSegment"

	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return "", fmt.Errorf("%s: %v", op, err)
	}

	defer tx.Rollback(context.Background())

	var segmentId string
	err = tx.QueryRow(context.Background(), "INSERT INTO segment(name) VALUES ($1) ON CONFLICT DO NOTHING RETURNING id", slug).Scan(&segmentId)
	if err != nil {
		return "", fmt.Errorf("%s: %v", op, err)
	}

	if percent > 0 {
		var ids []string
		rows, err := tx.Query(context.Background(), "SELECT id FROM \"user\" TABLESAMPLE BERNOULLI($1)", percent)
		if err != nil {
			return "", fmt.Errorf("%s: %v", op, err)
		}

		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				return "", fmt.Errorf("%s: %v", op, err)
			}
			ids = append(ids, id)
		}

		if err := rows.Err(); err != nil {
			return "", fmt.Errorf("%s: %v", op, err)
		}

		rows.Close()

		for _, id := range ids {
			_, err = tx.Exec(context.Background(), "INSERT INTO user_segment(user_id, segment_name) VALUES ($1, $2) ON CONFLICT DO NOTHING",
				id, slug)
			if err != nil {
				return "", fmt.Errorf("%s: %v", op, err)
			}
		}
	}

	_ = tx.Commit(context.Background())

	return segmentId, err
}

func (s *Storage) DeleteSegment(slug string) error {
	const op = "storage.postgres.DeleteSegment"

	_, err := s.db.Exec(context.Background(), "DELETE FROM segment WHERE name = $1", slug)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return err
}

func (s *Storage) AddUserToSegments(newSegments []api.Segment, oldSegments []api.Segment, userId string) error {
	const op = "storage.postgres.AddUserToSegment"

	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	defer tx.Rollback(context.Background())

	for _, segment := range oldSegments {
		_, err := tx.Exec(context.Background(), "DELETE FROM user_segment WHERE user_id = $1 AND segment_name = $2", userId, segment.Name)
		if err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}
	}

	for _, segment := range newSegments {
		_, err := tx.Exec(context.Background(), "INSERT INTO user_segment(user_id, segment_name, ttl, date_insert) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING",
			userId, segment.Name, segment.Ttl, time.Now().Unix())
		if err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}
	}

	_ = tx.Commit(context.Background())

	return err
}

func (s *Storage) GetUserSegments(id string) ([]string, error) {
	const op = "storage.postgres.GetUserSegments"

	rows, err := s.db.Query(context.Background(), "SELECT segment.name FROM user_segment JOIN segment ON user_segment.segment_name = segment.name AND user_segment.user_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}
		names = append(names, name)
	}

	return names, err
}

func (s *Storage) DeleteExpiredSegmentsFromUsers() error {
	const op = "storage.postgres.DeleteExpiredSegmentsFromUsers"

	_, err := s.db.Exec(context.Background(), "DELETE FROM user_segment WHERE current_timestamp - to_timestamp(date_insert) >= make_interval(0, 0, 0, 0, 0, 0, ttl);\n")
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return err
}
