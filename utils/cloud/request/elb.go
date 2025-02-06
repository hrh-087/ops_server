package request

type Listener struct {
	AK              string
	SK              string
	Region          string // 区域id
	LbId            string // 负载均衡id
	SubnetCidrId    string // vpc子网id
	ListenerId      string // 监听器id
	ListenerName    string // 监听器名称
	ListenerPort    int32  // 监听器端口
	Protocol        string // 监听器协议
	BackendPollId   string // 后端组id
	BackendPollName string // 后端组名称
	BackendAddr     string // 后端服务器ip
	BackendPort     int32  // 后端服务器端口
}
