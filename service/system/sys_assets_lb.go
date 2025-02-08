package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/system"
	"ops-server/utils"
	"ops-server/utils/cloud"
	"ops-server/utils/cloud/hw_cloud"
	"ops-server/utils/cloud/hw_cloud/elb"
	cloudRequest "ops-server/utils/cloud/request"
	"strconv"
)

type AssetsLbService struct{}

func (s *AssetsLbService) GetAssetsLbList(ctx *gin.Context, info request.PageInfo) (result interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Preload("CloudProduce").Preload("Platform").Model(&system.SysAssetsLb{})

	var resultList []system.SysAssetsLb

	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Set("skip_project_filter", true).Error
	if err != nil {
		return resultList, total, err
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	err = db.Preload("CloudProduce").Order(OrderStr).Find(&resultList).Error
	return resultList, total, err
}

func (s *AssetsLbService) GetAssetsLbById(ctx *gin.Context, id int) (result system.SysAssetsLb, err error) {
	err = global.OPS_DB.Preload("Listener").WithContext(ctx).Preload("Platform").Preload("CloudProduce").Where("id = ?", id).First(&result).Error
	return
}

func (s *AssetsLbService) CreateAssetsLb(ctx *gin.Context, lb system.SysAssetsLb) (err error) {
	err = global.OPS_DB.WithContext(ctx).Where("instance_id = ?", lb.InstanceId).FirstOrCreate(&lb).Error
	return
}

func (s *AssetsLbService) DeleteAssetsLb(ctx *gin.Context, id int) (err error) {
	return global.OPS_DB.WithContext(ctx).Unscoped().Delete(&system.SysAssetsLb{}, id).Error
}

func (s *AssetsLbService) RsyncAssetsCloudLb(ctx *gin.Context, assetsLb system.SysAssetsLb) (err error) {
	var hwElb elb.HwElb
	var cloudProduce system.SysCloudProduce
	var platform system.SysGamePlatform

	if err = global.OPS_DB.WithContext(ctx).Where("id = ?", assetsLb.CloudProduceId).First(&cloudProduce).Error; err != nil {
		global.OPS_LOG.Error("获取云平台信息失败:", zap.Error(err))
		return errors.New("获取云平台信息失败")
	}
	if err = global.OPS_DB.WithContext(ctx).Where("id = ?", assetsLb.PlatformId).First(&platform).Error; err != nil {
		global.OPS_LOG.Error("获取游戏平台信息失败:", zap.Error(err))
		return errors.New("获取游戏平台信息失败")
	}

	client := elb.NewHwElb(hw_cloud.HWCloud{
		AK:     cloudProduce.SecretId,
		SK:     cloudProduce.SecretKey,
		Region: cloudProduce.RegionId,
	})

	if client == nil {
		global.OPS_LOG.Error("初始化云商信息失败")
		return errors.New("初始化云商信息失败")
	}

	cloudLbList, err := hwElb.GetLbList(client, "", 0)
	if err != nil {
		global.OPS_LOG.Error("获取云负载均衡列表失败:", zap.Error(err))
		return
	}

	//var lbList []system.SysAssetsLb
	for _, lbItem := range *cloudLbList.Loadbalancers {
		var lb system.SysAssetsLb
		lb.Name = lbItem.Name
		lb.InstanceId = lbItem.Id
		lb.PubIp = *lbItem.Eips[0].EipAddress
		lb.PrivateIp = lbItem.VipAddress
		lb.CloudProduceId = assetsLb.CloudProduceId
		lb.PlatformId = assetsLb.PlatformId
		lb.SubnetCidrId = lbItem.VipSubnetCidrId

		//lbList = append(lbList, lb)

		err = s.CreateAssetsLb(ctx, lb)
		if err != nil {
			global.OPS_LOG.Error("创建负载均衡记录失败:", zap.Error(err))
			return
		}
	}
	return
}

func (s *AssetsLbService) RsyncLbListener(ctx *gin.Context) (err error) {

	var platformList []system.SysGamePlatform
	if err = global.OPS_DB.WithContext(ctx).Find(&platformList).Error; err != nil {
		global.OPS_LOG.Error("获取游戏平台信息失败:", zap.Error(err))
		return errors.New("获取游戏平台信息失败")
	}

	for _, platform := range platformList {
		var lbList []system.SysAssetsLb
		var gameServerList []system.SysGameServer

		// 获取游戏服信息
		if err = global.OPS_DB.WithContext(ctx).Preload("GameType").Preload("Platform").Preload("Host").Where("status = 2 and platform_id = ?", platform.ID).Find(&gameServerList).Error; err != nil {
			global.OPS_LOG.Error("获取游戏服信息失败:", zap.Error(err))
			return errors.New("获取游戏服信息失败")
		}

		// 获取负载均衡信息
		if err = global.OPS_DB.WithContext(ctx).Preload("CloudProduce").Preload("Listener").Where("platform_id = ?", platform.ID).Find(&lbList).Error; err != nil {
			global.OPS_LOG.Error("获取负载均衡信息失败:", zap.Error(err))
			return errors.New("获取负载均衡信息失败")
		}

		for _, gameServer := range gameServerList {
			// 跳过非战斗服跟游戏服
			if !gameServer.GameType.IsFight && gameServer.GameType.Code != "game" {
				fmt.Printf("跳过非战斗服跟游戏服: %s\n", gameServer.Name)
				continue
			}

			for _, lb := range lbList {
				listenerName := fmt.Sprintf("%s_%s_%d", gameServer.Platform.PlatformCode, gameServer.GameType.Code, gameServer.Vmid)
				err = global.OPS_DB.Debug().WithContext(ctx).Where("lb_id = ? and name = ?", lb.ID, listenerName).First(&system.SysAssetsListener{}).Error
				//fmt.Printf("listenerName: %s, err: %v\n", listenerName, err)
				// 不存在则创建该监听器
				if errors.Is(err, gorm.ErrRecordNotFound) {
					var listenerPort int
					var listener system.SysAssetsListener
					var listenerId string

					if len(lb.Listener) <= 0 {
						listenerPort = 10000
					} else {
						currentPort := lb.Listener[0].Port
						for _, num := range lb.Listener {
							if num.Port > currentPort {
								currentPort = num.Port
							}
						}
						listenerPort = currentPort + 1
					}

					// 定义请求参数
					lbRequestParams := cloudRequest.Listener{
						AK:              lb.CloudProduce.SecretId,
						SK:              lb.CloudProduce.SecretKey,
						Region:          lb.CloudProduce.RegionId,
						LbId:            lb.InstanceId,
						ListenerName:    listenerName,
						ListenerPort:    int32(listenerPort),
						Protocol:        "TCP",
						BackendPollName: fmt.Sprintf("%s-%d", gameServer.Host.PrivateIp, gameServer.TcpPort),
						BackendAddr:     gameServer.Host.PrivateIp,
						BackendPort:     int32(gameServer.TcpPort),
						SubnetCidrId:    lb.SubnetCidrId,
					}

					listenerId, err = cloud.CreateListener(lbRequestParams)
					if err != nil {
						return fmt.Errorf("%v, lbId: %s", err, lb.InstanceId)
					}

					listener.ProjectId = lb.ProjectId
					listener.InstanceId = listenerId
					listener.Name = listenerName
					listener.Port = listenerPort
					listener.LbId = lb.ID
					listener.Protocol = "TCP"
					listener.BackendIp = gameServer.Host.PrivateIp
					listener.BackendPort = int(gameServer.TcpPort)
					listener.HostId = gameServer.HostId

					if err = global.OPS_DB.WithContext(ctx).Create(&listener).Error; err != nil {
						return fmt.Errorf("创建负载均衡监听失败: %v", err)
					}

				} else if err != nil {
					global.OPS_LOG.Error("获取负载均衡监听信息失败:", zap.Error(err))
					return errors.New("获取负载均衡监听信息失败")
				}

			}
		}
	}

	return
}
func (s *AssetsLbService) WriteLBDataIntoRedis(ctx *gin.Context) (err error) {

	var lbList []system.SysAssetsLb

	redisKey := "LBDATA"

	resultMap := make(map[string]map[string][]string, 0)

	if err = global.OPS_DB.WithContext(ctx).Preload("Listener").Find(&lbList).Error; err != nil {
		global.OPS_LOG.Error("获取负载均衡信息失败:", zap.Error(err))
		return errors.New("获取负载均衡信息失败")
	}

	for _, lb := range lbList {
		// 根据 listener.BackendIp 和 listener.BackendPort 生成 key 和 value
		for _, listener := range lb.Listener {
			key := fmt.Sprintf("%s:%d", listener.BackendIp, listener.BackendPort)
			level := strconv.Itoa(int(lb.Level))
			value := fmt.Sprintf("%s:%d", lb.PubIp, listener.Port)

			if _, ok := resultMap[key]; !ok {
				resultMap[key] = make(map[string][]string)
			}
			if _, ok := resultMap[key][level]; !ok {
				resultMap[key][level] = make([]string, 0)
			}
			resultMap[key][level] = append(resultMap[key][level], value)
		}
	}

	var platformList []system.SysGamePlatform
	if err = global.OPS_DB.WithContext(ctx).Find(&platformList).Error; err != nil {
		global.OPS_LOG.Error("获取游戏平台信息失败:", zap.Error(err))
		return errors.New("获取游戏平台信息失败")
	}

	for _, platform := range platformList {
		var assetsRedis system.SysAssetsRedis
		if err = global.OPS_DB.WithContext(ctx).Where("platform_id = ?", platform.ID).First(&assetsRedis).Error; err != nil {
			global.OPS_LOG.Error("获取redis信息失败:", zap.Error(err))
			return errors.New("获取redis信息失败")
		}

		redisConn, err := utils.NewRedisConn(utils.RedisConfig{
			Addr:      fmt.Sprintf("%s:%d", assetsRedis.Host, assetsRedis.Port),
			DB:        0,
			Password:  assetsRedis.Password,
			IsCluster: assetsRedis.IsCluster,
		})

		if err != nil {
			return errors.New("redis连接失败")
		}
		for key, value := range resultMap {
			data, _ := json.Marshal(value)
			redisConn.HSet(context.Background(), redisKey, key, string(data))
		}

		if err := redisConn.Close(); err != nil {
			global.OPS_LOG.Error("redis关闭失败:", zap.Error(err))
		}
	}
	return
}
