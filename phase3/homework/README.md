# Phase 3 实战作业：订单管理模块

## 作业要求

在 Phase 3 主项目的基础上，扩展一个**订单管理模块**，要求：

1. **领域模型**
   - `Order` 结构体，包含 ID、UserID、Amount、Status、Items、CreatedAt、UpdatedAt
   - `OrderStatus`：pending / paid / shipped / completed / cancelled

2. **数据访问层**
   - `OrderRepository` 接口
   - `Create`、`GetByID`、`ListByUser`、`UpdateStatus` 方法
   - 使用 GORM + SQLite（接入主项目数据库）

3. **业务逻辑层**
   - `OrderService`
   - `CreateOrder`、`GetOrder`、`ListUserOrders`、`PayOrder`、`CancelOrder`
   - 状态转换校验（如 paid 不能直接转 completed）

4. **HTTP Handler**
   - `POST /api/v1/orders`：创建订单（需 JWT）
   - `GET /api/v1/orders/:id`：查询订单（需 JWT）
   - `GET /api/v1/orders`：查询当前用户订单列表（需 JWT）
   - `PUT /api/v1/orders/:id/status`：更新订单状态（需 JWT，仅 admin 或订单所有者）

5. **单元测试**
   - 覆盖正常流程、状态转换、非法输入
   - 使用 testify

## 运行方式

```bash
cd phase3/homework

# 运行命令行演示
go run .

# 运行单元测试
go test -v
```

## 参考答案

本目录下的 `order.go`、`main.go`、`order_test.go` 提供了一份简化版参考答案，使用内存 repository 便于理解。

要在主项目中真正接入订单模块，需要：
1. 在 `internal/domain/` 添加 `order.go`
2. 在 `internal/repository/` 添加 `order_repository.go`
3. 在 `internal/service/` 添加 `order_service.go`
4. 在 `internal/handler/` 添加 `order_handler.go`
5. 在 `internal/router/router.go` 注册订单路由
6. 在 `cmd/server/main.go` 的 `AutoMigrate` 中加入 `Order` 模型

## 扩展挑战

1. 订单创建时扣减库存（引入 Product 模型和事务）
2. 订单超时自动取消（使用 goroutine + timer）
3. 订单支付接口对接第三方支付 mock
4. 给订单列表增加分页和排序
