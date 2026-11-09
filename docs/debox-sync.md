# debox 通用能力同步到单机 dev

## 来源与边界

参考集群合并提交 `37dbf617` 以及 `debox`（核对到 `6a7121db`）中的原始提交，
并完成单机迁移链和活动统计时间单位复核修正。目标为单机 `dev`。

## 本次同步

- API 签名同时接受 `apps.app_secret` 和 `appexts.temporary_secret`；空密钥不参与验签。
- 加载应用扩展时过滤核心字段，防止覆盖主密钥、应用标识、状态等；保留原有 API 前回调用户白名单和订阅开关。
- 重复加群冷路径先批查请求成员；全部已存在时跳过总人数查询、人数限制检查与写入，返回幂等成功。增加 `cold_count_skipped` 指标。
- MySQL 迁移在同一物理连接上获取数据库级 advisory lock、执行 SQL、写入版本及释放锁；读取、SQL、版本写入和释放锁失败均返回错误，启动在迁移失败时退出。
  无版本记录或缺少版本表才走初始化；已存在的零值、负值、非法或超前版本拒绝启动。SQL 文件不再自行提前写版本。
- MySQL 初始化错误与迁移失败统一以 `os.Exit(1)` 退出。
- 核对全部 15 个嵌入迁移的建表、列和索引增删顺序：`20260720` 为已有 `appnavs.admin_url` 添加兼容定义守卫；`20260813` 为 `useractivities` 两个索引添加精确定义守卫，消除与 `20260617` 的重复 ADD。`pushtokens.idx_pushtoken` 同样使用精确定义守卫。已有不兼容同名定义仍报错，不能掩盖异常。
- 删除历史消息透传 `msg_time`、`sub_channel`，缺省时间由严格校验后的消息 ID 解析；处理空消息项。删除通知及群/私聊已读链路补齐子频道。
- 范围已读与阅后销毁查询统一使用实际存储字段 `msg_seq_no`（MySQL/Mongo），支持只传 `IndexScopes` 的请求。
- 群历史消息、定向成员关系和发送者删除记录在 MySQL 下事务写入；Mongo 保留原有写入模式。
- Mongo 群删除记录查询补齐子频道，并按实际的 `send_time` 字段排序。
- 新增独立 DAU CLI：`cmd/dau`，使用标准库 `flag`，不新增或升级依赖。
- 活动新写统一 UTC 日毫秒标记；仅活动写入改单位，原有消息日统计秒级存储/展示保持一致。
- 清理按双单位识别已结束日期，秒/毫秒同用户按 `COUNT(DISTINCT user_id)` 去重。日汇总覆盖写入而非累加，并移除对应旧秒级汇总；过期日先汇总再删除双单位明细，位于同一事务，任一步失败回滚。重复运行不重复计数，当天/未来日期不参与清理。

## 已存在，核对后复用

- API IP 白名单、请求来源 IP、批量未命中用户资料查询及用户类型保留。
- 会话目标用户资料批量加载、离线大群扇出避免冗余资料查询。
- HTTP interceptor 共享连接池、请求 context/取消处理。
- MySQL 空闲连接上限已为 300，不照搬集群连接寿命等调参。
- 撤回权限已比较 `SenderId`，扩展更新已有 `sub_channel=?` 条件。

## 明确排除

- 历史消息按月分区、历史 split/compat DAO、分区路由与双写切换。
- 月分区维护任务、分区 DDL、新增 `20261008/20261013/20261014` 历史迁移。
- `historyMsgPartition`、独立 `msgmysql` 的连接、配置、环境变量、模板与部署文档。
- 集群 RPC 背压隔离、ZK 重连注册、xxhash 路由、connectproxy 改动。
  单机 `commons/gmicro/actorsystem/msgsender.go` 仅将请求交给本地 receiver，没有远端 sender。
- debox 禁用服务/API、服务端 API 无条件跳过前回调等业务裁剪。
- 不修改现有 E2EE、设备 ban、metrics、标签、推送和后台实现；不更改原有消息 inbox/cmd 分表能力。
- 单机不支持 HBase，不新增该存储引擎。

## DAU 用法及单机适配

```sh
go run -mod=readonly ./cmd/dau --config launcher/conf/config.yml --app-key YOUR_APP_KEY
go run -mod=readonly ./cmd/dau --config launcher/conf/config.yml --app-key YOUR_APP_KEY --days 7
```

默认统计含今天的最近 30 天，按 UTC 日期倒序逐日输出。只读查询，不启动服务器或执行迁移。
默认查找 `conf/config.yml`、`launcher/conf/config.yml`。
活动新写使用毫秒标记，兼容已有秒记录；同一天秒/毫秒明细按用户去重。明细不存在时读取 `dailyactivities` 汇总，优先毫秒汇总，再回退旧秒汇总，不将二者相加。
维护任务补汇总所有仍有明细的已结束 UTC 日，保留最近 3 天明细，过期日按日事务归档删除。单日事务大小随该日明细量增长；大数据量下的执行时间和锁影响需实库评估。
结果依赖服务器已记录/汇总的数据；不会重建已清理且未汇总的历史明细。

## 验证与遗留

- `go build -mod=readonly ./...`、`git diff --check` 通过；`go.mod/go.sum` 无修改。
- `go test -mod=readonly`：DAU、dbcommons、tools、commonservices、group/services、historymsg、apigateway 通过。
- 扩展回归：launcher、message、conversation、actorsystem、metrics、connectmanager、usermanager 与推送服务测试通过。
- 现有推送存储测试依赖真实 MySQL/Mongo：MySQL 认证失败，Mongo 未启动，未作修改。
- 新增迁移脚本化 SQL driver 测试覆盖同连接加锁、失败停止、版本写入、释放锁与初始化判断，不依赖真实 DB。
- 新增全迁移列/索引清单重放、完整空库启动脚本、DDL 守卫生命周期测试；活动测试覆盖双单位查询边界、用户去重 SQL、覆盖汇总重试、失败回滚、保留期边界及当天/未来记录保护。
- 未进行真实 MySQL 的空库全迁移、DDL 守卫执行、事务回滚或生产数据验证。脚本化 driver 和静态迁移清单测试不能替代真实数据库验证。
