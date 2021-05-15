package db

import (
	"github.com/Masterminds/squirrel"
)

type PostgresStore struct {
	SB squirrel.StatementBuilderType
	DB *Driver
}

func (m *PostgresStore) Init() {
	m.SB = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	m.DB = GetPool()
}
