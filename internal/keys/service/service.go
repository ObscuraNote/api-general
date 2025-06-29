package service

import (
	"context"
	"fmt"

	"github.com/ObscuraNote/api-general/internal/keys/dto"
	"github.com/philippe-berto/logger"
)

const (
	ErrDatabase         = "DATABASE_ERROR"
	ErrWrongCredentials = "WRONG_CREDENTIALS"
)

type (
	KeysRepository interface {
		AddKey(userId int64, note dto.KeyImput) error
		GetKeysByUser(userId int64) ([]dto.KeyOutput, error)
		DeleteKey(id string) error
	}

	UserService interface {
		CreateUser(userAddress, password string) error
		GetUserId(userAddress, password string) (int64, error)
		CheckUserExists(userAddress, password string) (bool, error)
		UpdatePassword(userAddress, password string) error
		DeleteUser(userAddress string) (bool, error)
	}

	Service struct {
		ctx context.Context
		r   KeysRepository
		us  UserService
		log *logger.Logger
	}
)

func New(ctx context.Context, log logger.Logger, repo KeysRepository, us UserService) Service {
	return Service{
		ctx: ctx,
		log: &log,
		r:   repo,
		us:  us,
	}
}

func (s *Service) AddKey(note dto.KeyImput) error {
	userId, err := s.getUserId(note.UserAddress, note.Password)
	if err != nil {
		return err
	}

	err = s.r.AddKey(userId, note)
	if err != nil {
		s.log.WithFields(logger.Fields{"error": err.Error(), "component": "Note Service", "function": "AddKey"}).Error(ErrDatabase)
		return fmt.Errorf(ErrDatabase)
	}

	return nil
}

func (s *Service) GetKeysByUser(ctx context.Context, auth dto.AuthInput) ([]dto.KeyOutput, error) {
	userId, err := s.getUserId(auth.UserAddress, auth.Password)
	if err != nil {
		return nil, err
	}
	keys, err := s.r.GetKeysByUser(userId)
	if err != nil {
		s.log.WithFields(logger.Fields{"error": err.Error(), "component": "Note Service", "function": "GetKeysByUser"}).Error(ErrDatabase)
		return nil, fmt.Errorf(ErrDatabase)
	}

	return keys, nil
}

func (s *Service) DeleteKey(input dto.DeleteKeyInput) error {
	exists, err := s.us.CheckUserExists(input.UserAddress, input.Password)
	if err != nil {
		return err
	}

	if exists {
		err = s.r.DeleteKey(input.ID)
		if err != nil {
			s.log.WithFields(logger.Fields{"error": err.Error(), "component": "Note Service", "function": "DeleteKey"}).Error(ErrDatabase)
			return fmt.Errorf(ErrDatabase)
		}
	}

	return nil
}

func (s *Service) getUserId(userAddress, password string) (int64, error) {
	userId, err := s.us.GetUserId(userAddress, password)
	if err != nil {
		s.log.WithFields(logger.Fields{"error": err.Error(), "component": "Note Service", "function": "getUserId"}).Error(ErrDatabase)
		return 0, fmt.Errorf(ErrDatabase)
	}
	if userId <= 0 {
		return 0, fmt.Errorf(ErrWrongCredentials)
	}
	return userId, nil
}
