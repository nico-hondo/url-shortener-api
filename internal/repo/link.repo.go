package repo

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	apperror "github.com/nico-hondo/url-shortener/internal/error"
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

func (lr *LinkRepo) GetUserLinks(ctx context.Context, userId int) ([]model.Link, error) {
	sql := `SELECT user_id, original_url, slug, created_at FROM links WHERE user_id = $1 AND deleted_at IS NULL`

	rows, err := lr.db.Query(ctx, sql, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var links []model.Link = []model.Link{}

	for rows.Next() {
		var link model.Link
		if err := rows.Scan(&link.User_Id, &link.Original_Url, &link.Slug, &link.Created_At); err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func (lr *LinkRepo) DeleteLinkRepo(ctx context.Context, linkID int, userID int) error {
	sql := `UPDATE links SET deleted_at = NOW() WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`

	cmd, err := lr.db.Exec(ctx, sql, linkID, userID)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG DB] RowsAffected: %d | Input LinkID: %d | Input UserID: %d\n", cmd.RowsAffected(), linkID, userID)

	if cmd.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}

	return nil
}
