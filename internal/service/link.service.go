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

func (ls *LinkService) GetOriginalURL(ctx context.Context, slug string) (string, error) {
	originalURL, err := ls.lr.GetOgURLBySlug(ctx, slug)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}

func (ls *LinkService) GetUserLink(ctx context.Context, userId int) ([]dto.UserLinkListResponse, error) {
	links, err := ls.lr.GetUserLinks(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Inisialisasi slice penampung DTO (Gunakan slice kosong agar saat di-JSON-kan hasilnya [] bukan null)
	res := make([]dto.UserLinkListResponse, 0, len(links))

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}

	for _, link := range links {
		res = append(res, dto.UserLinkListResponse{
			User_Id:     link.User_Id,
			OriginalURL: link.Original_Url,
			Slug:        link.Slug,
			ShortURL:    fmt.Sprintf("%s/%s", baseURL, link.Slug),
			CreatedAt:   link.Created_At,
		})
	}

	return res, nil
}

func (ls *LinkService) DeleteLinkService(ctx context.Context, linkID int, userID int) error {
	err := ls.lr.DeleteLinkRepo(ctx, linkID, userID)
	if err != nil {
		return err
	}

	return nil
}
