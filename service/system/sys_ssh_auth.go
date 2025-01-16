package system

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/system"
)

type SshAuthService struct {
}

func (s *SshAuthService) CreateSshAuth(ctx context.Context, sshAuth system.SysSshAuth) (err error) {
	if global.OPS_DB.WithContext(ctx).First(&system.SysSshAuth{}).RowsAffected > 0 {
		return errors.New("一个项目只能设置一个ssh用户认证")
	}
	return global.OPS_DB.WithContext(ctx).Create(&sshAuth).Error
}

func (s *SshAuthService) UpdateSshAuth(ctx context.Context, sshAuth system.SysSshAuth) (err error) {
	var old system.SysSshAuth

	updateField := []string{
		"password",
		"private_key",
		"private_key_passphrase",
		"public_key",
		"use_pass",
	}
	if errors.Is(global.OPS_DB.WithContext(ctx).Where("id = ?", sshAuth.ID).First(&old).Error, gorm.ErrRecordNotFound) {
		return errors.New("记录不存在")
	}

	return global.OPS_DB.WithContext(ctx).Model(&old).Select(updateField).Updates(&sshAuth).Error
}

func (s *SshAuthService) DeleteSshAuth(ctx context.Context, id int) (err error) {
	return global.OPS_DB.WithContext(ctx).Where("id = ?", id).Unscoped().Delete(&system.SysSshAuth{}).Error
}

func (s *SshAuthService) GetSshAuthById(ctx context.Context, id int) (result system.SysSshAuth, err error) {
	err = global.OPS_DB.WithContext(ctx).Where("id = ?", id).First(&result).Error
	return
}

func (s *SshAuthService) GetSshAuthList(ctx context.Context, info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Model(&system.SysSshAuth{})

	var resultList []system.SysSshAuth
	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Set("skip_project_filter", true).Error
	if err != nil {
		return resultList, total, err
	}

	if err != nil {
		return resultList, total, err
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	err = db.Order(OrderStr).Find(&resultList).Error
	return resultList, total, err
}
