//go:build integration
// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/ObscuraNote/api-general/internal/keys/dto"
	"github.com/philippe-berto/database/postgresdb"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()
var cfg = postgresdb.Config{
	Host:         "localhost",
	Name:         "crypter",
	Password:     "password",
	User:         "user",
	Port:         5432,
	Driver:       "postgres",
	RunMigration: true,
}

func TestRepository(t *testing.T) {
	db, err := postgresdb.New(ctx, cfg, false, "file://../../../migrations")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	db.GetClient().Exec("TRUNCATE TABLE keys;")
	defer db.Close()
	defer db.GetClient().Exec("TRUNCATE TABLE keys;")

	db.GetClient().Exec("TRUNCATE TABLE users;")
	defer db.GetClient().Exec("TRUNCATE TABLE users;")

	db.GetClient().Exec(`
		INSERT INTO users (user_address, password)
		VALUES ('1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef', 'abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890'),
		       ('2222222222222222222222222222222222222222222222222222222222222222', 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb');
	`)

	repo := New(ctx, db)

	// Test user ID (should match an existing user in the database)
	var userId int64
	err = db.GetClient().QueryRow(`
		SELECT id FROM users WHERE user_address = '1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef';
	`).Scan(&userId)
	if err != nil {
		t.Fatalf("failed to get user id: %v", err)
	}

	t.Run("AddKey", func(t *testing.T) {
		note := dto.KeyImput{
			UserAddress:   "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			EncryptedKey:  []byte("key"),
			KeyIV:         []byte("key"),
			EncryptedData: []byte("enc"),
			DataIV:        []byte("iv"),
		}
		err := repo.AddKey(userId, note)
		assert.NoError(t, err)
	})

	t.Run("GetKeysByUser", func(t *testing.T) {
		keys, err := repo.GetKeysByUser(userId)
		assert.NoError(t, err)
		assert.Len(t, keys, 1)
	})

	t.Run("DeleteKey", func(t *testing.T) {
		note := dto.KeyImput{
			UserAddress:   "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			EncryptedKey:  []byte("key"),
			KeyIV:         []byte("key"),
			EncryptedData: []byte("enc"),
			DataIV:        []byte("iv"),
		}
		err := repo.AddKey(userId, note)
		assert.NoError(t, err)

		keys, err := repo.GetKeysByUser(userId)
		assert.NoError(t, err)
		assert.Len(t, keys, 2)

		err = repo.DeleteKey(keys[0].ID)
		assert.NoError(t, err)

		keys, err = repo.GetKeysByUser(userId)
		assert.NoError(t, err)
		assert.Len(t, keys, 1)

	})

	t.Run("GetKeysByUser_Error", func(t *testing.T) {
		nonExistentUserId := int64(99999)
		keys, err := repo.GetKeysByUser(nonExistentUserId)
		assert.NoError(t, err)
		assert.Len(t, keys, 0)
	})

}
