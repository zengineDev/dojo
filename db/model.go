package db

import (
	"github.com/Masterminds/squirrel"
)

type Model struct {
	SB squirrel.StatementBuilderType
	DB *Driver
}

func (m *Model) Init() {
	m.SB = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	m.DB = GetPool()
}
