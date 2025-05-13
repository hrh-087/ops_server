package ecs

import (
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	"go.uber.org/zap"
	"ops-server/global"
)

// 用于查询Ecs 的vpc信息

// GetEcsVpcInfo 获取指定ECS实例的VPC信息
//
// 参数:
//   - client: *ecs.EcsClient, ECS客户端实例，用于与ECS服务进行交互
//   - instanceId: string, 需要查询的ECS实例ID
//
// 返回值:
//   - vpcInfo: *model.ListServerInterfacesResponse, 返回的ECS实例的VPC信息
//   - err: error, 如果查询过程中发生错误，返回错误信息
func (h HwEcs) GetEcsVpcInfo(client *ecs.EcsClient, instanceId string) (vpcInfo *model.ListServerInterfacesResponse, err error) {
	request := &model.ListServerInterfacesRequest{
		ServerId: instanceId,
	}

	response, err := client.ListServerInterfaces(request)
	if err != nil {
		global.OPS_LOG.Error("查询ecs vpc信息失败:", zap.Error(err))
		return nil, err
	}

	return response, err
}
