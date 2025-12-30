# JuggleIM

一个高性能，可扩展的开源 IM 即时通讯系统。

<p align="center">
<img align="left" height="110" src="./docs/logo.png">
<ul>
<li><strong>官网</strong>: https://www.juggle.im</li>
<li><strong>提问</strong>: https://github.com/juggleim/im-server/issues</li>
<li><strong>文档</strong>: https://www.juggle.im/docs/guide/intro/</li>
</ul>
</p>

[![](https://img.shields.io/github/license/juggleim/im-server?color=yellow&style=flat-square)](./LICENSE)
[![](https://img.shields.io/badge/go-%3E%3D1.20-30dff3?style=flat-square&logo=go)](https://github.com/juggleim/im-server)

## 社群讨论

如果对 IM 感兴趣、有集成问题讨论的朋友，非常欢迎加入社群讨论～

[Telegram 中文群](https://t.me/juggleim_zh)

## 特性

* 灵活的部署模式，支持公有云，私有云，托管云等部署形态。
* Protobuf+Websocket 实现长连接，低流量，高性能，且在网络不佳的环境下具备较好的连通性。
* 性能强大，专业版支持集群部署，无限横向扩展，能支撑亿级日活应用。
* 协议及数据全链路加密，无数据泄露风险。
* 提供方便的运维工具和管理后台，简单好维护。
* 支持Android，iOS，Web，PC 等多平台 SDK，提供快捷集成的Demo和文档。
* 支持多端同时在线和消息多端同步，确保状态多端一致。
* 支持全球链路加速，可服务全球级应用。
* 提供丰富的API和WebHook，可方便的与现有系统集成。
* 支持万人，十万人大群，轻松沟通，不丢消息。
* 支持无上限直播聊天室。
* 具备 AI 机器人对接能力，可轻松对接大模型。

## JuggleIM 项目说明

对接文档：[https://juggle.im/docs/client/import](https://juggle.im/docs/client/import/)

|  仓库地址 | 说明 |
| :--------- | :----- |
| [im-server](https://github.com/juggleim/im-server/) | 底层 IM 核心服务，负责消息分发，存储等IM相关业务 |
| [jugglechat-server](https://github.com/juggleim/jugglechat-server) | Demo的业务服务，负责处理用户注册/登录，创建群组，添加好友等业务，可以在这个基础上二开自己特色的业务能力 |
| [jugglechat-server-java](https://github.com/juggleim/jugglechat-server-java)| Demo 业务服务的 Java 版本 | 
| [imserver-console](https://github.com/juggleim/imserver-console) | IM 服务的管理后台，用于操作IM相关配置，监控IM业务量 |
| [imsdk-android](https://github.com/juggleim/imsdk-android) | 安卓端 imsdk，内含一个 UI Demo，可用于二开 |
| [imsdk-ios](https://github.com/juggleim/imsdk-ios) | iOS 端 imsdk，内含一个 UI Demo，可用于二开|
| [imsdk-web](https://github.com/juggleim/imsdk-web) | web 端 imsdk |
| imsdk-pc | 桌面端 imsdk，暂未开源，可联系客服了解 |
| [imsdk-flutter](https://github.com/juggleim/imsdk-flutter)| imsdk 的 flutter 版本 |
| [imsdk-harmony](https://github.com/juggleim/imsdk-harmony) | 鸿蒙版本 imsdk，内含一个 UI Demo，可用于二开 |
| [jugglechat-web](https://github.com/juggleim/jugglechat-web) | 集成 imsdk-web 的 web 版 Demo，可用于二开 |
| [jugglechat-desktop](https://github.com/juggleim/jugglechat-desktop) | 集成 imsdk-pc 的桌面版 Demo，可用于二开 |
| [jugglelive-web](https://github.com/juggleim/jugglelive-web)| 集成 imsdk-web 的一个聊天室场景Demo，可用于二开 |

其他：
| 仓库地址 | 说明 |
| :------ | :----- |
| [bot-connector](https://github.com/juggleim/bot-connector) | 机器人对接服务，用于打通 im-server 和 三方机器人 | 
| [imserver-sdk-go](https://github.com/juggleim/imserver-sdk-go) | 封装 im-server 服务端 API 的 SDK，供业务方集成到自己业务系统中 |
| [imserver-sdk-java](https://github.com/juggleim/imserver-sdk-java) | imserver-sdk 的 java 版本|



## 快速部署体验

部署文档(https://www.juggle.im/docs/guide/deploy/quickdeploy/)

## 手动部署

### 1. 安装并初始化 MySQL

#### 1) 安装 MySQL
略

#### 2) 创建DB实例
```
CREATE SCHEMA `jim_db` ;
```

#### 3) 初始化表结构
初始化表结构的sql文件在  im-server/docs/jim.sql , 导入命令如下：
```
mysql -u{db_user} -p{db_password} jim_db < jim.sql
```

### 2. 安装MongoDB(可选)
略

### 3. 启动im-server

#### 1) 运行目录
运行目录为 im-server/launcher，其中 conf 目录下存放的是配置文件，logs目录下是服务的运行日志目录。

#### 2) 编辑配置文件

配置文件位置：im-server/launcher/conf/config.yml
```
defaultPort: 9003       # im-server 默认监听端口
nodeName: testNode      # im-server 的节点名称
nodeHost: 127.0.0.1     # im-server 的节点IP
msgStoreEngine: mysql   # 配置用什么存储来存消息数据，有两种存储引擎可选。mysql：使用mysql存储消息数据(默认)；mongo：使用MongoDB存储消息数据

log:
  logPath: ./logs       # 运行日志所在目录
  logName: jim-info     # 运行日志的前缀名
  visual: false         # 是否开启可视化日志。开启后，会同步将日志数据写入一个 KV 数据库，在管理后台”开发工具->连接排查“处，可界面化查询日志；

mysql:                  # im-server 所用的MySQL相关配置
  user: root
  password: 123456
  address: 127.0.0.1:3306
  name: im_db

# mongodb:                # im-server 所用的MongoDB相关配置，用于存储消息数据。该配置为可选，在 msgStoreEngine 配置为 "mongo" 时生效；
#   address: 127.0.0.1:27017
#   name: jim_msgs        # mongodb 表空间名称，im-server启动后，会自动在这个空间下初始化collection；

# apiGateway:             # im-server 的服务端 API 端口, 供业务APP的服务端调用；非必填项，默认复用 defaultPort 作为默认端口
#   httpPort: 9001

# connectManager:         # im-server 长连接端口；非必填项，默认复用 defaultPort 作为默认端口
#   wsPort: 9003

adminGateway:           # im-server 自带的管理后台地址，默认账号密码是：admin/123456
  httpPort: 8090
```

#### 3) 启动im-server

在 im-server/launcher 目录下，执行如下命令：
```
go run main.go
```

#### 4) 配置外网访问地址(域名/IP)

需要配置外网地址的端口列表：
| 端口 | 协议类型 | 说明 | 
| ----:|:-----:|:----|
| 9003| http | 服务端 API 服务监听端口，用于业务服务器，例如 jugglechat-server 配置文件中需要配置这个地址；|
| 9003| websocket | IM 长连接监听端口，用于客户端SDK与IM 服务建立长连接(websocket) |
| 8090| http | IM 服务的管理后台监听端口，默认账号和密码：admin/123456 |

配置外网地址的方法，这里不详细描述，大家可以根据自身环境来灵活配置(常用方式: 挂公网ip，nginx反向代理，负载均衡等).

注： 如果仅内网调试使用，可以不配置外网IP/域名，仅使用内网IP即可

#### 5) 将长连接地址配置到IM系统中

配置方式很简单，在数据库中插入一条配置数据即可：

```
insert into globalconfs (conf_key, conf_value)values('connect_address', '{"default":["127.0.0.1:9002"]}')
```
其中，将 127.0.0.1 替换成该机器的内网IP，或对外的公网IP/域名，这个是客户端SDK的长连接地址，将有导航服务(8081)下发给客户端SDK；

### 4. 创建应用(租户)
JuggleIM 本身是一套多租户的系统，可以在一套服务中创建多个appkey(租户)，租户之间的数据相互隔离，互不影响。

#### 方式一：登录管理后台，创建租户
待完善

#### 方式二：通过管理API，创建租户

其中，app_key 用于指定租户的标识，可自定义，要求在系统内唯一；app_name为租户的名称，可自定义； 
注：这里用的是IM服务管理后台(8090)的地址，127.0.0.1 替换成im服务的内网IP，或公网IP/域名。

```
curl --request POST \
  --url http://127.0.0.1:8090/admingateway/apps/create \
  --data '{
    "app_key":"appkey",
    "app_name":"appname"
}'
```

响应数据示例：

```
{
	"code": 0,
	"msg": "success",
	"data": {
		"app_name": "appname",
		"app_key": "appkey",
		"app_secret": "hciKcc6sXRDjYUQp"
	}
}
```
### 5. 登入管理后台

管理后台地址：http://127.0.0.1:8090    默认账号/密码： admin/123456

登入管理后台后，即可看到创建的应用列表，点击其中一个应用，可以对其配置进行修改和维护。

### 6. 业务服务器/客户端集成

这里汇总下业务集成所需的各项配置：

1) 业务服务器集成

| 配置项 | 示例|备注 | 
| ----: |:-----:|:----|
|IM 服务端 API 地址 |  http://127.0.0.1:9003 | 供业务服务器访问IM服务的API接口地址，使用该接口可以注册IM用户，创建群，发送系统消息等，接口文档参考：https://www.juggle.im/docs/server/api/|
|app_key | appkey1 |应用的租户标识，在第4步中创建，可自定义，但要保证在系统内唯一 |
|app_secret| hciKcc6sXRDjYUQp | 应用对应的鉴权秘钥，创建应用时自动生成。如果想自定义的话，确保配置为16位的字符串。注意：确保该秘钥仅在业务服务器端使用，不要泄露到客户端。 |

2) 客户端SDK集成

| 配置项 | 示例|备注 | 
| ----: |:-----:|:----|
|IM 服务的连接地址| ws://127.0.0.1:9003| IM 的连接地址，客户端SDK初始化时需要传入，参考文档：https://www.juggle.im/docs/client/quickstart/android/ |
|app_key|appkey1|应用的租户标识，确保与业务服务器端配置的保持一致。|
