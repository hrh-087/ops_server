package system

import (
	"errors"
	"github.com/gin-gonic/gin"
	"ops-server/model/common/request"
	gmRes "ops-server/model/common/response"
	"ops-server/utils/gm"
	"strconv"
)

func (g GmService) SetGameCron(ctx *gin.Context, cronParams request.GameServerActivityParams) (err error) {
	httpClient, err := gm.NewHttpClient(ctx, strconv.Itoa(cronParams.ServerId))
	if err != nil {
		return
	}

	_, err = httpClient.SetGameCron(cronParams.ServerId, cronParams.Key, cronParams.Cron)
	if err != nil {
		return errors.New("设置cron失败")
	}

	return
}

func (g GmService) GetGameCronList(ctx *gin.Context, serverId int) (data interface{}, err error) {

	var cronList []gmRes.GameServerCron
	httpClient, err := gm.NewHttpClient(ctx, strconv.Itoa(serverId))
	if err != nil {
		return
	}

	response, err := httpClient.GetGameCronList(serverId)
	if err != nil {
		return
	}

	for k, v := range response.Data.(map[string]interface{}) {
		cronList = append(cronList, gmRes.GameServerCron{
			ServerId: serverId,
			Key:      k,
			Cron:     v.(string),
		})
	}

	return cronList, nil
}

func (g GmService) GetActivityExtra(ctx *gin.Context, serverId int, key string) (data interface{}, err error) {
	httpClient, err := gm.NewHttpClient(ctx, strconv.Itoa(serverId))
	if err != nil {
		return
	}

	response, err := httpClient.GetActivityExtra(serverId, key)
	if err != nil {
		return
	}

	return response.Data, nil
}
