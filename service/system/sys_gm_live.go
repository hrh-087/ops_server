package system

import (
	"github.com/gin-gonic/gin"
	"ops-server/model/common/request"
	"ops-server/utils/gm"
	"strconv"
	"time"
)

func (g GmService) GetActivityServerList(ctx *gin.Context, serverId int) (data interface{}, err error) {

	httpClient, err := gm.NewHttpClient(ctx, strconv.Itoa(serverId))
	if err != nil {
		return
	}

	response, err := httpClient.GetActivityServerList(serverId)
	if err != nil {
		return
	}

	//data = make([]map[string]interface{}, 0)
	//
	//str, ok := response.Data.(string)
	//if !ok {
	//	return nil, errors.New("获取活动服务器列表失败")
	//}
	//
	//err = json.Unmarshal([]byte(str), &data)
	//
	//return data, nil
	return response.Data, nil
}

func (g GmService) KickLivePlayer(ctx *gin.Context, kickData request.LiveKickPlayer) (err error) {

	// todo 临时测试写死live的vmid为12001
	//kickData.Vmid = "12001"

	httpClient, err := gm.NewHttpClient(ctx, strconv.Itoa(kickData.ServerId))
	if err != nil {
		return
	}
	_, err = httpClient.KickLivePlayer(kickData.ServerId, kickData.GroupId, kickData.PlayerId)
	return err
}

func (g GmService) ClearActivityServer(ctx *gin.Context, serverId int) (err error) {
	httpClient, err := gm.NewHttpClient(ctx, strconv.Itoa(serverId))
	if err != nil {
		return
	}

	_, err = httpClient.ClearActivityServer(serverId)
	return
}

func (g GmService) SetLiveActivityServer(ctx *gin.Context, liveData request.LiveServerActivity) (err error) {
	if liveData.Id == 0 {
		liveData.Id = int(time.Now().Unix())
	}
	liveData.UpdateTime = time.Now().UnixMilli()

	httpClient, err := gm.NewHttpClient(ctx, strconv.Itoa(liveData.ServerId))
	if err != nil {
		return
	}

	_, err = httpClient.SetActivityServer(
		liveData.Id,
		liveData.ServerId,
		liveData.NameForTQT,
		liveData.Number,
		liveData.HomeOwner,
		liveData.Member,
		int(liveData.UpdateTime),
		liveData.StartTime,
		liveData.EndTime,
	)

	return
}
