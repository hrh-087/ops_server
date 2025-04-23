package gm

import (
	"encoding/json"
	"strconv"
)

func (h HttpClient) GetPlayerId(account string) (response *HttpResponse, err error) {
	data := map[string]interface{}{
		"account": account,
	}

	params, _ := json.Marshal(data)

	response, err = h.Post("/account/getPlayerIdByAccountName", params)

	return
}

func (h HttpClient) ExportPlayerData(playerId string) (response *HttpResponse, err error) {
	id, err := strconv.ParseInt(playerId, 10, 64)

	if err != nil {
		return
	}

	data := map[string]interface{}{
		"playerId": id,
	}

	params, _ := json.Marshal(data)
	response, err = h.Post("/role/exportPlayerData", params)

	return
}

func (h HttpClient) ImportPlayerData(importId, outputId string, playerData interface{}) (response *HttpResponse, err error) {

	importPlayerId, err := strconv.ParseInt(importId, 10, 64)
	if err != nil {
		return
	}
	outputPlayerId, err := strconv.ParseInt(outputId, 10, 64)
	if err != nil {
		return
	}

	jpData, err := json.Marshal(&playerData)
	if err != nil {
		return
	}

	data := map[string]interface{}{
		"importId":       importPlayerId,
		"outputPlayerId": outputPlayerId,
		"data":           string(jpData),
	}

	params, _ := json.Marshal(data)
	response, err = h.Post("/role/importPlayerData", params)

	return

}

func (h HttpClient) UpdateLoginMapping(playerId string, account string) (response *HttpResponse, err error) {

	id, err := strconv.ParseInt(playerId, 10, 64)
	if err != nil {
		return
	}

	data := map[string]interface{}{
		"playerId": id,
		"account":  account,
	}

	params, _ := json.Marshal(data)

	response, err = h.Post("/account/updateLoginMapping", params)
	return
}
