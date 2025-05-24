//go:build tool
// +build tool

package crypter

//go:generate mockgen -package mocks -destination mocks/mockKeys/repository_mock.go github.com/philippe-berto/internal/key/repository Statementer
