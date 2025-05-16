package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"ops-server/global"
	"ops-server/model/system"
	"ops-server/utils/game"
	"ops-server/utils/notice"
	"strconv"
	"strings"
)

type KickPlayer struct {
	ProjectId uint `json:"projectId"`
}

func HandleKickPlayer(ctx context.Context, t *asynq.Task) (err error) {
	//log.Println("cron kick player1")
	//global.OPS_LOG.Info("cron kick player")
	//return nil

	var params KickPlayer
	var resultList []string
	var platformList []system.SysGamePlatform

	defer func() {
		if err != nil {
			resultList = append(resultList, err.Error())
			WriteTaskResult(t, resultList)
		}
	}()

	if err = json.Unmarshal(t.Payload(), &params); err != nil {
		return errors.New("参数解析失败")
	}

	// 将projectId写入context
	ctx = context.WithValue(ctx, "projectId", strconv.Itoa(int(params.ProjectId)))

	projectId := ctx.Value("projectId").(string)
	fmt.Println("projectId:", projectId)

	err = global.OPS_DB.Debug().WithContext(ctx).Where("project_id = ?", params.ProjectId).Find(&platformList).Error
	//err = global.OPS_DB.Debug().WithContext(ctx).Where("project_id = ? and platform_code = 887706", params.ProjectId).Find(&platformList).Error
	if err != nil {
		return errors.New("获取项目游戏平台失败")
	}

	var message string

	var messageList []string

	for _, platform := range platformList {
		serverId, err := strconv.ParseInt(platform.PlatformCode, 10, 64)
		if err != nil {
			return err
		}
		err = game.KickPlayer(ctx, int(serverId))
		if err != nil {
			message = fmt.Sprintf("平台ID:%d,踢人失败", serverId)
			return errors.New(message)
		}

		message = fmt.Sprintf("平台ID:%d,踢人成功", serverId)
		resultList = append(resultList, message)
		messageList = append(messageList, message)
	}

	// 通知dingding
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": strings.Join(messageList, "\n"),
		},
	}

	webhook := "https://oapi.dingtalk.com/robot/send?access_token=6693350e344a346c332bcfb6a059d9edd2e4c70461e525dcf84e8532418bde49"
	secret := "SEC4e35b9b0e11bbe423dcd930d0a0e1993c67cb766d7435cbc66e8851f594c59c0"

	err = notice.SendDingTalkMessage(webhook, secret, msg)
	if err != nil {
		return errors.New("通知钉钉失败")
	}

	WriteTaskResult(t, resultList)
	return nil
}
