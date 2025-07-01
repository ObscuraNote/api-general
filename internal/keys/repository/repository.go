package repository

import (
	"context"
	"log"

	"github.com/ObscuraNote/api-general/internal/keys/dto"
	"github.com/philippe-berto/database/postgresdb"
)

var _ KeysRepository = (*Repository)(nil)

type (
	KeysRepository interface {
		AddKey(userId int64, note dto.KeyImput) error
		GetKeysByUser(userId int64) ([]dto.KeyOutput, error)
		DeleteKey(id string) error
	}
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

func (r *Repository) AddKey(userId int64, note dto.KeyImput) error {
	_, err := r.statements.addKey.statement.
		ExecContext(r.ctx, userId, note.UserAddress, note.EncryptedKey, note.KeyIV, note.EncryptedData, note.DataIV)
	if err != nil {
		log.Println("Error adding note")

		return err
	}

	return nil
}

func (r *Repository) GetKeysByUser(userId int64) ([]dto.KeyOutput, error) {
	rows, err := r.statements.getKeysByUser.statement.
		QueryContext(r.ctx, userId)
	if err != nil {
		log.Println("Error getting notes by user")

		return nil, err
	}
	defer rows.Close()

	var notes []dto.KeyOutput
	for rows.Next() {
		var note dto.KeyOutput
		if err := rows.Scan(&note.ID, &note.EncryptedKey, &note.KeyIV, &note.EncryptedData, &note.DataIV); err != nil {
			log.Println("Error scanning note")

			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (r *Repository) DeleteKey(id string) error {
	_, err := r.statements.deleteKey.statement.
		ExecContext(r.ctx, id)
	if err != nil {
		log.Println("Error deleting note")

		return err
	}

	return nil
}

func (r *Repository) prepareStatements() (statements, error) {
	var err error

	statementsList.addKey.statement, err = r.db.PrepareStatement(statementsList.addKey.query)
	if err != nil {
		return statements{}, err
	}

	statementsList.getKeysByUser.statement, err = r.db.PrepareStatement(statementsList.getKeysByUser.query)
	if err != nil {
		return statements{}, err
	}

	statementsList.deleteKey.statement, err = r.db.PrepareStatement(statementsList.deleteKey.query)
	if err != nil {
		return statements{}, err
	}

	return statementsList, nil
}
