package service

import (
	"context"
	"fmt"

	ur "github.com/ObscuraNote/api-general/internal/users/repository"
	"github.com/ObscuraNote/api-general/internal/utils"
	"github.com/philippe-berto/logger"
)

var _ UserService = (*Service)(nil)

type (
	UserService interface {
		CreateUser(userAddress, password string) error
		GetUserId(userAddress, password string) (int64, error)
		CheckUserExists(userAddress, password string) (bool, error)
		UpdatePassword(userAddress, password, newPassword string) error
		DeleteUser(userAddress, password string) (bool, error)
	}

	Service struct {
		ctx  context.Context
		repo ur.UsersRepository
		log  *logger.Logger
	}
)

func New(ctx context.Context, repo ur.UsersRepository) *Service {
	s := &Service{
		ctx:  ctx,
		repo: repo,
		log:  logger.New(ctx),
	}

	return s
}

func (s *Service) CreateUser(userAddress, password string) error {
	if err := s.repo.CreateUser(userAddress, password); err != nil {
		s.log.WithFields(logger.Fields{"error": err.Error(), "component": "user service", "function": "CreateUser"}).
			Error(utils.ErrDatabase)

		return err
	}

	return nil
}

func (s *Service) GetUserId(userAddress, password string) (int64, error) {
	userId, err := s.repo.GetUserId(userAddress, password)
	if err != nil {
		s.log.WithFields(logger.Fields{"error": err.Error(), "component": "user service", "function": "GetUserId"}).
			Error(utils.ErrDatabase)

		return userId, err
	}

	return userId, nil
}

func (s *Service) CheckUserExists(userAddress, password string) (bool, error) {
	exists, err := s.repo.CheckUserExists(userAddress, password)
	if err != nil {
		s.log.WithFields(logger.Fields{"error": err.Error(), "component": "user service", "function": "CheckUserExists"}).
			Error(utils.ErrDatabase)

		return exists, err
	}

	return exists, nil
}

func (s *Service) UpdatePassword(userAddress, currentPassword, newPassword string) error {
	userId, err := s.repo.GetUserId(userAddress, currentPassword)
	if err != nil {
		s.log.WithFields(logger.Fields{"error": err.Error(), "component": "user service", "function": "UpdatePassword"}).
			Error(utils.ErrDatabase)
		return err
	}

	if userId > 0 {
		return s.repo.UpdatePassword(userId, newPassword)
	} else {
		s.log.WithFields(logger.Fields{"error": err.Error(), "component": "user service", "function": "UpdatePassword"}).
			Error(utils.UserNotFound)

		return fmt.Errorf(utils.UserNotFound)
	}
}

func (s *Service) DeleteUser(userAddress, password string) (bool, error) {
	userId, err := s.repo.GetUserId(userAddress, password)
	if err != nil {
		s.log.WithFields(logger.Fields{"error": err.Error(), "component": "user service", "function": "DeleteUser"}).
			Error(utils.ErrDatabase)
		return false, err
	}

	if userId > 0 {
		return s.repo.DeleteUser(userId)
	} else {
		s.log.WithFields(logger.Fields{"error": err.Error(), "component": "user service", "function": "DeleteUser"}).
			Error(utils.UserNotFound)

		return false, fmt.Errorf(utils.UserNotFound)
	}
}
