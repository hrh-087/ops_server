package system

import (
	"github.com/gin-gonic/gin"
	"ops-server/utils/gm"
)

type GmService struct {
}

type GmSwitch struct {
	ServerId int    `json:"serverId"`
	TypeKey  string `json:"typeKey"`
	TypeName string `json:"typeName"`
	State    bool   `json:"state"`
}

func (g GmService) SetSwitch(ctx *gin.Context, serverId int, typeKey string, state bool) (msg string, err error) {

	projectId := ctx.GetString("projectId")

	if projectId == "" {
		return
	}

	httpClient, err := gm.NewHttpClient(projectId)
	//httpClient, err := gm.NewHttpClient("http://10.3.21.48:10060")
	if err != nil {
		return
	}

	response, err := httpClient.SetSwitch(serverId, typeKey, state)
	//response, err := httpClient.SetSwitch(1, "MATCH_LIMIT", true)
	if err != nil {
		return
	}

	return response.Msg, nil
}

func (g GmService) GetSwitchList(ctx *gin.Context, serverId int) (data interface{}, err error) {

	var switchList []GmSwitch

	var switchMap = map[string]string{
		"MATCH_LIMIT": "匹配限制开关",
		"MAC_LIMIT":   "限制mac地址登录开关",
	}
	projectId := ctx.GetString("projectId")

	if projectId == "" {
		return
	}

	httpClient, err := gm.NewHttpClient(projectId)
	if err != nil {
		return
	}

	response, err := httpClient.GetSwitchList(serverId)
	if err != nil {
		return
	}

	for k, v := range response.Data.(map[string]interface{}) {
		switchList = append(switchList, GmSwitch{
			ServerId: serverId,
			TypeKey:  k,
			TypeName: switchMap[k],
			State:    v.(bool),
		})
	}
	return switchList, nil
}
