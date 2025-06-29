package repository

import (
	"context"
	"log"

	"github.com/philippe-berto/database/postgresdb"
)

type (
	Repository struct {
		ctx        context.Context
		db         *postgresdb.Client
		statements statements
	}
)

func New(ctx context.Context, db *postgresdb.Client) *Repository {
	r := &Repository{
		ctx:        ctx,
		db:         db,
		statements: statements{},
	}
	statements, err := r.prepareStatements()
	if err != nil {
		panic(err)
	}

	r.statements = statements

	return r
}

func (r *Repository) CreateUser(userAddress, password string) error {
	_, err := r.statements.createUser.statement.
		ExecContext(r.ctx, userAddress, password)
	if err != nil {
		log.Println("Error creating user")
		return err
	}

	return nil
}

func (r *Repository) GetUserId(userAddress, password string) (int64, error) {
	var id int64
	err := r.statements.getUserId.statement.
		QueryRowContext(r.ctx, userAddress, password).Scan(&id)
	if err != nil {
		log.Println("Error getting user id")
		return 0, err
	}

	return id, nil
}

func (r *Repository) CheckUserExists(userAddress, password string) (bool, error) {
	var exists bool
	err := r.statements.checkUserExists.statement.
		QueryRowContext(r.ctx, userAddress, password).Scan(&exists)
	if err != nil {
		log.Println("Error checking if user exists")
		return false, err
	}

	return exists, nil
}

func (r *Repository) UpdatePassword(userId int64, password string) error {
	_, err := r.statements.updatePassword.statement.
		ExecContext(r.ctx, userId, password)
	if err != nil {
		log.Println("Error updating password")
		return err
	}

	return nil
}

func (r *Repository) DeleteUser(userAddress string) (bool, error) {
	result, err := r.statements.deleteUser.statement.
		ExecContext(r.ctx, userAddress)
	if err != nil {
		log.Println("Error deleting user")
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected")
		return false, err
	}

	return rowsAffected > 0, nil
}

func (r *Repository) prepareStatements() (statements, error) {
	var err error

	statementsList.createUser.statement, err = r.db.PrepareStatement(statementsList.createUser.query)
	if err != nil {
		return statements{}, err
	}

	statementsList.getUserId.statement, err = r.db.PrepareStatement(statementsList.getUserId.query)
	if err != nil {
		return statements{}, err
	}

	statementsList.checkUserExists.statement, err = r.db.PrepareStatement(statementsList.checkUserExists.query)
	if err != nil {
		return statements{}, err
	}

	statementsList.updatePassword.statement, err = r.db.PrepareStatement(statementsList.updatePassword.query)
	if err != nil {
		return statements{}, err
	}

	statementsList.deleteUser.statement, err = r.db.PrepareStatement(statementsList.deleteUser.query)
	if err != nil {
		return statements{}, err
	}

	return statementsList, nil
}
