package ecs

import (
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/region"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/utils/cloud/hw_cloud"
)

type HwEcs struct {
}

func NewHwEcs(hw hw_cloud.HWCloud) *ecs.EcsClient {
	basicAuth, err := hw_cloud.NewHwCloud(hw)
	if err != nil {
		global.OPS_LOG.Error("华为云AK认证失败:", zap.Error(err))
		return nil
	}
	_region, _ := region.SafeValueOf(hw.Region)

	auth, err := ecs.EcsClientBuilder().WithRegion(_region).WithCredential(basicAuth).SafeBuild()
	if err != nil {
		global.OPS_LOG.Error("华为云Ecs初始化失败:", zap.Error(err))
		return nil
	}

	return ecs.NewEcsClient(auth)
}
