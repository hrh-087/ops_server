package system

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/system"
	"ops-server/utils/game"
	"strings"
)

type GamePlatformService struct {
}

func (g *GamePlatformService) CreatePlatform(ctx context.Context, platform system.SysGamePlatform) (err error) {
	if !errors.Is(global.OPS_DB.WithContext(ctx).Where("platform_code = ?", platform.PlatformCode).First(&system.SysGamePlatform{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("项目已存在该渠道")
	}

	platform.PlatformCode = strings.TrimSpace(platform.PlatformCode)

	// 如果不以/结尾则添加根号
	if !strings.HasSuffix(platform.ImageUri, "/") {
		platform.ImageUri += "/"
	}

	return global.OPS_DB.WithContext(ctx).Create(&platform).Error
}

func (g *GamePlatformService) UpdatePlatform(ctx context.Context, platform system.SysGamePlatform) (err error) {
	var old system.SysGamePlatform
	err = global.OPS_DB.WithContext(ctx).Where("id = ?", platform.ID).First(&old).Error
	if err != nil {
		return
	}

	// 如果不以/结尾则添加根号
	if !strings.HasSuffix(old.ImageUri, "/") {
		platform.ImageUri += "/"
	}
	return global.OPS_DB.WithContext(ctx).Model(&old).Updates(platform).Error
}

func (g *GamePlatformService) DeletePlatform(ctx context.Context, id int) (err error) {
	return global.OPS_DB.WithContext(ctx).Where("id = ?", id).Delete(&system.SysGamePlatform{}).Error
}

func (g *GamePlatformService) GetPlatformById(ctx context.Context, id int) (platform system.SysGamePlatform, err error) {
	err = global.OPS_DB.WithContext(ctx).Preload("SysProject").First(&platform, "id = ?", id).Error
	return
}

func (g *GamePlatformService) GetPlatformAll(ctx context.Context) (result []system.SysGamePlatform, err error) {
	err = global.OPS_DB.WithContext(ctx).Find(&result).Error
	return
}

func (g *GamePlatformService) GetPlatformList(ctx context.Context, platform system.SysGamePlatform, info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Model(&system.SysGamePlatform{})

	var resultList []system.SysGamePlatform

	if platform.PlatformCode != "" {
		db = db.Where("platform_code LIKE ?", "%"+platform.PlatformCode+"%")
	}

	if platform.PlatformName != "" {
		db = db.Where("platform_name LIKE ?", "%"+platform.PlatformName+"%")
	}

	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Set("skip_project_filter", true).Error

	if err != nil {
		return resultList, total, err
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	err = db.Preload("SysProject").Order(OrderStr).Find(&resultList).Error
	return resultList, total, err
}

func (g GamePlatformService) KickGameServer(ctx *gin.Context, serverId int) (err error) {
	err = game.KickPlayer(ctx, serverId)
	return
}
