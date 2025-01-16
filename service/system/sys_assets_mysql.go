package system

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/system"
)

type AssetsMysqlService struct {
}

func (a *AssetsMysqlService) CreateMysql(ctx context.Context, mysql system.SysAssetsMysql) (err error) {
	if !errors.Is(global.OPS_DB.WithContext(ctx).Where("host = ?", mysql.Host).First(&system.SysAssetsMysql{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("记录已存在,请检查后重新添加")
	}
	err = global.OPS_DB.WithContext(ctx).Create(&mysql).Error
	return
}

func (a *AssetsMysqlService) UpdateMysql(ctx context.Context, mysql system.SysAssetsMysql) (err error) {
	var old system.SysAssetsMysql

	updateField := []string{
		"platform_id",
		"host",
		"name",
		"port",
		"user",
		"pass",
	}
	if errors.Is(global.OPS_DB.WithContext(ctx).Where("id = ?", mysql.ID).First(&old).Error, gorm.ErrRecordNotFound) {
		return errors.New("记录不存在")
	}
	err = global.OPS_DB.Debug().WithContext(ctx).Model(&old).Select(updateField).Updates(mysql).Error
	return
}

func (a *AssetsMysqlService) DeleteMysql(ctx context.Context, id int) (err error) {
	return global.OPS_DB.WithContext(ctx).Where("id = ?", id).Delete(&system.SysAssetsMysql{}).Error
}

func (a *AssetsMysqlService) GetAssetsMysqlById(ctx context.Context, id int) (result system.SysAssetsMysql, err error) {
	err = global.OPS_DB.WithContext(ctx).Preload("Platform").Preload("SysProject").First(&result, "id = ?", id).Error
	return
}

func (a *AssetsMysqlService) GetAssetsMysqlAll(ctx context.Context, id uint) (result []system.SysAssetsMysql, err error) {
	err = global.OPS_DB.WithContext(ctx).Where("platform_id = ?", id).Find(&result).Error
	return
}

func (a *AssetsMysqlService) GetAssetsMysqlList(ctx context.Context, info request.PageInfo, server request.NameAndPlatformSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Model(&system.SysAssetsMysql{})

	var resultList []system.SysAssetsMysql

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
