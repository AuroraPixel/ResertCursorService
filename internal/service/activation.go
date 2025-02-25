package service

import (
	"time"

	"github.com/ResertCursorService/internal/models"
	"github.com/ResertCursorService/internal/repository/postgres"
	"github.com/ResertCursorService/pkg/jwt"
	"github.com/ResertCursorService/pkg/utils"
)

type ActivationService struct {
	repo *postgres.ActivationCodeRepository
}

func NewActivationService(repo *postgres.ActivationCodeRepository) *ActivationService {
	return &ActivationService{repo: repo}
}

func (s *ActivationService) CreateActivationCode(duration int, maxAccounts int) (*models.ActivationCode, error) {
	code := &models.ActivationCode{
		Code:        utils.GenerateRandomString(18),
		ExpiresAt:   time.Now().AddDate(0, 0, duration),
		MaxAccounts: maxAccounts,
		Status:      "enabled",
	}

	if err := s.repo.Create(code); err != nil {
		return nil, err
	}

	return code, nil
}

// GetActivationCode 获取激活码信息，包括所有关联的账号
func (s *ActivationService) GetActivationCode(id uint) (*models.ActivationCode, error) {
	return s.repo.FindByID(id)
}

func (s *ActivationService) GetActivationCodeByCode(code string) (*models.ActivationCode, error) {
	return s.repo.FindByCode(code)
}

func (s *ActivationService) ListActivationCodes(page, pageSize int) (*postgres.PaginatedResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.List(page, pageSize)
}

func (s *ActivationService) UpdateActivationCodeStatus(id uint, status string) error {
	return s.repo.UpdateStatus(id, status)
}

func (s *ActivationService) AddAccount(code *models.ActivationCode, account *models.CursorAccount) error {
	return s.repo.AddAccount(code, account)
}

// GenerateAppToken 生成 app 端使用的 token
func (s *ActivationService) GenerateAppToken(codeID uint) (string, error) {
	return jwt.GenerateAppToken(codeID)
}
