package models

import (
	"time"

	"gorm.io/gorm"
)

// Admin 管理员模型
type Admin struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	Username     string         `json:"username" gorm:"unique;not null"`
	PasswordHash string         `json:"-" gorm:"not null"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// ActivationCode 激活码模型
type ActivationCode struct {
	ID          uint            `json:"id" gorm:"primarykey"`
	Code        string          `json:"code" gorm:"unique;not null;size:18"`
	ExpiresAt   time.Time       `json:"expiresAt" gorm:"not null"`
	MaxAccounts int             `json:"maxAccounts" gorm:"not null;default:1"`
	Status      string          `json:"status" gorm:"not null;default:'enabled'"` // enabled or disabled
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt  `json:"-" gorm:"index"`
	Accounts    []CursorAccount `json:"accounts"`
}

// CursorAccount Cursor账号模型
type CursorAccount struct {
	ID               uint           `json:"id" gorm:"primarykey"`
	ActivationCode   ActivationCode `json:"-" gorm:"foreignKey:ActivationCodeID"`
	ActivationCodeID uint           `json:"activationCodeId" gorm:"not null"`
	Email            string         `json:"email" gorm:"not null"`
	EmailPassword    string         `json:"emailPassword" gorm:"not null"`  // 邮箱密码
	CursorPassword   string         `json:"cursorPassword" gorm:"not null"` // Cursor密码
	AccessToken      string         `json:"accessToken" gorm:"type:text"`
	RefreshToken     string         `json:"refreshToken" gorm:"type:text"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}
