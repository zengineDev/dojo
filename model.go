package dojo

import (
	"github.com/Masterminds/squirrel"
	"github.com/zengineDev/dojo/db"
)

type User struct {
	Model
}

type Model struct {
	sb squirrel.StatementBuilderType
	db *db.Driver
}

func (m Model) NewModel() *Model {
	return &Model{
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		db: db.GetPool(),
	}
}
