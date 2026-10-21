package service

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/nico-hondo/url-shortener/internal/dto"
	"github.com/nico-hondo/url-shortener/internal/repo"
)

type LinkService struct {
	lr *repo.LinkRepo
}

func NewLinkService(lr *repo.LinkRepo) *LinkService {
	return &LinkService{
		lr: lr,
	}
}

func (ls *LinkService) CreateLink(ctx context.Context, userId int, body dto.CreateLinkRequest) (dto.CreateLinkResponse, error) {
	slug := ""
	randomSlugLength := 6

	if body.Slug != nil && *body.Slug != "" {
		slug = *body.Slug
	} else {
		slug = ls.GenerateRandSlug(randomSlugLength)
	}

	link, err := ls.lr.CreateLink(ctx, userId, body.OriginalURL, slug)
	if err != nil {
		return dto.CreateLinkResponse{}, err
	}

	// Opsional: Buat full ShortURL untuk response
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}

	return dto.CreateLinkResponse{
		ID:          link.User_Id,
		OriginalURL: link.Original_Url,
		Slug:        link.Slug,
		ShortURL:    fmt.Sprintf("%s/%s", baseURL, link.Slug),
	}, nil
}

func (ls *LinkService) GenerateRandSlug(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.IntN(len(letterBytes))]
	}
	return string(b)
}
