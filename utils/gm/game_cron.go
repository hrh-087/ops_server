package gm

import "encoding/json"

func (h HttpClient) SetGameCron(serverId int, key, cron string) (response *HttpResponse, err error) {
	data := map[string]interface{}{
		"serverId": serverId,
		"key":      key,
		"cron":     cron,
	}

	params, _ := json.Marshal(data)

	return h.Post("/cron/cronReset", params)

}

func (h HttpClient) GetGameCronList(serverId int) (response *HttpResponse, err error) {
	data := map[string]interface{}{
		"serverId": serverId,
	}

	params, _ := json.Marshal(data)

	return h.Post("/cron/cronList", params)
}

func (h HttpClient) SetActivityExtra(serverId int, key string, extraData map[string]interface{}) (response *HttpResponse, err error) {
	data := map[string]interface{}{
		"serverId":  serverId,
		"key":       key,
		"extraData": extraData,
	}

	params, _ := json.Marshal(data)

	return h.Post("/activity/activityExtra", params)

}

func (h HttpClient) GetActivityExtra(serverId int, key string) (response *HttpResponse, err error) {
	data := map[string]interface{}{
		"serverId": serverId,
		"key":      key,
	}

	params, _ := json.Marshal(data)

	return h.Post("/activity/getActivityExtra", params)

}
