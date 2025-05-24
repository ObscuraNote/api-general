package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/philippe-berto/database/postgresdb"
)

type (
	Statementer interface {
		ExecContext(ctx context.Context, args ...any) (sql.Result, error)
		QueryContext(ctx context.Context, args ...any) (*sql.Rows, error)
	}

	Repository struct {
		db        *postgresdb.Client
		statments map[string]*sqlx.Stmt
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
		db:        db,
		statments: map[string]*sqlx.Stmt{},
	}
	if err := r.prepareStatements(); err != nil {
		panic(err)
	}

	return r
}

func (r *Repository) AddKey(ctx context.Context, note *keyImput) error {
	_, err := statementsList.addNote.statement.
		ExecContext(ctx, note.UserAddress, note.Key, note.EncryptedData, note.IV)
	if err != nil {
		log.Println("Error adding note")

		return err
	}

	return nil
}

func (r *Repository) GetKeysByUser(ctx context.Context, userAddress string) ([]keyOutput, error) {
	rows, err := statementsList.getNotesByUser.statement.
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

func (r *Repository) prepareStatements() error {
	var err error

	statementsList.addNote.statement, err = r.db.PrepareStatement(statementsList.addNote.query)
	if err != nil {
		return err
	}

	statementsList.getNotesByUser.statement, err = r.db.PrepareStatement(statementsList.getNotesByUser.query)
	if err != nil {
		return err
	}

	return nil
}
