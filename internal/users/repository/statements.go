package repository

import "github.com/jmoiron/sqlx"

type statementsItem struct {
	name      string
	query     string
	statement *sqlx.Stmt
}

type statements struct {
	createUser      statementsItem
	getUserId       statementsItem
	checkUserExists statementsItem
	updatePassword  statementsItem
	deleteUser      statementsItem
}

var statementsList = statements{
	createUser: statementsItem{
		name: "createUser",
		query: `
            INSERT INTO users (user_address, password)
            VALUES ($1, $2);`,
	},
	getUserId: statementsItem{
		name: "getUserId",
		query: `
            SELECT id
            FROM users
            WHERE user_address = $1
            AND password = $2;`,
	},
	checkUserExists: statementsItem{
		name: "checkUserExists",
		query: `
            SELECT EXISTS(
            SELECT 1 FROM users 
            WHERE user_address = $1
						AND password = $2
        );`,
	},
	updatePassword: statementsItem{
		name: "updatePassword",
		query: `
            UPDATE users
            SET password = $2, updated_at = CURRENT_TIMESTAMP
            WHERE id = $1;`,
	},
	deleteUser: statementsItem{
		name: "deleteUser",
		query: `
            DELETE FROM users
            WHERE id = $1;`,
	},
}
