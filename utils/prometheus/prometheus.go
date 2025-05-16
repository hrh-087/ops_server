package prometheus

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/job/task"
	"ops-server/model/common/request"
	"ops-server/model/system"
	"ops-server/utils"
	"path/filepath"
	"strconv"
	"time"
)

func GenerateGameMonitorFile(projectId uint, gameServerList []system.SysGameServer, isMaintenance bool) (err error) {

	var prometheusGameConfig []request.PrometheusConfig

	sshConfig, err := task.GetSSHKey(projectId, global.OPS_CONFIG.Prometheus.Addr, global.OPS_CONFIG.Prometheus.SshPort)
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

	for _, gameServer := range gameServerList {
		var config request.PrometheusConfig
		var targets []string

		labels := make(map[string]string)
		if gameServer.SysProject.ProjectName == "剑气劫" {
			targets = append(targets, fmt.Sprintf("%s:%d", gameServer.Host.PrivateIp, gameServer.HttpPort))
		} else {
			targets = append(targets, fmt.Sprintf("%s:%d", gameServer.Host.PubIp, gameServer.HttpPort))
		}

		labels["platform"] = gameServer.Platform.PlatformCode
		labels["job"] = gameServer.SysProject.ProjectName
		labels["hostname"] = gameServer.Host.ServerName
		labels["gamename"] = fmt.Sprintf("%s_%d", gameServer.GameType.Code, gameServer.Vmid)
		labels["type"] = "game"
		labels["isMaintenance"] = strconv.FormatBool(isMaintenance)
		labels["isFight"] = strconv.FormatBool(gameServer.GameType.IsFight)

		config.Targets = targets
		config.Labels = labels

		prometheusGameConfig = append(prometheusGameConfig, config)
	}

	jsonData, err := json.MarshalIndent(prometheusGameConfig, "", "  ")
	if err != nil {
		global.OPS_LOG.Error("json序列化失败:", zap.Error(err))
		return errors.New("json序列化失败")
	}

	configFilePath, err := utils.CreateFile(
		filepath.Join(global.OPS_CONFIG.Local.Path, "prometheus", time.Now().Format("2006-01-02")),
		"game_server.json",
		string(jsonData),
	)
	if err != nil {
		global.OPS_LOG.Error("生成文件失败:", zap.Error(err))
		return errors.New("生成文件失败")
	}

	err = utils.UploadFile(sshClient, configFilePath, fmt.Sprintf("%s/game_server.json", global.OPS_CONFIG.Prometheus.GameServerJsonDir))
	if err != nil {
		return fmt.Errorf("上传游戏服监控文件失败:%v", err)
	}

	return
}
