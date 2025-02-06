package elb

import (
	elb "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/region"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/utils/cloud/hw_cloud"
)

type HwElb struct {
}

func NewHwElb(hw hw_cloud.HWCloud) *elb.ElbClient {
	basicAuth, err := hw_cloud.NewHwCloud(hw)
	if err != nil {
		global.OPS_LOG.Error("华为云AK认证失败:", zap.Error(err))
		return nil
	}
	_region, _ := region.SafeValueOf(hw.Region)

	auth, err := elb.ElbClientBuilder().WithRegion(_region).WithCredential(basicAuth).SafeBuild()
	if err != nil {
		global.OPS_LOG.Error("华为云ELB初始化失败:", zap.Error(err))
		return nil
	}

	client := elb.NewElbClient(auth)
	return client
}
