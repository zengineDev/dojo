package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/zengineDev/dojo"
	"sync"
)

var once sync.Once

type Driver struct {
	Pool *pgxpool.Pool
}

var (
	instance *Driver
)

func GetPool() *Driver {
	once.Do(func() {
		cfg := dojo.GetConfig()
		var err error

		config, err := pgxpool.ParseConfig(cfg.DB.DSN())
		if err != nil {
			panic(err)
		}
		config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			// do something with every new connection
			//fmt.Println("connect")
			return nil
		}

		pool, err := pgxpool.ConnectConfig(context.Background(), config)

		if err != nil {
			panic(err)
		}

		instance = &Driver{Pool: pool}
	})

	return instance
}
