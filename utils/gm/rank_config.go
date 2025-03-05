package gm

import (
	"encoding/json"
	gmRes "ops-server/model/common/response"
)

func (h HttpClient) GetRankList(serverId int) (response *HttpResponse, err error) {

	data := map[string]interface{}{
		"serverId": serverId,
	}

	params, _ := json.Marshal(data)

	response, err = h.Post("/config/getRankOpenConfig", params)

	if err != nil {
		return
	}

	return
}

func (h HttpClient) GetRewardConfig(serverId, rankId int) (response *HttpResponse, err error) {
	data := map[string]interface{}{
		"serverId": serverId,
		"rankId":   rankId,
	}

	params, _ := json.Marshal(data)

	response, err = h.Post("/config/getRankRewardConfig", params)

	if err != nil {
		return
	}

	return
}

func (h HttpClient) SetRankConfig(serverId int, openConfig []gmRes.RankOpenConfig, rewardConfig []gmRes.RankRewardConfig) (response *HttpResponse, err error) {

	data := map[string]interface{}{
		"serverId":         serverId,
		"rankOpenConfig":   openConfig,
		"rankRewardConfig": rewardConfig,
	}

	params, _ := json.Marshal(data)

	response, err = h.Post("/config/setRankConfig", params)

	if err != nil {
		return
	}

	return
}
