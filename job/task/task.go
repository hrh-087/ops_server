package task

import (
	"errors"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"ops-server/global"
	"ops-server/model/system"
	"ops-server/utils"
	"strings"
	"time"
)

const (
	BatchCommandTypeName = "game:batchCommand" // 批量命令

	// 游戏服
	UpdateGameConfigTypeName    = "game:updateGameConfig"    // 更新游戏服配置
	RsyncGameConfigTypeName     = "game:rsyncGameConfig"     // 同步游戏服配置
	RsyncGameJsonConfigTypeName = "game:rsyncGameJsonConfig" // 同步游戏服json配置文件
	RsyncGameScriptTypeName     = "game:rsyncGameScript"     // 同步游戏服脚本
	InstallGameServerTypeName   = "game:installGameServer"   // 安装游戏服
	StartGameTypeName           = "game:startGame"           // 启动游戏服
	StopGameTypeName            = "game:stopGame"            // 关闭游戏服
	UpdateGameImageTypeName     = "game:updateGameImage"     // 更新游戏服镜像
	CheckGameVersionTypeName    = "game:checkGameVersion"    // 检查游戏服版本

	// 热更游戏服代码
	HotGameUnzipFileTypeName   = "game:HotGameUnzipFile"   // 解压热更文件包
	HotGameRsyncHostTypeName   = "game:HotGameRsyncHost"   // 同步到服务器
	HotGameRsyncServerTypeName = "game:HotGameRsyncServer" // 同步到游戏服

)

// 任务存储时间
func retention() asynq.Option {
	return asynq.Retention(time.Duration(global.OPS_CONFIG.Asynq.Retention) * time.Hour * 24)
}

func retryCount() asynq.Option {
	return asynq.MaxRetry(global.OPS_CONFIG.Asynq.MaxRetryCount)
}

func taskTimeout() asynq.Option {
	return asynq.Timeout(time.Duration(global.OPS_CONFIG.Asynq.Timeout) * time.Second)
}

func NewTask(serverType string, payload []byte, opts ...asynq.Option) *asynq.Task {
	opts = append(opts, retention(), retryCount(), taskTimeout())
	return asynq.NewTask(serverType, payload, opts...)
}

func WriteTaskResult(t *asynq.Task, result []string) {

	resultStr := strings.Join(result, "\r\n")

	_, err := t.ResultWriter().Write([]byte(resultStr))
	if err != nil {
		global.OPS_LOG.Error("写入任务结果失败", zap.Error(err))
	}
}

func GetSSHKey(projectId uint, host, port string) (auth utils.SShConfig, err error) {
	var sshKey system.SysSshAuth
	err = global.OPS_DB.First(&sshKey, "project_id = ?", projectId).Error
	if err != nil {
		return auth, err
	}
	if sshKey.UsePass {
		if strings.TrimSpace(sshKey.Password) == "" {
			return auth, errors.New("密码为空")
		}
		return utils.SShConfig{
			User:     sshKey.User,
			Password: sshKey.Password,
			Host:     host,
			Port:     port,
		}, err
	} else {
		if strings.TrimSpace(sshKey.PrivateKey) == "" {
			return utils.SShConfig{}, errors.New("私钥为空")
		}
		return utils.SShConfig{
			User:                 sshKey.User,
			Host:                 host,
			Port:                 port,
			PrivateKey:           sshKey.PrivateKey,
			PrivateKeyPassphrase: sshKey.PrivateKeyPassphrase,
		}, err
	}
}

func GetSSHConn(projectId uint, pubIp string, port string) (client *ssh.Client, err error) {
	sshConfig, err := GetSSHKey(projectId, pubIp, port)
	if err != nil {
		global.OPS_LOG.Error("获取ssh配置失败", zap.Error(err))
		return nil, err
	}

	return utils.NewSSHClient(&sshConfig)
}
