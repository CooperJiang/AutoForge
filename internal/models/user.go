package models

import (
	"auto-forge/pkg/common"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	BaseModel
	Username   string `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password   string `gorm:"size:100" json:"-"` // OAuth2用户可以为空
	Email      string `gorm:"size:100;uniqueIndex" json:"email"`
	Avatar     string `gorm:"size:255" json:"avatar"`
	Bio        string `gorm:"size:500" json:"bio"`
	Status     int    `gorm:"default:1" json:"status"`    // 1:正常 2:禁用 3:删除
	Role       int    `gorm:"default:3" json:"role"`      // 1:超级管理员 2:管理员 3:普通用户
	Provider   string `gorm:"size:50;index" json:"provider"` // 登录提供商: local, oauth2_linuxdo
	ExternalID string `gorm:"size:100;index" json:"external_id"` // 第三方平台用户ID
	TrustLevel int    `gorm:"default:0" json:"trust_level"` // Linux.do信任等级 (0-4)
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的BeforeCreate
	if err := u.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	if u.Status == 0 {
		u.Status = common.UserStatusNormal
	}
	if u.Role == 0 {
		u.Role = common.UserRoleUser
	}
	return nil
}

// IsAdmin 是否是管理员
func (u *User) IsAdmin() bool {
	return u.Role == common.UserRoleAdmin || u.Role == common.UserRoleSuperAdmin
}

// IsSuperAdmin 是否是超级管理员
func (u *User) IsSuperAdmin() bool {
	return u.Role == common.UserRoleSuperAdmin
}

// IsNormal 是否是正常状态
func (u *User) IsNormal() bool {
	return u.Status == common.UserStatusNormal
}

// IsDisabled 是否被禁用
func (u *User) IsDisabled() bool {
	return u.Status == common.UserStatusDisabled
}

// IsDeleted 是否被删除
func (u *User) IsDeleted() bool {
	return u.Status == common.UserStatusDeleted
}
