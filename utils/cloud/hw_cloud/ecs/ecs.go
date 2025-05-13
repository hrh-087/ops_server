package ecs

import (
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	"go.uber.org/zap"
	"ops-server/global"
)

func (h HwEcs) GetEcsInfo(client *ecs.EcsClient, instanceId string) (vpcInfo *model.ShowServerResponse, err error) {
	request := &model.ShowServerRequest{
		ServerId: instanceId,
	}

	response, err := client.ShowServer(request)
	if err != nil {
		global.OPS_LOG.Error("查询ecs信息失败:", zap.Error(err))
		return nil, err
	}

	return response, err
}
