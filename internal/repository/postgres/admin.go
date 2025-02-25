package postgres

import (
	"github.com/ResertCursorService/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) FindByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepository) Create(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.Admin{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	return r.db.Create(&admin).Error
}

func (r *AdminRepository) VerifyPassword(admin *models.Admin, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password))
	return err == nil
}
