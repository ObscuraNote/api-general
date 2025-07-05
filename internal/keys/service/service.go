package service

import (
	"context"
	"fmt"

	"github.com/ObscuraNote/api-general/internal/keys/dto"
	r "github.com/ObscuraNote/api-general/internal/keys/repository"
	u "github.com/ObscuraNote/api-general/internal/users/service"
	"github.com/ObscuraNote/api-general/internal/utils"
	"github.com/philippe-berto/logger"
)

var _ KeysService = (*Service)(nil)

type (
	KeysService interface {
		AddKey(note dto.KeyImput) (*dto.KeyOutput, error)
		GetKeysByUser(ctx context.Context, auth dto.AuthInput) ([]dto.KeyOutput, error)
		DeleteKey(keyId string, auth dto.AuthInput) error
	}

	Service struct {
		ctx context.Context
		r   r.KeysRepository
		us  u.UserService
		log *logger.Logger
	}
)

func New(ctx context.Context, log logger.Logger, repo r.KeysRepository, us u.UserService) Service {
	return Service{
		ctx: ctx,
		log: &log,
		r:   repo,
		us:  us,
	}
}

func (s *Service) AddKey(note dto.KeyImput) (*dto.KeyOutput, error) {
	userId, err := s.getUserId(note.UserAddress, note.Password)
	if err != nil {
		return nil, err
	}

	createdKey, err := s.r.AddKey(userId, note)
	if err != nil {
		s.log.WithFields(logger.Fields{"error": err.Error(), "component": "keys_service", "function": "AddKey"}).Error(utils.ErrDatabase)
		return nil, fmt.Errorf(utils.ErrDatabase)
	}

	return createdKey, nil
}

func (s *Service) GetKeysByUser(ctx context.Context, auth dto.AuthInput) ([]dto.KeyOutput, error) {
	userId, err := s.getUserId(auth.UserAddress, auth.Password)
	if err != nil {
		return nil, err
	}
	keys, err := s.r.GetKeysByUser(userId)
	if err != nil {
		s.log.WithFields(logger.Fields{"error": err.Error(), "component": "keys_service", "function": "GetKeysByUser"}).Error(utils.ErrDatabase)
		return nil, fmt.Errorf(utils.ErrDatabase)
	}

	return keys, nil
}

func (s *Service) DeleteKey(keyId string, auth dto.AuthInput) error {
	exists, err := s.us.CheckUserExists(auth.UserAddress, auth.Password)
	if err != nil {
		return fmt.Errorf(utils.ErrUnauthorized)
	}

	if exists {
		err = s.r.DeleteKey(keyId)
		if err != nil {
			s.log.WithFields(logger.Fields{"error": err.Error(), "component": "keys_service", "function": "DeleteKey"}).Error(utils.ErrDatabase)
			return fmt.Errorf(utils.ErrDatabase)
		}
	}

	return nil
}

func (s *Service) getUserId(userAddress, password string) (int64, error) {
	userId, err := s.us.GetUserId(userAddress, password)
	if err != nil {
		s.log.WithFields(logger.Fields{"error": err.Error(), "component": "keys_service", "function": "getUserId"}).Error(utils.ErrDatabase)
		return 0, fmt.Errorf(utils.ErrDatabase)
	}
	if userId <= 0 {
		return 0, fmt.Errorf(utils.ErrUnauthorized)
	}
	return userId, nil
}
