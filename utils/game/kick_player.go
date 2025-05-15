package game

import (
	"context"
	"errors"
	"ops-server/utils/gm"
	"strconv"
)

func KickPlayer(ctx context.Context, serverId int) (err error) {
	httpClient, err := gm.NewHttpClient(ctx, strconv.Itoa(serverId))
	if err != nil {
		return
	}
	// 踢战斗服玩家
	fightRes, err := httpClient.KickFightServer(serverId)
	if err != nil {
		return
	} else if fightRes.Code != 0 {
		return errors.New(fightRes.Msg)
	}
	// 踢游戏服玩家
	gameRes, err := httpClient.KickGameServer(serverId)
	if err != nil {
		return
	} else if gameRes.Code != 0 {
		return errors.New(gameRes.Msg)
	}
	return
}
