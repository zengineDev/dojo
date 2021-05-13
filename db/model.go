package db

import (
	"github.com/Masterminds/squirrel"
)

type User struct {
	Model
}

type Model struct {
	sb squirrel.StatementBuilderType
	db *Driver
}

func (m Model) NewModel() *Model {
	return &Model{
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		db: GetPool(),
	}
}
