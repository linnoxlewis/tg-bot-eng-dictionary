package manager

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"linnoxlewis/tg-bot-eng-dictionary/internal/api"
	"linnoxlewis/tg-bot-eng-dictionary/internal/models"
	"linnoxlewis/tg-bot-eng-dictionary/internal/repository"
)

type Manager struct {
	logger           *logrus.Logger
	userRepo         *repository.UserRepo
	translateService api.TranslateInterface
}

var (
	errUserAlreadyExist = errors.New("user already exist")
	errUserNotFound     = errors.New("user not found")
	errCreateUser       = errors.New("creating user error")
	errUpdateStatus     = errors.New("include learn mode error")
)

func NewManager(logger *logrus.Logger,
	userRepo *repository.UserRepo,
	translateService api.TranslateInterface) *Manager {
	return &Manager{
		logger:           logger,
		userRepo:         userRepo,
		translateService: translateService,
	}
}

func (m *Manager) CreateUser(ctx context.Context, user *models.User) error {
	existUser, err := m.userRepo.GetByTgId(ctx, user.ChatId)
	if err != nil {
		m.logger.Println(err)
		return errCreateUser
	}

	if !existUser.Empty() {
		m.logger.Println(errUserAlreadyExist)
		return nil
	}

	if err := m.userRepo.CreateUser(ctx, user); err != nil {
		m.logger.Println(err)
		return errCreateUser
	}

	return err
}

func (m *Manager) TranslateWord(word string) (*models.Word, error) {
	trl, err := m.translateService.GetTranslate(word)
	if err != nil {
		m.logger.Error(err)
		return nil, errors.New("can not translate this word")
	}
	if trl == nil {
		return nil, errors.New("translate not found")
	}

	return trl, nil
}

func (m *Manager) UpdateLearnStatus(ctx context.Context, chatId int64) error {
	user, err := m.userRepo.GetByTgId(ctx, chatId)
	if err != nil {
		m.logger.Println(err)
		return err
	}
	if user.Empty() {
		return errUserNotFound
	}

	if err := m.userRepo.UpdateLearnStatusUser(ctx, user); err != nil {
		m.logger.Println(err)
		return errUpdateStatus
	}

	return nil
}
