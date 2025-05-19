## 1.项目介绍

>`多项目管理运维后台` 是基于 Vue 和 gin 开发的全栈前后端分离的运维后台， 集成了jwt鉴权， 动态路由， 动态菜单， cmdb管理， 游戏服业务管理， 定时任务及异步任务管理等功能。

参考文档: 

后端:`https://github.com/flipped-aurora/gin-vue-admin`

前端:`https://github.com/youlaitech/vue3-element-admin-thin`

## 2. 使用说明

```
- node 版本 >=v20.16.0
- golang版本 >= v1.22.4
- 开发工具推荐: goland、vscode
```

### 2.1 server项目

使用`Goland`等编辑工具打开

```
# 克隆代码
git clone https://github.com/hrh-087/ops_server.git
# 切换目录
cd ops_server
# 安装相应的包
go mod tidy
# 编译代码
go build -ldflags="-s -w" -o server
# 初始化后台数据
./server -c config.yaml -o initData
# 运行(请修改相应的mysql，reids，asynq，prometheus配置)
./server -c config.yaml 
# 运行异步任务
./server -c  config.yaml -o worker
# 运行定时任务
./server -c config.yaml -o scheduler

# 账号密码
账号: admin
密码: dianchu666
```



### 2.2 web项目

使用`vscode`等编辑工具打开

```
git clone https://github.com/hrh-087/ops-web.git
# 安装 pnpm
npm install pnpm -g
# 设置镜像源(可忽略)
pnpm config set registry https://registry.npmmirror.com
# 安装依赖
pnpm install
# 启动运行
pnpm run dev
# 项目打包
pnpm run build
# 上传文件至远程服务器
将本地打包生成的 dist 目录下的所有文件拷贝至服务器的 /usr/share/nginx/html 目录
# nginx.cofig 配置
server {
	listen     80;
	server_name  localhost;
	location / {
			root /usr/share/nginx/html;
			index index.html index.htm;
	}
	# 反向代理配置
	location /prod-api/ {
            # vapi.youlai.tech 替换后端API地址，注意保留后面的斜杠 /
            proxy_pass http://vapi.youlai.tech/; 
	}
}


```

