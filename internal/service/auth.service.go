package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nico-hondo/url-shortener/internal/dto"
	apperror "github.com/nico-hondo/url-shortener/internal/error"
	"github.com/nico-hondo/url-shortener/internal/model"
	"github.com/nico-hondo/url-shortener/internal/repo"
	"github.com/nico-hondo/url-shortener/pkg"
)

type AuthService struct {
	ar *repo.AuthRepo
}

func NewAuthService(ar *repo.AuthRepo) *AuthService {
	return &AuthService{
		ar: ar,
	}
}

func (as *AuthService) RegisterService(ctx context.Context, body dto.RegisterRequest) error {
	//validasi
	if len(body.First_Name) == 0 || len(body.Last_Name) == 0 || len(body.Email) == 0 || len(body.Password) == 0 {
		return errors.New("name, email and password tidak boleh kosong")
	}

	//Cek exists di repo
	_, err := as.ar.FindUser(ctx, body.Email)
	if err == nil {
		return apperror.ErrEmailAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	//hash pwd
	hc := pkg.NewRecommendedHashConfig()
	hashedPwd := hc.GenHash(body.Password)

	//menambahkan data ke db
	if err := as.ar.CreateUser(ctx, model.Users{
		Email:      body.Email,
		Password:   hashedPwd,
		First_Name: body.First_Name,
		Last_Name:  body.Last_Name,
	}); err != nil {
		return err
	}

	return nil
}

func (as *AuthService) LoginService(ctx context.Context, body dto.LoginRequest) (string, dto.LoginResponse, error) {
	//validas
	if len(body.Email) == 0 || len(body.Password) == 0 {
		return "", dto.LoginResponse{}, apperror.ErrCredentialsEmpty
	}

	//cek account apakah sudah ada di db (melalui repo)
	acc, err := as.ar.FindUser(ctx, body.Email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", dto.LoginResponse{}, apperror.ErrInvalidCredentials
		}
		return "", dto.LoginResponse{}, err
	}

	if err := pkg.Compare(body.Password, acc.Password); err != nil {
		log.Printf("DEBUG compare failed: %v", err)
		return "", dto.LoginResponse{}, apperror.ErrInvalidCredentials
	}

	//Jika tidak terjadi error, maka berikan nilai truenya
	claims := pkg.NewJwtClaims(acc.Id)
	token, err := claims.GenToken()
	if err != nil {
		return "", dto.LoginResponse{}, err
	}

	userRespon := dto.LoginResponse{
		ID:         acc.Id,
		Email:      acc.Email,
		First_Name: acc.First_Name,
		Last_Name:  acc.Last_Name,
	}

	return token, userRespon, err
}

func (as *AuthService) Logout(ctx context.Context, jti string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)

	// MENCEGAH POLLUTION MEMORY:
	// Jika token sudah expired secara alami, jangan simpan ke Redis!
	if ttl <= 0 {
		return nil
	}

	return as.ar.BlackListToken(ctx, jti, ttl)
}
