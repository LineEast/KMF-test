package database

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Database struct {
		Pool *pgxpool.Pool
	}

	Configuration struct {
		User     string `toml:"user"`
		Database string `toml:"database"`
		Password string `toml:"password"`
		Host     string `toml:"host"`
		Port     string `toml:"port"`
	}
)

func Conn(configuration *Configuration) (database *Database, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DSN := strings.Join(
		[]string{
			"user=" + configuration.User,
			"database=" + configuration.Database,
			"password=" + configuration.Password,
			"host=" + configuration.Host,
			"port=" + configuration.Port,
		},
		" ",
	)

	database = new(Database)
	database.Pool, err = pgxpool.New(ctx, DSN)
	if err != nil {
		panic(err)
	}

	return
}
