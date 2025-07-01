//go:build integration
// +build integration

package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/philippe-berto/database/postgresdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testUserAddress     = "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	testPassword        = "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	testPassword2       = "fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321"
	testWrongPassword   = "1111111111111111111111111111111111111111111111111111111111111111"
	testNonExistentAddr = "9999999999999999999999999999999999999999999999999999999999999999"
)

type RepositoryTestSuite struct {
	suite.Suite
	db   *postgresdb.Client
	repo *Repository
	ctx  context.Context
}

func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	var cfg = postgresdb.Config{
		Host:         "localhost",
		Name:         "crypter",
		Password:     "password",
		User:         "user",
		Port:         5432,
		Driver:       "postgres",
		RunMigration: true,
	}

	// Connect to test database
	db, err := postgresdb.New(suite.ctx, cfg, false, "file://../../../migrations")
	require.NoError(suite.T(), err)

	suite.db = db
	suite.repo = New(suite.ctx, db)
}

func (suite *RepositoryTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

func (suite *RepositoryTestSuite) SetupTest() {
	// Clean users table before each test
	_, err := suite.db.GetClient().Exec("DELETE FROM users")
	require.NoError(suite.T(), err)
}

func (suite *RepositoryTestSuite) TestCreateUser() {
	err := suite.repo.CreateUser(testUserAddress, testPassword)
	assert.NoError(suite.T(), err)

	// Verify user was created
	exists, err := suite.repo.CheckUserExists(testUserAddress, testPassword)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *RepositoryTestSuite) TestCreateUser_DuplicateAddress() {
	// Create first user
	err := suite.repo.CreateUser(testUserAddress, testPassword)
	assert.NoError(suite.T(), err)

	// Try to create user with same address
	err = suite.repo.CreateUser(testUserAddress, testPassword2)
	assert.Error(suite.T(), err)
}

func (suite *RepositoryTestSuite) TestCheckUserExists() {
	// Check non-existent user
	exists, err := suite.repo.CheckUserExists(testUserAddress, testPassword)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)

	// Create user
	err = suite.repo.CreateUser(testUserAddress, testPassword)
	require.NoError(suite.T(), err)

	// Check existing user with correct credentials
	exists, err = suite.repo.CheckUserExists(testUserAddress, testPassword)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)

	// Check existing user with wrong password
	exists, err = suite.repo.CheckUserExists(testUserAddress, testWrongPassword)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *RepositoryTestSuite) TestGetUserId() {
	// Create user
	err := suite.repo.CreateUser(testUserAddress, testPassword)
	require.NoError(suite.T(), err)

	// Get user ID with correct credentials
	id, err := suite.repo.GetUserId(testUserAddress, testPassword)
	assert.NoError(suite.T(), err)
	assert.Greater(suite.T(), id, int64(0))

	// Try with wrong password
	_, err = suite.repo.GetUserId(testUserAddress, testWrongPassword)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), sql.ErrNoRows, err)

	// Try with non-existent user
	_, err = suite.repo.GetUserId(testNonExistentAddr, testPassword)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), sql.ErrNoRows, err)
}

func (suite *RepositoryTestSuite) TestUpdatePassword() {
	// Create user
	err := suite.repo.CreateUser(testUserAddress, testPassword)
	require.NoError(suite.T(), err)

	// Get user ID with correct credentials
	id, err := suite.repo.GetUserId(testUserAddress, testPassword)
	assert.NoError(suite.T(), err)
	assert.Greater(suite.T(), id, int64(0))

	// Update password using user ID
	err = suite.repo.UpdatePassword(id, testPassword2)
	assert.NoError(suite.T(), err)

	// Verify old password no longer works
	_, err = suite.repo.GetUserId(testUserAddress, testPassword)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), sql.ErrNoRows, err)

	// Verify new password works
	newId, err := suite.repo.GetUserId(testUserAddress, testPassword2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), id, newId)
}

func (suite *RepositoryTestSuite) TestUpdatePassword_NonExistentUser() {
	// Try to update password for non-existent user ID
	nonExistentUserId := int64(99999)
	err := suite.repo.UpdatePassword(nonExistentUserId, testPassword)
	assert.NoError(suite.T(), err) // UPDATE returns no error even if no rows affected
}

func (suite *RepositoryTestSuite) TestDeleteUser() {
	// Create user
	err := suite.repo.CreateUser(testUserAddress, testPassword)
	require.NoError(suite.T(), err)

	// Verify user exists
	exists, err := suite.repo.CheckUserExists(testUserAddress, testPassword)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)

	// Get user ID for deletion
	userId, err := suite.repo.GetUserId(testUserAddress, testPassword)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), userId, int64(0))

	// Delete user
	deleted, err := suite.repo.DeleteUser(userId)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), deleted) // Verify that a row was actually deleted

	// Verify user no longer exists
	exists, err = suite.repo.CheckUserExists(testUserAddress, testPassword)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *RepositoryTestSuite) TestDeleteUser_NonExistent() {
	// Try to delete non-existent user
	nonExistentUserId := int64(99999)
	deleted, err := suite.repo.DeleteUser(nonExistentUserId)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), deleted) // Verify that no rows were affected
}

func (suite *RepositoryTestSuite) TestMultipleUsers() {
	users := []struct {
		address  string
		password string
	}{
		{"1111111111111111111111111111111111111111111111111111111111111111", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{"2222222222222222222222222222222222222222222222222222222222222222", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
		{"3333333333333333333333333333333333333333333333333333333333333333", "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"},
	}

	// Create multiple users
	for _, user := range users {
		err := suite.repo.CreateUser(user.address, user.password)
		assert.NoError(suite.T(), err)
	}

	// Verify all users exist and can authenticate
	for _, user := range users {
		exists, err := suite.repo.CheckUserExists(user.address, user.password)
		assert.NoError(suite.T(), err)
		assert.True(suite.T(), exists)

		id, err := suite.repo.GetUserId(user.address, user.password)
		assert.NoError(suite.T(), err)
		assert.Greater(suite.T(), id, int64(0))
	}

	// Delete one user
	// Get user ID for deletion
	userIdToDelete, err := suite.repo.GetUserId(users[1].address, users[1].password)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), userIdToDelete, int64(0))

	deleted, err := suite.repo.DeleteUser(userIdToDelete)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), deleted) // Verify that a row was actually deleted

	// Verify deleted user no longer exists
	exists, err := suite.repo.CheckUserExists(users[1].address, users[1].password)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)

	// Verify other users still exist
	for i, user := range users {
		if i == 1 {
			continue // Skip deleted user
		}
		exists, err := suite.repo.CheckUserExists(user.address, user.password)
		assert.NoError(suite.T(), err)
		assert.True(suite.T(), exists)
	}
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
