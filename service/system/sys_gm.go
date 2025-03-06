package system

import (
	"errors"
	"github.com/gin-gonic/gin"
	"ops-server/model/common/request"
	gmRes "ops-server/model/common/response"
	"ops-server/utils/gm"
	"time"
)

type GmService struct {
}

// SetSwitch
// 设置gm开关
func (g GmService) SetSwitch(ctx *gin.Context, serverId int, typeKey string, state bool) (msg string, err error) {

	httpClient, err := gm.NewHttpClient(ctx)
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

// GetSwitchList
// 获取gm开关列表
func (g GmService) GetSwitchList(ctx *gin.Context, serverId int) (data interface{}, err error) {

	var switchList []gmRes.GmSwitch

	var switchMap = map[string]string{
		"MATCH_LIMIT": "匹配限制开关",
		"MAC_LIMIT":   "限制mac地址登录开关",
	}
	//projectId := ctx.GetString("projectId")
	//
	//if projectId == "" {
	//	return
	//}

	httpClient, err := gm.NewHttpClient(ctx)
	if err != nil {
		return
	}

	response, err := httpClient.GetSwitchList(serverId)
	if err != nil {
		return
	}

	for k, v := range response.Data.(map[string]interface{}) {
		switchList = append(switchList, gmRes.GmSwitch{
			ServerId: serverId,
			TypeKey:  k,
			TypeName: switchMap[k],
			State:    v.(bool),
		})
	}
	return switchList, nil
}

// GetRankList
// 获取排行榜列表
func (g GmService) GetRankList(ctx *gin.Context, serverId int) (data interface{}, err error) {

	httpClient, err := gm.NewHttpClient(ctx)
	if err != nil {
		return
	}

	response, err := httpClient.GetRankList(serverId)
	if err != nil {
		return
	}

	return response.Data, nil
}

// GetRankRewardList
// 获取排行榜奖励配置
func (g GmService) GetRankRewardList(ctx *gin.Context, serverId, id int) (data interface{}, err error) {
	var rewardList []request.GmRankReward
	var rewards []request.Reward
	httpClient, err := gm.NewHttpClient(ctx)
	if err != nil {
		return
	}

	response, err := httpClient.GetRewardConfig(serverId, id)
	if err != nil {
		return
	}

	for _, v := range response.Data.([]interface{}) {
		var reward request.GmRankReward
		rewardData, ok := v.(map[string]interface{})
		if !ok {
			return nil, errors.New("获取排行榜奖励配置失败")
		}

		for rewardId, rewardNum := range rewardData["rewards"].(map[string]interface{}) {
			rewards = append(rewards, request.Reward{
				RewardId:  rewardId,
				RewardNum: int(rewardNum.(float64)),
			})
		}

		reward.Id = int(rewardData["id"].(float64))
		reward.Rank = int(rewardData["rank"].(float64))
		reward.RankId = int(rewardData["rankId"].(float64))
		reward.Rewards = rewards

		rewardList = append(rewardList, reward)

	}

	return rewardList, nil
}

// SetRankConfig
// 设置排行榜配置
func (g GmService) SetRankConfig(ctx *gin.Context, serverId int, rankConfig []request.GmRankConfig) (err error) {
	var rankOpenConfig []gmRes.RankOpenConfig
	var rankRewardConfig []gmRes.RankRewardConfig

	httpClient, err := gm.NewHttpClient(ctx)
	if err != nil {
		return
	}

	for _, rank := range rankConfig {
		if rank.Id == 0 || rank.StartTime == "" || rank.EndTime == "" || rank.RankId == 0 || rank.ShowCount == 0 {
			return errors.New("排行榜配置有误")
		}

		startTime, err := time.Parse("2006-01-02 15:04:05", rank.StartTime)
		if err != nil {
			return errors.New("开始时间配置有误,请确认格式为:(2006-01-02 15:04:05)")
		}

		endTime, err := time.Parse("2006-01-02 15:04:05", rank.EndTime)
		if err != nil {
			return errors.New("结束时间配置有误,请确认格式为:(2006-01-02 15:04:05)")
		}

		if endTime.Before(startTime) {
			return errors.New("结束时间不能小于开始时间")
		}

		rankOpenConfig = append(rankOpenConfig, gmRes.RankOpenConfig{
			ID:        rank.Id,
			RankID:    rank.RankId,
			StartTime: rank.StartTime,
			EndTime:   rank.EndTime,
			ShowCount: rank.ShowCount,
		})

		for _, reward := range rank.RewardList {

			rewardData := make(map[string]int)
			for _, v := range reward.Rewards {
				rewardData[v.RewardId] = v.RewardNum
			}

			rankRewardConfig = append(rankRewardConfig, gmRes.RankRewardConfig{
				ID:      reward.Id,
				RankID:  reward.RankId,
				Rank:    reward.Rank,
				Rewards: rewardData,
			})
		}
	}

	//fmt.Printf("rankOpenConfig: %+v, rankRewardConfig: %+v\n", rankOpenConfig, rankRewardConfig)

	_, err = httpClient.SetRankConfig(serverId, rankOpenConfig, rankRewardConfig)
	return err
}
