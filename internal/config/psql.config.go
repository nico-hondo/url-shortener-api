package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PsqlDb struct {
	User   string
	Pass   string
	Host   string
	Port   string
	DbName string
}

func NewPsqlDb(user, pass, host, port, dbName string) *PsqlDb {
	return &PsqlDb{
		User:   user,
		Pass:   pass,
		Host:   host,
		Port:   port,
		DbName: dbName,
	}
}

func (p *PsqlDb) Connect() (*pgxpool.Pool, error) {
	// Client Connection
	// pgx.Connect()
	// Pool Connection
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", p.User, p.Pass, p.Host, p.Port, p.DbName)
	// log.Println(connStr)
	// context.Background()
	return pgxpool.New(context.Background(), connStr)
	// pool.Query()
}
