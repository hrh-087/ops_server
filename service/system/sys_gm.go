package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"ops-server/model/common/request"
	gmRes "ops-server/model/common/response"
	"ops-server/utils/gm"
	"sort"
	"strconv"
	"strings"
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
		reward.Rank = rewardData["rank"].(string)
		reward.OpenId = int(rewardData["openId"].(float64))
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

	checkRankConfigMap := make(map[int][]request.GmRankConfig)

	// 给排行榜分组
	for _, rank := range rankConfig {
		if _, ok := checkRankConfigMap[rank.RankId]; !ok {
			checkRankConfigMap[rank.RankId] = []request.GmRankConfig{rank}
		} else {
			checkRankConfigMap[rank.RankId] = append(checkRankConfigMap[rank.RankId], rank)
		}
	}

	// 根据排行榜开始时间排序
	for _, rankInfo := range checkRankConfigMap {
		sort.Slice(rankInfo, func(i, j int) bool {
			firstTime, err := time.Parse("2006-01-02 15:04:05", rankInfo[i].StartTime)
			if err != nil {
				return false
			}

			secondTime, err := time.Parse("2006-01-02 15:04:05", rankInfo[j].StartTime)
			if err != nil {
				return false
			}
			return firstTime.Before(secondTime)
		})
	}

	// 检测排行榜配置时间是否有交叉
	for _, rankInfo := range checkRankConfigMap {

		for i := 1; i < len(rankInfo); i++ {
			// 检测时间是否有交叉
			startTime, err := time.Parse("2006-01-02 15:04:05", rankInfo[i].StartTime)
			if err != nil {
				return errors.New(fmt.Sprintf("排行榜唯一id(%d)配置时间无法解析", rankInfo[i].Id))
			}

			closeTime, err := time.Parse("2006-01-02 15:04:05", rankInfo[i-1].CloseTime)
			if err != nil {
				return errors.New(fmt.Sprintf("排行榜唯一id(%d)配置时间无法解析", rankInfo[i-1].Id))
			}

			if startTime.Before(closeTime) {
				return errors.New(fmt.Sprintf("排行榜id%d配置时间有交叉", rankInfo[i].RankId))
			}
		}
	}

	rewardId := 1

	for _, rank := range rankConfig {

		if rank.Id == 0 || rank.StartTime == "" || rank.EndTime == "" || rank.RankId == 0 || rank.ShowCount == 0 || rank.CloseTime == "" {
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

		closeTime, err := time.Parse("2006-01-02 15:04:05", rank.CloseTime)
		if err != nil {
			return errors.New("关闭时间配置有误,请确认格式为:(2006-01-02 15:04:05)")
		}

		if endTime.Before(startTime) {
			return errors.New("结束时间不能小于开始时间")
		} else if closeTime.Before(endTime) {
			return errors.New("关闭时间不能小于结束时间")
		}

		rankOpenConfig = append(rankOpenConfig, gmRes.RankOpenConfig{
			ID:        rank.Id,
			RankID:    rank.RankId,
			StartTime: rank.StartTime,
			EndTime:   rank.EndTime,
			CloseTime: rank.CloseTime,
			ShowCount: rank.ShowCount,
		})

		for _, reward := range rank.RewardList {

			rewardData := make(map[string]int)
			for _, v := range reward.Rewards {
				rewardData[v.RewardId] = v.RewardNum
			}

			rankInfo := strings.Split(reward.Rank, ",")

			// 获取排名
			if len(rankInfo) > 1 {
				startRewardRank, err := strconv.ParseInt(rankInfo[0], 10, 64)
				if err != nil {
					return errors.New("解析排行榜排名失败")
				}

				endRewardRank, err := strconv.ParseInt(rankInfo[1], 10, 64)
				if err != nil {
					return errors.New("解析排行榜排名失败")
				}

				for i := startRewardRank; i <= endRewardRank; i++ {

					rankRewardConfig = append(rankRewardConfig, gmRes.RankRewardConfig{
						ID:      rewardId,
						OpenId:  reward.OpenId,
						Rank:    int(i),
						Rewards: rewardData,
					})
					rewardId++
				}
			} else {

				rewardRank, err := strconv.ParseInt(rankInfo[0], 10, 64)
				if err != nil {
					return errors.New("解析排行榜排名失败")
				}
				rankRewardConfig = append(rankRewardConfig, gmRes.RankRewardConfig{
					ID:      rewardId,
					OpenId:  reward.OpenId,
					Rank:    int(rewardRank),
					Rewards: rewardData,
				})
			}
			// 累加奖励唯一id
			rewardId++
		}
	}

	//fmt.Printf("rankOpenConfig: %+v, rankRewardConfig: %+v\n", rankOpenConfig, rankRewardConfig)

	_, err = httpClient.SetRankConfig(serverId, rankOpenConfig, rankRewardConfig)
	return err
}

// 上传游戏服策划配置
func (g GmService) UploadGameConfig(ctx *gin.Context, data map[string]interface{}) (err error) {

	for k, v := range data {
		fmt.Println(v)
		switch k {
		case "rank":
			var rankList []request.Rank

			err = json.Unmarshal([]byte(v.(string)), &rankList)
			if err != nil {
				return errors.New("解析rank表失败")
			}

			for _, rank := range rankList {
				fmt.Printf("rank: %+v\n", rank)
			}

		case "item":

			var itemList []request.Item
			err = json.Unmarshal([]byte(v.(string)), &itemList)
			if err != nil {
				return errors.New("解析item表失败")
			}

			for _, item := range itemList {
				fmt.Printf("item: %+v\n", item)
			}
		}
	}
	return err
}
