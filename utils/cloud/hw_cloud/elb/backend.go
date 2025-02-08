package elb

import (
	elb "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
	"go.uber.org/zap"
	"ops-server/global"
)

// ListBackendPolls 查询后端监听池信息。
// 该方法根据提供的poolId来查询后端监听池的详细信息。
// 参数:
//
//	client *elb.ElbClient: ELB服务的客户端，用于发送请求。
//	listenerIid string: 监听器的ID，用于指定查询的绑定改监听器的后端服务器组。
//
// 返回值:
//
//	*model.ListPoolsResponse: 包含后端监听池信息的响应对象。
//	error: 查询过程中发生的错误，如果查询成功则为nil。
func (h HwElb) ListBackendPolls(client *elb.ElbClient, listenerIid string) (*model.ListPoolsResponse, error) {

	var ids []string

	if listenerIid != "" {
		ids = append(ids, listenerIid)
	}
	request := &model.ListPoolsRequest{
		ListenerId: &ids,
	}
	response, err := client.ListPools(request)
	if err != nil {
		global.OPS_LOG.Error("查询后端监听池失败:", zap.Error(err))
		return nil, err
	}

	return response, err
}

// CreateBackendPoll 创建负载均衡的后端监听池
// 参数:
//
//	client: elb客户端，用于发送创建后端监听池的请求
//	listenerId: 监听器ID，关联后端监听池与特定监听器
//	name: 后端监听池的名称
//
// 返回值:
//
//	*model.CreatePoolResponse: 创建后端监听池的响应对象，包含创建结果
//	error: 错误对象，如果创建过程中发生错误
func (h HwElb) CreateBackendPoll(client *elb.ElbClient, listenerId string, name string) (*model.CreatePoolResponse, error) {
	request := &model.CreatePoolRequest{}
	request.Body = &model.CreatePoolRequestBody{
		Pool: &model.CreatePoolOption{
			LbAlgorithm: "ROUND_ROBIN", // 使用加权轮询
			ListenerId:  &listenerId,
			Name:        &name,
			Protocol:    "TCP",
		},
	}

	response, err := client.CreatePool(request)
	if err != nil {
		global.OPS_LOG.Error("创建负载均衡后端失败:", zap.Error(err), zap.String("listenerId", listenerId), zap.String("name", name))
		return nil, err
	}
	return response, err
}

// CreateBackendMember 创建负载均衡的后端成员
// 参数:
//
//	client *elb.ElbClient: ELB客户端，用于调用ELB API
//	poolId string: 负载均衡池的ID
//	ipAddr string: 后端成员的IP地址
//	port int32: 后端成员的端口号
//
// 返回值:
//
//	*model.CreateMemberResponse: 创建后端成员的响应对象
//	error: 错误对象，如果调用API过程中发生错误
func (h HwElb) CreateBackendMember(client *elb.ElbClient, poolId string, ipAddr string, port int32, subnetCidrId string) (*model.CreateMemberResponse, error) {
	request := &model.CreateMemberRequest{}

	request.PoolId = poolId
	request.Body = &model.CreateMemberRequestBody{
		Member: &model.CreateMemberOption{
			Address:      ipAddr,
			ProtocolPort: &port,
			SubnetCidrId: &subnetCidrId,
		},
	}

	response, err := client.CreateMember(request)
	if err != nil {
		global.OPS_LOG.Error("创建负载均衡后端成员失败:", zap.Error(err), zap.String("poolId", poolId), zap.String("ipAddr", ipAddr), zap.Int32("port", port))
		return nil, err
	}
	return response, err
}

// BatchDeleteMembers 批量删除负载均衡的后端成员。
// 参数:
//
//	client (*elb.ElbClient): 用于调用负载均衡服务的客户端。
//	poolId (string): 负载均衡池的ID，用于指定要删除后端成员的负载均衡池。
//
// 返回值:
//
//	err (error): 删除操作中可能发生的错误。
func (h HwElb) BatchDeleteMembers(client *elb.ElbClient, poolId string, ids []model.MemberRef) (err error) {

	var memberIdList []model.BatchDeleteMembersOption
	request := &model.BatchDeleteMembersRequest{}
	request.PoolId = poolId

	if len(ids) > 0 {
		for _, memberId := range ids {
			memberIdList = append(memberIdList, model.BatchDeleteMembersOption{
				Id: &memberId.Id,
			})
		}
	}
	request.Body = &model.BatchDeleteMembersRequestBody{
		Members: memberIdList,
	}

	_, err = client.BatchDeleteMembers(request)
	if err != nil {
		global.OPS_LOG.Error("删除负载均衡后端成员失败:", zap.Error(err), zap.String("poolId", poolId))
		return
	}
	return err
}

// DeleteBackendPoll 删除负载均衡的后端监听池。
// 参数:
//
//	client (*elb.ElbClient): 用于调用ELB服务的客户端。
//	poolId (string): 要删除的后端监听池的ID。
//
// 返回值:
//
//	err (error): 删除操作可能遇到的错误，如果没有错误则为nil。
func (h HwElb) DeleteBackendPoll(client *elb.ElbClient, poolId string) (err error) {
	request := &model.DeletePoolRequest{}
	request.PoolId = poolId
	_, err = client.DeletePool(request)
	if err != nil {
		global.OPS_LOG.Error("删除负载均衡后端监听池失败:", zap.Error(err), zap.String("poolId", poolId))
		return
	}
	return err
}
