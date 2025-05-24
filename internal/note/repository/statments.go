package repository

import "github.com/jmoiron/sqlx"

type statementsItem struct {
	name      string
	query     string
	statement *sqlx.Stmt
}

type statements struct {
	addNote        statementsItem
	getNotesByUser statementsItem
}

var statementsList = statements{
	addNote: statementsItem{
		name: "addNote",
		query: `
			INSERT INTO keys (user_address, key, encrypted_data, iv)
        VALUES ($1, $2, $3, $4);`,
	},
	getNotesByUser: statementsItem{
		name: "getNotesByUser",
		query: `
      SELECT key, encrypted_data, iv
      FROM keys
      WHERE user_address = $1;`,
	},
}
