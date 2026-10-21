package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nico-hondo/url-shortener/internal/model"
)

type LinkRepo struct {
	db *pgxpool.Pool
}

func NewLinkRepo(db *pgxpool.Pool) *LinkRepo {
	return &LinkRepo{
		db: db,
	}
}

func (lr *LinkRepo) CreateLink(ctx context.Context, userID int, ogURL string, slug string) (model.Link, error) {
	sql := `insert into links(user_id, original_url, slug) values($1, $2, $3) RETURNING user_id, original_url, slug`

	args := []any{userID, ogURL, slug}

	var body model.Link

	if err := lr.db.QueryRow(ctx, sql, args...).Scan(&body.User_Id, &body.Original_Url, &body.Slug); err != nil {
		return model.Link{}, err
	}

	return body, nil
}

func (lr *LinkRepo) GetOgURLBySlug(ctx context.Context, slug string) (string, error) {
	sql := `SELECT original_url FROM links WHERE slug = $1 AND deleted_at IS NULL`

	var originalURL string
	if err := lr.db.QueryRow(ctx, sql, slug).Scan(&originalURL); err != nil {
		return "", err
	}

	return originalURL, nil
}
