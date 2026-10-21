package repo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	apperror "github.com/nico-hondo/url-shortener/internal/error"
	"github.com/nico-hondo/url-shortener/internal/model"
	"github.com/redis/go-redis/v9"
)

type AuthRepo struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewAuthRepo(db *pgxpool.Pool, rdb *redis.Client) *AuthRepo {
	return &AuthRepo{
		db:  db,
		rdb: rdb,
	}
}

func (ar *AuthRepo) FindUser(ctx context.Context, email string) (model.Users, error) {

	sql := "SELECT id, email, password, first_name, last_name from users WHERE email=$1"
	args := []any{email}

	var acc model.Users
	if err := ar.db.QueryRow(ctx, sql, args...).Scan(&acc.Id, &acc.Email, &acc.Password, &acc.First_Name, &acc.Last_Name); err != nil {
		return model.Users{}, err
	}

	return acc, nil
}

func (ar *AuthRepo) CreateUser(ctx context.Context, acc model.Users) error {
	sql := "INSERT INTO users (email, password, first_name, last_name) VALUES($1, $2, $3, $4)"

	args := []any{acc.Email, acc.Password, acc.First_Name, acc.Last_Name}

	cmd, err := ar.db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}

	return nil
}

// BlackListToken menyimpan JTI yang di-revoke dengan payload "revoke"
func (ar *AuthRepo) BlackListToken(ctx context.Context, jti string, ttl time.Duration) error {
	// Format key di redis
	key := fmt.Sprintf("%s:bl:%s", os.Getenv("RDB_USER"), jti)

	// Cukup simpan "revoke" sebagai value penanda
	return ar.rdb.Set(ctx, key, "revoke", ttl).Err()
}

// IsBlacklisted mengecek keberadaan key di Redis menggunakan EXISTS
func (a *AuthRepo) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("%s:bl:%s", os.Getenv("RDB_USER"), jti)

	// Exists jauh lebih efisien daripada Get karena tidak perlu mentransfer data value
	count, err := a.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
