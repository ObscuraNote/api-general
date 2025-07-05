package repository

import (
	"context"

	"github.com/philippe-berto/database/postgresdb"
)

var _ UsersRepository = (*Repository)(nil)

type (
	UsersRepository interface {
		CreateUser(userAddress, password string) error
		GetUserId(userAddress, password string) (int64, error)
		CheckUserExists(userAddress, password string) (bool, error)
		UpdatePassword(userId int64, password string) error
		DeleteUser(userId int64) (bool, error)
	}
	Repository struct {
		ctx        context.Context
		db         *postgresdb.Client
		statements statements
	}
)

func New(ctx context.Context, db *postgresdb.Client) (*Repository, error) {
	r := &Repository{
		ctx:        ctx,
		db:         db,
		statements: statements{},
	}
	statements, err := r.prepareStatements()
	if err != nil {
		return &Repository{}, err
	}

	r.statements = statements

	return r, nil
}

func (r *Repository) CreateUser(userAddress, password string) error {
	_, err := r.statements.createUser.statement.
		ExecContext(r.ctx, userAddress, password)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserId(userAddress, password string) (int64, error) {
	var id int64
	err := r.statements.getUserId.statement.
		QueryRowContext(r.ctx, userAddress, password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) CheckUserExists(userAddress, password string) (bool, error) {
	var exists bool
	err := r.statements.checkUserExists.statement.
		QueryRowContext(r.ctx, userAddress, password).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) UpdatePassword(userId int64, password string) error {
	_, err := r.statements.updatePassword.statement.
		ExecContext(r.ctx, userId, password)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteUser(userId int64) (bool, error) {
	result, err := r.statements.deleteUser.statement.
		ExecContext(r.ctx, userId)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
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
