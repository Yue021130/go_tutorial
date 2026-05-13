# 源码阅读：etcd

## 阅读目标

理解 etcd 作为分布式 KV 存储的核心模块：Raft 共识、WAL 日志、存储引擎、API 服务端。

## 推荐版本

- GitHub: https://github.com/etcd-io/etcd
- 建议阅读 v3.5.x 稳定版，从 `server/` 子模块入手

## 核心模块路线

```text
server/
├── etcdserver/         # etcd 服务端核心
│   ├── server.go       # EtcdServer 主循环、请求处理
│   ├── raft.go         # 与 raft 模块交互
│   └── api/            # v3 API 实现
├── mvcc/               # 多版本并发控制存储
│   ├── kvstore.go      # 键值存储
│   ├── watcher.go      # Watch 机制
│   └── index.go        # 索引
├── wal/                # 预写日志（Write Ahead Log）
├── snap/               # 快照管理
├── lease/              # 租约（TTL）
├── auth/               # 认证
└── embed/              # 嵌入式启动

raft/
├── raft.go             # Raft 状态机
├── log.go              # Raft 日志
├── node.go             # Node 接口
└── read_only_queue.go  # ReadIndex 只读请求队列
```

## 重点解析

### 1. etcdserver：请求入口

- `EtcdServer` 实现了 `etcdserverpb.KVServer`、`etcdserverpb.WatchServer` 等 gRPC 接口。
- 写请求（Put/Delete/Txn）通过 Raft 提案提交。
- 读请求默认走 `Range`，可配置为 `serializable`（不经过 Raft，快但可能读到旧值）或 `linearizable`（经过 ReadIndex，保证线性一致）。

### 2. Raft 模块

- etcd 使用自己实现的 Raft 库，位于 `raft/` 目录。
- 核心状态：Follower -> Candidate -> Leader。
- Leader 处理写请求，生成 Entry，通过 `Ready` 通道通知上层持久化并发送给 Followers。
- 多数派确认后 Entry 被提交，应用到状态机（mvcc）。

```text
Raft 状态机简化流程：

Client -> EtcdServer -> Raft.Propose -> WAL + Snapshot -> Apply -> mvcc
                              |
                              v
                       发送 AppendEntries 给其他节点
```

### 3. mvcc：多版本并发控制

- 每个 key 保存多个版本，用 revision（main + sub）标识。
- `kvstore` 维护当前索引和历史版本。
- `watchableStore` 实现 Watch，支持从指定 revision 开始监听。
- 后台 compaction 定期清理过旧版本。

### 4. WAL 与 Snapshot

- WAL 记录所有 Raft 日志，用于崩溃恢复。
- Snapshot 定期对 mvcc 状态做快照，避免日志无限增长。
- 启动时：加载最新 Snapshot + 回放后续 WAL。

### 5. Lease 租约

- 用于给 key 设置 TTL，常用于服务注册与发现。
- 租约由 Leader 统一续约（heartbeat），避免每个 key 单独计时。

## 阅读建议

1. 从 `server/embed/etcd.go` 看启动流程。
2. 跟踪一个 `Put` 请求：`grpc gateway -> EtcdServer.Put -> raftNode.Propose -> WAL -> apply -> mvcc`。
3. 阅读 `raft/raft.go` 的 `becomeLeader`、`Step`、`tick`。
4. 阅读 `mvcc/watchable_store.go` 的 `watch` 实现。

## 可迁移到自己的项目中的设计

- Raft 日志 + 快照的持久化与恢复模式。
- Watch 机制：事件通知 + revision 回溯。
- Lease 集中式心跳续约。
- 线性一致读 vs 串行读的选择策略。
