package system

import (
	"context"
	"errors"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/system"
	"strconv"
	"strings"
)

type AssetsServerService struct{}

var AssetsServerApp = new(AssetsServerService)

func (s *AssetsServerService) CreateAssetsServer(ctx context.Context, server system.SysAssetsServer) (err error) {
	if !errors.Is(global.OPS_DB.WithContext(ctx).Where("pub_ip = ? or private_ip = ?", server.PubIp, server.PrivateIp).First(&system.SysAssetsServer{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("内网ip或外网ip已存在,请检查后重新添加")
	}
	server.UUID = uuid.Must(uuid.NewV4())
	err = global.OPS_DB.WithContext(ctx).Create(&server).Error
	return
}

func (s *AssetsServerService) UpdateAssetsServer(ctx context.Context, server system.SysAssetsServer) (err error) {
	var oldServer system.SysAssetsServer

	updateField := []string{
		"platform_id",
		"private_ip",
		"pub_ip",
		"ssh_port",
		"ServerName",
		"CloudProduceId",
		"InstanceId",
		"ServerType",
	}

	if errors.Is(global.OPS_DB.WithContext(ctx).Where("id = ?", server.ID).First(&oldServer).Error, gorm.ErrRecordNotFound) {
		return errors.New("记录不存在")
	}

	err = global.OPS_DB.Debug().WithContext(ctx).Model(&oldServer).Select(updateField).Updates(server).Error
	return
}

func (s *AssetsServerService) DeleteAssetsServer(ctx context.Context, id int) (err error) {
	return global.OPS_DB.WithContext(ctx).Where("id = ?", id).Unscoped().Delete(&system.SysAssetsServer{}).Error
}

func (s *AssetsServerService) GetAssetsServerList(ctx context.Context, info request.PageInfo, server request.NameAndPlatformSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Model(&system.SysAssetsServer{})

	var resultList []system.SysAssetsServer

	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Set("skip_project_filter", true).Error
	if err != nil {
		return resultList, total, err
	}

	if server.PlatformId != 0 {
		db = db.Where("platform_id = ?", server.PlatformId)
	}

	if server.Name != "" {
		db = db.Where("server_name like ?", "%"+server.Name+"%")
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	err = db.Preload("Platform").Preload("Cloud").Order(OrderStr).Find(&resultList).Error
	return resultList, total, err
}

func (s *AssetsServerService) GetAssetsServerById(ctx context.Context, id int) (result system.SysAssetsServer, err error) {
	err = global.OPS_DB.WithContext(ctx).Preload("Platform").Preload("Ports").First(&result, "id = ?", id).Error
	return
}

func (s *AssetsServerService) GetAssetsServerAll(ctx context.Context, id uint) (result []system.SysAssetsServer, err error) {
	if id == 0 {
		err = global.OPS_DB.WithContext(ctx).Find(&result).Error
	} else {
		err = global.OPS_DB.WithContext(ctx).Where("platform_id = ?", id).Find(&result).Error
	}

	return
}

func (s *AssetsServerService) getServerPort(serverId uint, ruleRange string, tx *gorm.DB) (port int64, err error) {
	ports := strings.Split(ruleRange, ",")
	if len(ports) <= 1 {
		return 0, errors.New("端口规则不正确")
	}

	tx.Model(&system.SysAssetsServerPort{}).Select("max(port) as max").Where("server_id = ? and port BETWEEN ? and ?", serverId, ports[0], ports[1]).Pluck("max", &port)
	if port == 0 {
		port, err = strconv.ParseInt(ports[0], 10, 64)
		if err != nil {
			return 0, errors.New("端口解析失败")
		}
	}
	newPort := port + 1

	serverPort := system.SysAssetsServerPort{
		ServerId: serverId,
		Port:     newPort,
	}

	return newPort, tx.Create(&serverPort).Error
}
