package hw_cloud

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
)

type HWCloud struct {
	AK     string
	SK     string
	Region string
}

func NewHwCloud(hw HWCloud) (*basic.Credentials, error) {
	return basic.NewCredentialsBuilder().WithAk(hw.AK).WithSk(hw.SK).SafeBuild()
}
