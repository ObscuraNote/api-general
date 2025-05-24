package repository

import (
	"context"
	"log"

	"github.com/philippe-berto/database/postgresdb"
)

type (
	Repository struct {
		db         *postgresdb.Client
		statements statements
	}
	keyImput struct {
		UserAddress   string
		Key           []byte
		EncryptedData []byte
		IV            []byte
	}
	keyOutput struct {
		ID            [16]byte
		Key           []byte
		EncryptedData []byte
		IV            []byte
	}
)

func New(db *postgresdb.Client) *Repository {
	r := &Repository{
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

func (r *Repository) AddKey(ctx context.Context, note keyImput) error {
	_, err := r.statements.addNote.statement.
		ExecContext(ctx, note.UserAddress, note.Key, note.EncryptedData, note.IV)
	if err != nil {
		log.Println("Error adding note")

		return err
	}

	return nil
}

func (r *Repository) GetKeysByUser(ctx context.Context, userAddress string) ([]keyOutput, error) {
	rows, err := r.statements.getNotesByUser.statement.
		QueryContext(ctx, userAddress)
	if err != nil {
		log.Println("Error getting notes by user")

		return nil, err
	}
	defer rows.Close()

	var notes []keyOutput
	for rows.Next() {
		var note keyOutput
		if err := rows.Scan(&note.Key, &note.EncryptedData, &note.IV); err != nil {
			log.Println("Error scanning note")

			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (r *Repository) prepareStatements() (statements, error) {
	var err error

	statementsList.addNote.statement, err = r.db.PrepareStatement(statementsList.addNote.query)
	if err != nil {
		return statements{}, err
	}

	statementsList.getNotesByUser.statement, err = r.db.PrepareStatement(statementsList.getNotesByUser.query)
	if err != nil {
		return statements{}, err
	}

	return statementsList, nil
}
