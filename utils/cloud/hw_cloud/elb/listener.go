package elb

import (
	elb "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
	"go.uber.org/zap"
	"ops-server/global"
)

// CreateLbListener 创建负载均衡监听器。
// 参数:
//
//	client: elb客户端，用于发送创建监听器的请求。
//	lbId: 负载均衡器的ID，用于关联监听器与负载均衡器。
//	protocol: 监听器使用的协议，决定了监听器如何处理入站流量。
//	name: 监听器的名称，用于标识特定的监听器。
//	protocolPort: 监听器监听的端口，决定了监听器监听的网络流量。
//
// 返回值:
//
//	*model.CreateListenerResponse: 创建监听器的响应对象，包含创建结果。
//	error: 错误对象，如果执行过程中发生错误则返回。
func (h HwElb) CreateLbListener(client *elb.ElbClient, lbId, protocol, name string, protocolPort int32) (*model.CreateListenerResponse, error) {

	transparentClientIpEnable := true

	request := &model.CreateListenerRequest{}
	request.Body = &model.CreateListenerRequestBody{
		Listener: &model.CreateListenerOption{
			LoadbalancerId:            lbId,
			Protocol:                  protocol,
			Name:                      &name,
			ProtocolPort:              &protocolPort,
			TransparentClientIpEnable: &transparentClientIpEnable,
		},
	}
	response, err := client.CreateListener(request)
	if err != nil {
		global.OPS_LOG.Error("创建负载均衡监听失败:", zap.Error(err), zap.String("lbId", lbId), zap.String("listener_name", name))
		return nil, err
	}
	return response, err
}

// DeleteListenerForce 删除负载均衡监听器。
//
// client: elb客户端，用于发送删除请求。
// listenerId: 监听器的ID，用于指定要删除的监听器。
// 返回error类型，表示删除操作可能遇到的错误。
// 该函数通过构造删除监听器的请求，并调用客户端的DeleteListenerForce方法来执行删除操作。
// 如果删除操作失败，会记录错误日志，并返回错误信息。
func (h HwElb) DeleteListenerForce(client *elb.ElbClient, listenerId string) error {
	request := &model.DeleteListenerForceRequest{}
	request.ListenerId = listenerId

	_, err := client.DeleteListenerForce(request)
	if err != nil {
		global.OPS_LOG.Error("删除负载均衡监听失败:", zap.Error(err), zap.String("listenerId", listenerId))
		return err
	}
	return err
}
