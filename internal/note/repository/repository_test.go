package repository

import (
	"context"
	"testing"

	"github.com/philippe-berto/database/postgresdb"
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
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
	t.Run("GetKeysByUser", func(t *testing.T) {
		userAddress := "user1"
		keys, err := repo.GetKeysByUser(ctx, userAddress)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(keys) == 0 {
			t.Errorf("expected keys, got empty slice")
		}
	})
	t.Run("GetKeysByUser_Error", func(t *testing.T) {
		userAddress := "user2"
		keys, err := repo.GetKeysByUser(ctx, userAddress)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(keys) != 0 {
			t.Errorf("expected empty slice, got %v", keys)
		}
	})

}
