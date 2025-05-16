package gm

import (
	"encoding/json"
)

// 设置游戏服开关
func (h HttpClient) SetSwitch(serverId int, typeKey string, state bool) (response *HttpResponse, err error) {

	data := map[string]interface{}{
		"serverId": serverId,
		"key":      typeKey,
		"state":    state,
	}

	params, _ := json.Marshal(data)

	return h.Post("/switch/setSwitch", params)
}

func (h HttpClient) GetSwitchList(serverId int) (response *HttpResponse, err error) {
	data := map[string]interface{}{
		"serverId": serverId,
	}

	params, _ := json.Marshal(data)

	return h.Post("/switch/getSwitchList", params)

}
