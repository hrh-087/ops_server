package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/model/system"
	"ops-server/utils"
	"ops-server/utils/cloud"
	"ops-server/utils/cloud/request"
	"path/filepath"
	"time"
)

type InstallServerParams struct {
	TaskId       uuid.UUID
	GameServerId uint
}

func NewInstallServerTask(params InstallServerParams) (*asynq.TaskInfo, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	task := NewTask(InstallServerTypeName, payload)
	return global.AsynqClinet.Enqueue(task)
}

func HandleInstallServerTask(ctx context.Context, t *asynq.Task) (err error) {
	var gameServer system.SysGameServer
	//var task system.JobTask
	var resultList []string
	var params InstallServerParams
	var lbList []system.SysAssetsLb
	var successListenerList []system.SysAssetsListener

	err = json.Unmarshal(t.Payload(), &params)
	if err != nil {
		return fmt.Errorf("参数解析失败:%v", err)
	}

	if err = global.OPS_DB.Preload("Platform").Preload("GameType").Preload("Host").Where("id = ?", params.GameServerId).First(&gameServer).Error; err != nil {
		return fmt.Errorf("获取游戏服信息失败:%v", err)
	}
	if gameServer.Status == 1 || gameServer.Status == 2 {
		return errors.New("该服务已安装或正在安装")
	}

	// 获取负载均衡信息
	if err = global.OPS_DB.Preload("CloudProduce").Preload("Listener").Where("platform_id = ?", gameServer.PlatformId).Find(&lbList).Error; err != nil {
		return fmt.Errorf("获取负载均衡信息失败:%v", err)
	}

	// 修改游戏服状态
	gameServer.Status = 1
	global.OPS_DB.Save(&gameServer)

	// 捕捉err是否有值，有值时，记录日志，并返回错误
	defer func() {
		if err != nil {
			global.OPS_LOG.Error("安装游戏服失败", zap.Error(err))
			// 修改安装状态为失败
			gameServer.Status = 4
			global.OPS_DB.Save(&gameServer)

			// 删除已添加的监听器
			if len(successListenerList) > 0 {
				for _, listener := range successListenerList {

					var lb system.SysAssetsLb
					if err := global.OPS_DB.Preload("CloudProduce").Preload("Listener").Where("id = ?", listener.LbId).First(&lb).Error; err != nil {
						global.OPS_LOG.Error("获取负载均衡信息失败:%v", zap.Error(err))
						continue
					}
					deleteBackendMemberParams := request.Listener{
						AK:         lb.CloudProduce.SecretId,
						SK:         lb.CloudProduce.SecretKey,
						Region:     lb.CloudProduce.RegionId,
						ListenerId: listener.InstanceId,
					}
					if err := cloud.DeleteListener(deleteBackendMemberParams); err != nil {
						global.OPS_LOG.Error("删除监听器失败", zap.Error(err))
					}
				}
			}

			// 写入任务结果
			resultList = append(resultList, err.Error())
			WriteTaskResult(t, resultList)
		}
	}()

	sshConfig, err := GetSSHKey(gameServer.ProjectId, gameServer.Host.PubIp, gameServer.Host.SSHPort)
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

	// 定义游戏服目录
	gameServerDir := fmt.Sprintf("%s/%s/%s_%d",
		global.OPS_CONFIG.Game.GamePath,
		gameServer.Platform.PlatformCode,
		gameServer.GameType.Code,
		gameServer.Vmid,
	)
	// 本地保存目录
	localGameServerDir := fmt.Sprintf("%s/%s/%s_%d",
		filepath.Join(global.OPS_CONFIG.Local.Path, "gameConfig", time.Now().Format("2006-01-02")),
		gameServer.Platform.PlatformCode,
		gameServer.GameType.Code,
		gameServer.Vmid,
	)

	command := fmt.Sprintf("mkdir -p %s/data/hotswap", gameServerDir)
	resultList = append(resultList, command)
	_, err = utils.ExecuteSSHCommand(sshClient, command)
	if err != nil {
		return fmt.Errorf("创建游戏服目录失败:%v", err)
	}

	configFilePath, err := utils.CreateFile(localGameServerDir, "application.yml", gameServer.ConfigFile)
	if err != nil {
		return fmt.Errorf("创建配置文件失败:%v", err)
	}

	composeFilePath, err := utils.CreateFile(localGameServerDir, "docker-compose.yml", gameServer.ComposeFile)
	if err != nil {
		return fmt.Errorf("创建docker-compose文件失败:%v", err)
	}

	err = utils.UploadFile(sshClient, configFilePath, fmt.Sprintf("%s/data/application.yml", gameServerDir))
	if err != nil {
		return fmt.Errorf("上传配置文件失败:%v", err)
	}

	err = utils.UploadFile(sshClient, composeFilePath, fmt.Sprintf("%s/docker-compose.yml", gameServerDir))
	if err != nil {
		return fmt.Errorf("上传docker-compose文件失败:%v", err)
	}

	// 只有战斗服类型跟游戏服类型才需要创建监听器
	if gameServer.GameType.IsFight || gameServer.GameType.Code == "game" {
		// 添加监听器
		err = global.OPS_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var listenerPort int
			for _, lb := range lbList {
				var listener system.SysAssetsListener
				var listenerId string

				if len(lb.Listener) <= 0 {
					listenerPort = 10000
				} else {
					currentPort := lb.Listener[0].Port
					for _, num := range lb.Listener {
						if num.Port > currentPort {
							currentPort = num.Port
						}
					}
					listenerPort = currentPort + 1
				}
				listenerName := fmt.Sprintf("%s_%s_%d", gameServer.Platform.PlatformCode, gameServer.GameType.Code, gameServer.Vmid)
				// 定义请求参数
				lbRequestParams := request.Listener{
					AK:              lb.CloudProduce.SecretId,
					SK:              lb.CloudProduce.SecretKey,
					Region:          lb.CloudProduce.RegionId,
					LbId:            lb.InstanceId,
					ListenerName:    listenerName,
					ListenerPort:    int32(listenerPort),
					Protocol:        "TCP",
					BackendPollName: fmt.Sprintf("%s-%d", gameServer.Host.PrivateIp, gameServer.TcpPort),
					BackendAddr:     gameServer.Host.PrivateIp,
					BackendPort:     int32(gameServer.TcpPort),
					SubnetCidrId:    lb.SubnetCidrId, // 使用负载均衡的子网ID
				}

				listenerId, err = cloud.CreateListener(lbRequestParams)

				if err != nil {
					return fmt.Errorf("%v, lbId: %s", err, lb.InstanceId)
				}
				listener.ProjectId = lb.ProjectId
				listener.InstanceId = listenerId
				listener.Name = listenerName
				listener.Port = listenerPort
				listener.LbId = lb.ID
				listener.Protocol = "TCP"
				listener.BackendIp = gameServer.Host.PrivateIp
				listener.BackendPort = int(gameServer.TcpPort)
				listener.HostId = gameServer.HostId

				successListenerList = append(successListenerList, listener)
			}

			if err = tx.CreateInBatches(&successListenerList, 100).Error; err != nil {
				return fmt.Errorf("创建监听器失败: %v", err)
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	resultList = append(resultList, "安装成功")
	WriteTaskResult(t, resultList)

	gameServer.Status = 2
	global.OPS_DB.Save(&gameServer)

	return err
}
