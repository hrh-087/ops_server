package gm

import "encoding/json"

func (h HttpClient) KickGameServer(serverId int) (response *HttpResponse, err error) {
	data := map[string]interface{}{
		"serverId": serverId,
		"blocked":  true,
	}

	params, _ := json.Marshal(data)

	response, err = h.Post("/server/kickGame", params)
	if err != nil {
		return
	}

	return
}

func (h HttpClient) KickFightServer(serverId int) (response *HttpResponse, err error) {
	data := map[string]interface{}{
		"serverId": serverId,
		"blocked":  true,
	}

	params, _ := json.Marshal(data)

	response, err = h.Post("/server/kickFight", params)
	if err != nil {
		return
	}

	return
}
