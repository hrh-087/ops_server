package system

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/system"
)

type AssetsRedisService struct {
}

func (a *AssetsRedisService) CreateRedis(ctx context.Context, redis system.SysAssetsRedis) (err error) {
	if !errors.Is(global.OPS_DB.WithContext(ctx).Where("host = ?", redis.Host).First(&system.SysAssetsRedis{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("记录已存在,请检查后重新添加")
	}
	err = global.OPS_DB.WithContext(ctx).Create(&redis).Error
	return
}

func (a *AssetsRedisService) UpdateRedis(ctx context.Context, redis system.SysAssetsRedis) (err error) {
	var old system.SysAssetsRedis

	updateField := []string{
		"platform_id",
		"host",
		"port",
		"name",
		"password",
	}
	if errors.Is(global.OPS_DB.WithContext(ctx).Where("id = ?", redis.ID).First(&old).Error, gorm.ErrRecordNotFound) {
		return errors.New("记录不存在")
	}
	err = global.OPS_DB.Debug().WithContext(ctx).Model(&old).Select(updateField).Updates(redis).Error
	return
}

func (a *AssetsRedisService) DeleteRedis(ctx context.Context, id int) (err error) {
	return global.OPS_DB.WithContext(ctx).Where("id = ?", id).Delete(&system.SysAssetsRedis{}).Error
}

func (a *AssetsRedisService) GetAssetsRedisById(ctx context.Context, id int) (result system.SysAssetsRedis, err error) {
	err = global.OPS_DB.WithContext(ctx).Preload("Platform").Preload("SysProject").First(&result, "id = ?", id).Error
	return
}

func (a *AssetsRedisService) GetAssetsRedisAll(ctx context.Context, id uint) (result []system.SysAssetsRedis, err error) {
	err = global.OPS_DB.WithContext(ctx).Where("platform_id = ?", id).Find(&result).Error
	return
}

func (a *AssetsRedisService) GetAssetsRedisList(ctx context.Context, info request.PageInfo, server request.NameAndPlatformSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Model(&system.SysAssetsRedis{})

	var resultList []system.SysAssetsRedis

	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Set("skip_project_filter", true).Error
	if err != nil {
		return resultList, total, err
	}

	if server.PlatformId != 0 {
		db = db.Where("platform_id = ?", server.PlatformId)
	}

	if server.Name != "" {
		db = db.Where("name like ?", "%"+server.Name+"%")
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	err = db.Preload("SysProject").Preload("Platform").Order(OrderStr).Find(&resultList).Error
	return resultList, total, err
}
