package service

import (
	"context"

	"github.com/philippe-berto/logger"
)

type (
	UsersRepository interface {
		CreateUser(userAddress, password string) error
		GetUserId(userAddress, password string) (int64, error)
		CheckUserExists(userAddress, password string) (bool, error)
		UpdatePassword(userId int64, password string) error
		DeleteUser(userAddress string) (bool, error)
	}

	Service struct {
		ctx  context.Context
		repo UsersRepository
		log  *logger.Logger
	}
)

func New(ctx context.Context, repo UsersRepository) *Service {
	s := &Service{
		ctx:  ctx,
		repo: repo,
		log:  logger.New(ctx),
	}

	return s
}

func (s *Service) CreateUser(userAddress, password string) error {
	return s.repo.CreateUser(userAddress, password)
}

func (s *Service) GetUserId(userAddress, password string) (int64, error) {
	return s.repo.GetUserId(userAddress, password)
}

func (s *Service) CheckUserExists(userAddress, password string) (bool, error) {
	return s.repo.CheckUserExists(userAddress, password)
}

func (s *Service) UpdatePassword(userAddress, password string) error {
	id, err := s.repo.GetUserId(userAddress, password)
	if err != nil {
		s.log.WithFields(logger.Fields{"error": err.Error()}).Error("Failed to check user existence")
		return err
	}
	if id > 0 {
		return s.repo.UpdatePassword(id, password)
	}
	return nil
}

func (s *Service) DeleteUser(userAddress string) (bool, error) {
	return s.repo.DeleteUser(userAddress)
}
