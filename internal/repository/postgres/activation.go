package postgres

import (
	"time"

	"github.com/ResertCursorService/internal/models"
	"gorm.io/gorm"
)

type ActivationCodeRepository struct {
	db *gorm.DB
}

func NewActivationCodeRepository(db *gorm.DB) *ActivationCodeRepository {
	return &ActivationCodeRepository{db: db}
}

func (r *ActivationCodeRepository) Create(code *models.ActivationCode) error {
	return r.db.Create(code).Error
}

func (r *ActivationCodeRepository) FindByCode(code string) (*models.ActivationCode, error) {
	var activationCode models.ActivationCode
	if err := r.db.Preload("Accounts").Where("code = ?", code).First(&activationCode).Error; err != nil {
		return nil, err
	}
	return &activationCode, nil
}

func (r *ActivationCodeRepository) FindByID(id uint) (*models.ActivationCode, error) {
	var activationCode models.ActivationCode
	if err := r.db.Preload("Accounts").First(&activationCode, id).Error; err != nil {
		return nil, err
	}
	return &activationCode, nil
}

type PaginatedResult struct {
	Items      []models.ActivationCode
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

func (r *ActivationCodeRepository) List(page, pageSize int) (*PaginatedResult, error) {
	var total int64
	if err := r.db.Model(&models.ActivationCode{}).Count(&total).Error; err != nil {
		return nil, err
	}

	var codes []models.ActivationCode
	offset := (page - 1) * pageSize
	if err := r.db.Preload("Accounts").Order("created_at desc").Offset(offset).Limit(pageSize).Find(&codes).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PaginatedResult{
		Items:      codes,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *ActivationCodeRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.ActivationCode{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ActivationCodeRepository) AddAccount(code *models.ActivationCode, account *models.CursorAccount) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 检查是否超过最大账户数
		var count int64
		if err := tx.Model(&models.CursorAccount{}).Where("activation_code_id = ?", code.ID).Count(&count).Error; err != nil {
			return err
		}

		if int(count) >= code.MaxAccounts {
			return gorm.ErrInvalidData
		}

		// 检查激活码是否过期
		if code.ExpiresAt.Before(time.Now()) {
			return gorm.ErrInvalidData
		}

		// 添加账户
		account.ActivationCodeID = code.ID
		if err := tx.Create(account).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetAccountByCodeAndToken 根据激活码和 token 获取账号信息
func (r *ActivationCodeRepository) GetAccountByCodeAndToken(code, token string) (*models.CursorAccount, error) {
	var account models.CursorAccount
	err := r.db.Joins("JOIN activation_codes ON cursor_accounts.activation_code_id = activation_codes.id").
		Where("activation_codes.code = ?", code).
		First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}
