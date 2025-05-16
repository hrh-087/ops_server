package system

import (
	"errors"
	"github.com/gin-gonic/gin"
	"ops-server/utils/gm"
)

func (g GmService) ImportPlayerData(ctx *gin.Context, account, importId, outputId, serverId string) (err error) {
	// 解析playerId

	onlineHttpClient, err := gm.NewHttpClient(ctx, "online")
	if err != nil {
		return
	}

	localHttpClient, err := gm.NewHttpClient(ctx, serverId)
	if err != nil {
		return
	}
	// 导出玩家数据
	sourcePlayerData, err := onlineHttpClient.ExportPlayerData(outputId)
	if err != nil {
		return
	}
	// 导入玩家数据
	_, err = localHttpClient.ImportPlayerData(importId, outputId, sourcePlayerData.Data)
	if err != nil {
		return
	}
	// 更改登录映射
	_, err = localHttpClient.UpdateLoginMapping(outputId, account)
	if err != nil {
		return
	}

	return
}

func (g GmService) GetPlayerIdByAccountName(ctx *gin.Context, account, serverId string) (result interface{}, err error) {
	httpClient, err := gm.NewHttpClient(ctx, serverId)
	if err != nil {
		return
	}

	response, err := httpClient.GetPlayerId(account)
	if err != nil {
		return
	}

	if response.Code != 0 {
		return nil, errors.New(response.Msg)
	}

	return response.Data, nil

}
