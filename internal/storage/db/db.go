package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strings"
	"time"
	"url-saver-bot/internal/storage"
)

const (
	table = "links"
)

type DBStorage struct {
	pool *pgxpool.Pool
}

func NewDBStorage(ctx context.Context, DSN string) *DBStorage {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, DSN)
	if err != nil {
		log.Fatal(err)
	}

	exist, err := isTableExists(ctx, pool)
	if err != nil {
		log.Fatal(err)
	}
	if !exist {
		err = createTable(ctx, pool)
		if err != nil {
			log.Fatal(err)
		}
	}
	return &DBStorage{pool: pool}
}

// Save check if page already exists and save if not
func (s *DBStorage) Save(ctx context.Context, p *storage.Page) error {
	var u string
	err := s.pool.QueryRow(ctx, "SELECT url FROM links WHERE url = $1 AND user_name = $2", p.URL, p.UserName).Scan(&u)
	if u != "" {
		return storage.NewAlreadyExistsError()
	}
	_, err = s.pool.Exec(ctx, "INSERT INTO links (url, user_name, tags, created_time) VALUES ($1, $2, $3, $4)", p.URL, p.UserName, p.Tags, p.Created)
	if err != nil {
		return fmt.Errorf("storage can't save page: %w", err)
	}
	return nil
}

func (s *DBStorage) Pick(ctx context.Context, userName string) (*storage.Page, error) {
	var p storage.Page
	err := s.pool.QueryRow(ctx, "SELECT url, user_name, tags, created_time FROM links WHERE user_name = $1 ORDER BY created_time LIMIT 1", userName).Scan(&p.URL, &p.UserName, &p.Tags, &p.Created)
	if err == pgx.ErrNoRows {
		return &storage.Page{}, storage.NewNoResultError()
	} else if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *DBStorage) Remove(ctx context.Context, p *storage.Page) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM links WHERE url = $1 AND user_name = $2", p.URL, p.UserName)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBStorage) PickAll(ctx context.Context, userName string) ([]storage.Page, error) {
	pages := make([]storage.Page, 0, 20)

	rows, err := s.pool.Query(ctx, "SELECT url, user_name, tags, created_time FROM links WHERE user_name = $1 ORDER BY created_time", userName)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("can't pick all rows: %w", err)
	}

	for rows.Next() {
		var p storage.Page
		err = rows.Scan(&p.URL, &p.UserName, &p.Tags, &p.Created)
		if err != nil {
			return nil, fmt.Errorf("can't scan row: %w", err)
		}
		pages = append(pages, p)
	}

	return pages, nil
}

func (s *DBStorage) SelectTags(ctx context.Context, userName string) ([]string, error) {
	tags := make([]string, 0, 10)

	rows, err := s.pool.Query(ctx, "SELECT DISTINCT tags FROM links WHERE user_name = $1 AND tags != '' ORDER BY tags", userName)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("can't select tags: %w", err)
	}

	for rows.Next() {
		var t string
		err = rows.Scan(&t)
		if err != nil {
			return nil, fmt.Errorf("can't scan row: %w", err)
		}
		tags = append(tags, t)
	}

	return tags, nil
}

func (s *DBStorage) SelectByTag(ctx context.Context, tag string, userName string) ([]string, error) {
	rows, err := s.pool.Query(ctx, "SELECT url FROM links WHERE user_name = $1 AND tags = $2 ORDER BY created_time", userName, tag)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("can't select rows: %w", err)
	}

	urls := make([]string, 0, 10)
	for rows.Next() {
		var url string
		err = rows.Scan(&url)
		if err != nil {
			return nil, fmt.Errorf("can't scan row: %w", err)
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (s *DBStorage) BatchUpdate(ctx context.Context, pages []storage.Page) error {
	b := &pgx.Batch{}
	for _, v := range pages {
		b.Queue("UPDATE links SET tags = $1 WHERE url = $2 AND user_name = $3", v.Tags, v.URL, v.UserName)
	}
	con, err := s.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer con.Release()
	res := con.SendBatch(ctx, b)
	_, err = res.Exec()
	if err != nil {
		return fmt.Errorf("error exec batch: %w", err)
	}
	return nil
}

func createTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "CREATE TABLE "+table+" (id serial primary key, url varchar, user_name varchar, tags varchar, created_time timestamptz)")
	if err != nil {
		return err
	}
	return nil
}

func isTableExists(ctx context.Context, pool *pgxpool.Pool) (bool, error) {
	_, err := pool.Query(ctx, "SELECT * FROM links LIMIT 1")
	if err != nil && strings.Contains(err.Error(), "не существует") {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
