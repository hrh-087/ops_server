package gm

import (
	"encoding/json"
)

// 直播服管理

// KickLivePlayer 直播服踢人
func (h HttpClient) KickLivePlayer(serverId int, groupId int, playerId int) (response *HttpResponse, err error) {

	data := map[string]interface{}{
		"serverId": serverId,
		//"vmid":     vmid,
	}

	if groupId == 0 && playerId != 0 {
		data["playerId"] = playerId
	} else {
		data["groupId"] = groupId
	}

	params, _ := json.Marshal(data)
	return h.Post("/server/kickLivePlayer", params)
}

// SetActivityServer 直播服白名单设置
func (h HttpClient) SetActivityServer(id, serverId int, nameForTQT string, number int, homeowner, member []int, updateTime int, startTime, endTime string) (response *HttpResponse, err error) {
	data := map[string]interface{}{
		"id":         id,
		"serverId":   serverId,
		"nameForTQT": nameForTQT,
		"number":     number,
		"homeowner":  homeowner,
		"member":     member,
		"updateTime": updateTime,
		"startTime":  startTime,
		"endTime":    endTime,
	}

	params, _ := json.Marshal(data)

	//return h.Post("/activityServer/setActivityServer", params)
	return h.Post("/activityServer/setActivityServerTask", params)
}

// GetActivityServerList 获取直播服列表
func (h HttpClient) GetActivityServerList(serverId int) (response *HttpResponse, err error) {
	data := map[string]interface{}{
		"serverId": serverId,
	}

	params, _ := json.Marshal(data)

	return h.Post("/activityServer/getActivityServerTaskParams", params)
}

// ClearActivityServer 清除直播服白名单
func (h HttpClient) ClearActivityServer(serverId int) (response *HttpResponse, err error) {
	data := map[string]interface{}{
		"serverId": serverId,
	}

	params, _ := json.Marshal(data)

	return h.Post("/activityServer/clearActivityServerData", params)
}
