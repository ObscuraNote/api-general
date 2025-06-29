package repository

import "github.com/jmoiron/sqlx"

type statementsItem struct {
	name      string
	query     string
	statement *sqlx.Stmt
}

type statements struct {
	addKey        statementsItem
	getKeysByUser statementsItem
	deleteKey     statementsItem
}

var statementsList = statements{
	addKey: statementsItem{
		name: "addKey",
		query: `
			INSERT INTO keys (user_id, user_address, encrypted_key, key_iv, encrypted_data, data_iv)
        VALUES ($1, $2, $3, $4, $5, $6);`,
	},
	getKeysByUser: statementsItem{
		name: "getKeysByUser",
		query: `
      SELECT id, encrypted_key, key_iv, encrypted_data, data_iv
      FROM keys
      WHERE user_id = $1
			ORDER BY created_at DESC;`,
	},
	deleteKey: statementsItem{
		name: "deleteKey",
		query: `
			DELETE FROM keys
			WHERE id = $1;`,
	},
}
