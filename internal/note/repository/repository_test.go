package repository

import (
	"context"
	"testing"

	"github.com/philippe-berto/database/postgresdb"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()
var cfg = postgresdb.Config{
	Host:         "localhost",
	Name:         "test",
	Password:     "postgres",
	User:         "postgres",
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

	repo := New(db)

	t.Run("AddKey", func(t *testing.T) {
		note := keyImput{
			UserAddress:   "user1",
			Key:           []byte("key"),
			EncryptedData: []byte("enc"),
			IV:            []byte("iv"),
		}
		err := repo.AddKey(ctx, note)
		assert.NoError(t, err)
	})

	t.Run("GetKeysByUser", func(t *testing.T) {
		userAddress := "user1"
		keys, err := repo.GetKeysByUser(ctx, userAddress)
		assert.NoError(t, err)
		assert.Len(t, keys, 1)
	})

	t.Run("DeleteKey", func(t *testing.T) {
		note := keyImput{
			UserAddress:   "user1",
			Key:           []byte("key"),
			EncryptedData: []byte("enc"),
			IV:            []byte("iv"),
		}
		err := repo.AddKey(ctx, note)
		assert.NoError(t, err)

		userAddress := "user1"
		keys, err := repo.GetKeysByUser(ctx, userAddress)
		assert.NoError(t, err)
		assert.Len(t, keys, 2)

		err = repo.DeleteKey(ctx, keys[0].ID)
		assert.NoError(t, err)

		keys, err = repo.GetKeysByUser(ctx, userAddress)
		assert.NoError(t, err)
		assert.Len(t, keys, 1)

	})

	t.Run("GetKeysByUser_Error", func(t *testing.T) {
		userAddress := "user2"
		keys, err := repo.GetKeysByUser(ctx, userAddress)
		assert.NoError(t, err)
		assert.Len(t, keys, 0)
	})

}
