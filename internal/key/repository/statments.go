package repository

import "github.com/jmoiron/sqlx"

type statementsItem struct {
	name      string
	query     string
	statement *sqlx.Stmt
}

type statments struct {
	addNote        statementsItem
	getNotesByUser statementsItem
}

var statementsList = statments{
	addNote: statementsItem{
		name: "addNote",
		query: `
			INSERT INTO user_data (user_address, key, encrypted_data, iv)
        VALUES ($1, $2, $3, $4);`,
	},
	getNotesByUser: statementsItem{
		name: "getNotesByUser",
		query: `
      SELECT key, encrypted_data, iv
      FROM user_data
      WHERE user_address = $1;`,
	},
}
