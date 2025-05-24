//go:build tool
// +build tool

package crypter

//go:generate mockgen -package mocks -destination mocks/mockKeys/repository_mock.go github.com/ObscuraNote/api-general/internal/key/repository Statementer
