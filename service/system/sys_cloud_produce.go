package system

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/system"
	"strings"
)

type CloudProduceService struct {
}

var CloudProduceServiceApp = new(CloudProduceService)

func (c *CloudProduceService) CreateCloudProduce(ctx context.Context, cloud system.SysCloudProduce) (err error) {

	if !errors.Is(global.OPS_DB.WithContext(ctx).Where("region_id = ?", cloud.RegionId).First(&system.SysCloudProduce{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("项目已存在该地域记录")
	}

	cloud.SecretId = strings.TrimSpace(cloud.SecretId)
	cloud.SecretKey = strings.TrimSpace(cloud.SecretKey)

	return global.OPS_DB.WithContext(ctx).Create(&cloud).Error
}

func (c *CloudProduceService) GetCloudProduceList(ctx context.Context, cloud system.SysCloudProduce, info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Model(&system.SysCloudProduce{})

	var cloudList []system.SysCloudProduce

	if cloud.CloudName != "" {
		db = db.Where("cloud_name LIKE ?", "%"+cloud.CloudName+"%")
	}

	if cloud.RegionId != "" {
		db = db.Where("region_id LIKE ?", "%"+cloud.RegionId+"%")
	}

	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Set("skip_project_filter", true).Error

	if err != nil {
		return cloudList, total, err
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	err = db.Preload("SysProject").Order(OrderStr).Find(&cloudList).Error
	return cloudList, total, err
}

func (c *CloudProduceService) GetCloudProduceById(ctx context.Context, id int) (cloud system.SysCloudProduce, err error) {
	err = global.OPS_DB.WithContext(ctx).Preload("SysProject").First(&cloud, "id = ?", id).Error
	return
}

func (c *CloudProduceService) UpdateCloudProduce(ctx context.Context, cloud system.SysCloudProduce) (err error) {
	var oldCloud system.SysCloudProduce
	err = global.OPS_DB.WithContext(ctx).Where("id = ?", cloud.ID).First(&oldCloud).Error
	if err != nil {
		return
	}
	return global.OPS_DB.WithContext(ctx).Model(&oldCloud).Select("*").Updates(&cloud).Error
}

func (c *CloudProduceService) DeleteCloudProduce(ctx context.Context, id int) (err error) {
	return global.OPS_DB.WithContext(ctx).Where("id = ?", id).Delete(&system.SysCloudProduce{}).Error
}

func (c *CloudProduceService) GetCloudProduceAll(ctx context.Context) (result []system.SysCloudProduce, err error) {
	err = global.OPS_DB.WithContext(ctx).Find(&result).Error
	return
}
