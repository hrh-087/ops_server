package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/job/task"
	"ops-server/model/common/request"
	"ops-server/model/system"
	"ops-server/utils"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	var serverPort system.SysAssetsServerPort

	ports := strings.Split(ruleRange, ",")
	if len(ports) <= 1 {
		return 0, errors.New("端口规则不正确")
	}

	err = tx.Debug().Where("server_id = ? and port BETWEEN ? and ?", serverId, ports[0], ports[1]).Order("port desc").First(&serverPort).Error
	if err == gorm.ErrRecordNotFound {
		port, err = strconv.ParseInt(ports[0], 10, 64)
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	} else {
		port = serverPort.Port + 1
	}

	return port, tx.Create(&system.SysAssetsServerPort{ServerId: serverId, Port: port}).Error
}

func (s AssetsServerService) GeneratePrometheusHostConfig(ctx *gin.Context) (err error) {
	var hostList []system.SysAssetsServer
	var prometheusHostConfig []request.PrometheusConfig

	projectId, err := strconv.ParseUint(ctx.GetString("projectId"), 10, 64)
	if err != nil {
		return errors.New("项目id解析失败")
	}

	sshConfig, err := task.GetSSHKey(uint(projectId), global.OPS_CONFIG.Prometheus.Addr, global.OPS_CONFIG.Prometheus.SshPort)
	if err != nil {
		return fmt.Errorf("获取ssh配置失败:%v", err)
	}

	sshClient, err := utils.NewSSHClient(&sshConfig)
	if err != nil {
		return fmt.Errorf("ssh连接失败:%v", err)
	}
	defer func() {
		if err := sshClient.Close(); err != nil {
			global.OPS_LOG.Error("ssh关闭失败", zap.Error(err))
		}
	}()

	err = global.OPS_DB.Where("status = ? and server_type != ?", 1, 3).Preload("SysProject").Preload("Platform").Find(&hostList).Error
	if err != nil {
		global.OPS_LOG.Error("获取主机信息失败:", zap.Error(err))
		return errors.New("获取主机信息失败")
	}

	for _, host := range hostList {
		var config request.PrometheusConfig
		var targets []string

		labels := make(map[string]string)
		// todo 这里写死剑气劫项目(后续查看如何优化到相应的配置中，避免写死)
		if host.SysProject.ProjectName == "剑气劫" {
			targets = append(targets, fmt.Sprintf("%s:%s", host.PrivateIp, global.OPS_CONFIG.Prometheus.NodeExporterPort))
		} else {
			targets = append(targets, fmt.Sprintf("%s:%s", host.PubIp, global.OPS_CONFIG.Prometheus.NodeExporterPort))
		}

		labels["platform"] = host.Platform.PlatformCode
		labels["job"] = host.SysProject.ProjectName
		labels["hostname"] = host.ServerName
		labels["type"] = "host"

		config.Targets = targets
		config.Labels = labels
		prometheusHostConfig = append(prometheusHostConfig, config)
	}

	jsonData, err := json.MarshalIndent(prometheusHostConfig, "", "  ")
	if err != nil {
		global.OPS_LOG.Error("json序列化失败:", zap.Error(err))
		return errors.New("json序列化失败")
	}

	configFilePath, err := utils.CreateFile(
		filepath.Join(global.OPS_CONFIG.Local.Path, "prometheus", time.Now().Format("2006-01-02")),
		"host.json",
		string(jsonData),
	)
	if err != nil {
		global.OPS_LOG.Error("生成文件失败:", zap.Error(err))
		return errors.New("生成文件失败")
	}

	err = utils.UploadFile(sshClient, configFilePath, fmt.Sprintf("%s/host.json", global.OPS_CONFIG.Prometheus.HostServerJsonDir))
	if err != nil {
		return fmt.Errorf("上传主机监控文件失败:%v", err)
	}

	return
}
