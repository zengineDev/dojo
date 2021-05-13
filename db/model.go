package db

import (
	"github.com/Masterminds/squirrel"
)

type Model struct {
	SB squirrel.StatementBuilderType
	DB *Driver
}

func (m Model) NewModel() *Model {
	return &Model{
		SB: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		DB: GetPool(),
	}
}
