package cloud

import (
	"errors"
	"ops-server/utils/cloud/hw_cloud"
	"ops-server/utils/cloud/hw_cloud/elb"
	"ops-server/utils/cloud/request"
)

func CreateListener(listener request.Listener) (string, error) {
	var hwElb elb.HwElb
	client := elb.NewHwElb(hw_cloud.HWCloud{
		AK:     listener.AK,
		SK:     listener.SK,
		Region: listener.Region,
	})

	if client == nil {
		return "", errors.New("初始化云商信息失败")
	}

	cloudListener, err := hwElb.CreateLbListener(
		client,
		listener.LbId,
		listener.Protocol,
		listener.ListenerName,
		listener.ListenerPort,
	)

	if err != nil {
		return "", errors.New("创建负载均衡监听失败")
	}
	// 获取监听器id
	listener.ListenerId = cloudListener.Listener.Id

	backendPoll, err := hwElb.CreateBackendPoll(client, listener.ListenerId, listener.BackendPollName)
	if err != nil {
		return "", errors.New("创建负载均衡后端服务器组失败")
	}
	// 获取后端服务器组id
	listener.BackendPollId = backendPoll.Pool.Id

	_, err = hwElb.CreateBackendMember(client, listener.BackendPollId, listener.BackendAddr, listener.BackendPort, listener.SubnetCidrId)
	if err != nil {
		return "", errors.New("创建负载均衡后端服务器失败")
	}

	return listener.ListenerId, err
}

func DeleteListener(listener request.Listener) (err error) {
	var hwElb elb.HwElb
	client := elb.NewHwElb(hw_cloud.HWCloud{
		AK:     listener.AK,
		SK:     listener.SK,
		Region: listener.Region,
	})

	if client == nil {
		return errors.New("初始化云商信息失败")
	}

	// 根据监听器id删除后端服务器
	backendList, err := hwElb.ListBackendPolls(client, listener.ListenerId)
	if err != nil {
		return errors.New("获取负载均衡后端服务器组失败")
	}

	for _, poll := range *backendList.Pools {
		if err = hwElb.BatchDeleteMembers(client, poll.Id, poll.Members); err != nil {
			return errors.New("删除负载均衡后端服务器失败")
		}
		if err = hwElb.DeleteBackendPoll(client, poll.Id); err != nil {
			return errors.New("删除负载均衡后端服务器组失败")
		}
	}

	if err = hwElb.DeleteListenerForce(client, listener.ListenerId); err != nil {
		return errors.New("删除负载均衡监听失败")
	}
	//global.OPS_LOG.Info("成功删除监听器:", zap.String("listenerId", listener.ListenerId))
	return
}
