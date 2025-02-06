package elb

import (
	elb "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
)

// GetLbList 获取负载均衡列表
// 该方法通过指定的ELB客户端，分页获取负载均衡器列表
// 参数:
//
//	client *elb.ElbClient: ELB客户端，用于发送请求并获取数据
//	page string: 分页标记，表示从哪一页开始获取数据
//	limit int32: 每页数据的数量限制
//
// 返回值:
//
//	*model.ListLoadBalancersResponse: 包含负载均衡列表的响应对象
//	error: 错误信息，如果执行过程中遇到错误则返回
func (h HwElb) GetLbList(client *elb.ElbClient, page string, limit int32) (result *model.ListLoadBalancersResponse, err error) {
	if limit == 0 {
		limit = 2000
	}

	request := &model.ListLoadBalancersRequest{
		Limit: &limit,
	}

	if page != "" {
		request.Marker = &page
	}

	return client.ListLoadBalancers(request)
}
