package service

import (
	"errors"

	"github.com/ResertCursorService/internal/repository/postgres"
	"github.com/ResertCursorService/pkg/jwt"
)

type AdminService struct {
	repo *postgres.AdminRepository
}

func NewAdminService(repo *postgres.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) Login(username, password string) (string, error) {
	admin, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("用户名或密码错误")
	}

	if !s.repo.VerifyPassword(admin, password) {
		return "", errors.New("用户名或密码错误")
	}

	// 生成JWT令牌
	token, err := jwt.GenerateToken(admin.ID)
	if err != nil {
		return "", errors.New("生成令牌失败")
	}

	return token, nil
}

func (s *AdminService) CreateDefaultAdmin() error {
	_, err := s.repo.FindByUsername("wang")
	if err == nil {
		return nil // 已存在默认管理员
	}

	return s.repo.Create("wang", "wangtest")
}
