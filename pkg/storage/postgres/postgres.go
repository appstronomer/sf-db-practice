package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"GoNews/pkg/storage"
)

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

// Конструктор объекта хранилища.
func New(constr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			p.id,
			p.title,
			p.content,
			p.created_at,
			p.published_at,
			p.author_id,
			a.name
		FROM posts p
		LEFT JOIN authors a ON p.author_id = a.id
		ORDER BY id;
	`)
	if err != nil {
		return nil, err
	}

	var posts []storage.Post
	for rows.Next() {
		var t storage.Post
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&t.CreatedAt,
			&t.PublishedAt,
			&t.AuthorID,
			&t.AuthorName,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, t)

	}
	return posts, rows.Err()
}

func (s *Store) AddPost(post storage.Post) error {
	curTimestamp := time.Now().UnixMilli()
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO posts (
			title,
			content,
			created_at,
			published_at,
			author_id
		)
		VALUES ($1, $2, $3, $4, $5);
		`,
		post.Title,
		post.Content,
		curTimestamp,
		curTimestamp,
		post.AuthorID,
	)
	return err
}

func (s *Store) UpdatePost(post storage.Post) error {
	curTimestamp := time.Now().UnixMilli()
	tag, err := s.db.Exec(context.Background(), `
		UPDATE posts
		SET
			title = $1,
			content = $2,
			published_at = $3,
			author_id = $4
		WHERE id = $5;
		`,
		post.Title,
		post.Content,
		curTimestamp,
		post.AuthorID,
		post.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("update post id=%v : no rows affected", post.ID)
	}
	return nil
}

func (s *Store) DeletePost(post storage.Post) error {
	tag, err := s.db.Exec(context.Background(), `
		DELETE 
		FROM posts
		WHERE id = $1;
		`,
		post.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delete post id=%v: no rows affected", post.ID)
	}
	return nil
}
