# SMS 订单查询功能

## Context
用户需要在首页添加一个入口，点击进入一个公开的查询页面。用户输入订单号，查询对应的虚拟手机号和短信内容。这是独立于 sub2api 核心功能的业务模块，数据存在本项目 PostgreSQL 中。

## 数据模型

新建 `sms_orders` 表，一个订单对应一条短信：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 (auto) | 主键 |
| order_no | string (unique) | 订单号，查询用 |
| phone_number | string | 虚拟手机号 |
| sms_content | text | 短信内容 |
| status | string | 状态（pending/received/expired） |
| created_at | timestamptz | 创建时间 |
| updated_at | timestamptz | 更新时间 |

## 后端改动（共 10 个文件）

### 1. Ent Schema
- **新建** `backend/ent/schema/sms_order.go`
  - 使用 `mixins.TimeMixin` 提供 created_at/updated_at
  - 字段：order_no、phone_number、sms_content、status
  - 索引：order_no（唯一）

### 2. 代码生成
- `cd backend && go generate ./ent && go generate ./cmd/server`

### 3. Service 层（领域模型 + 接口）
- **新建** `backend/internal/service/sms_order.go`
  - `SmsOrder` 结构体
  - `SmsOrderRepository` 接口（`GetByOrderNo`）
  - `SmsOrderService` + `NewSmsOrderService`
  - 方法：`GetByOrderNo(ctx, orderNo) (*SmsOrder, error)`

### 4. Repository 层
- **新建** `backend/internal/repository/sms_order_repo.go`
  - 实现 `SmsOrderRepository` 接口
  - `NewSmsOrderRepository` 返回 `service.SmsOrderRepository`
  - Entity→Service 映射函数

### 5. Handler 层
- **新建** `backend/internal/handler/sms_order_handler.go`
  - `SmsOrderHandler` + `NewSmsOrderHandler`
  - `Query` 方法：`GET /api/v1/sms-orders/:order_no`（公开，无需认证）

### 6. DTO
- **新建** `backend/internal/handler/dto/sms_order.go`
  - `SmsOrderResponse` 结构体（JSON 响应）
  - `SmsOrderFromService` 映射函数

### 7. 路由注册
- **修改** `backend/internal/server/routes/auth.go`
  - 在公开路由区域添加：`v1.GET("/sms-orders/:order_no", h.SmsOrder.Query)`

### 8. DI 注册
- **修改** `backend/internal/repository/wire.go` — 添加 `NewSmsOrderRepository`
- **修改** `backend/internal/service/wire.go` — 添加 `NewSmsOrderService`
- **修改** `backend/internal/handler/wire.go` — 添加 `NewSmsOrderHandler`，扩展 `Handlers` 结构体和 `ProvideHandlers`
- **修改** `backend/internal/handler/handler.go` — `Handlers` 结构体添加 `SmsOrder` 字段

### 9. Wire 代码生成
- `cd backend && go generate ./cmd/server`

## 前端改动（共 5 个文件）

### 1. 查询页面
- **新建** `frontend/src/views/SmsQueryView.vue`
  - 公开页面，布局参考 ModelsView.vue（独立页面，不带侧边栏）
  - 输入框：订单号
  - 查询按钮
  - 结果展示：手机号、短信内容、状态、时间
  - 空状态和错误处理

### 2. API 模块
- **新建** `frontend/src/api/smsQuery.ts`
  - `queryByOrderNo(orderNo: string)` 调用 `GET /api/v1/sms-orders/:order_no`

### 3. 路由
- **修改** `frontend/src/router/index.ts`
  - 添加 `/sms-query` 路由，`requiresAuth: false`

### 4. 首页入口
- **修改** `frontend/src/views/HomeView.vue`
  - 在 header 导航区域（Doc Link 旁边）添加一个查询入口图标/按钮，链接到 `/sms-query`

### 5. i18n
- **修改** `frontend/src/i18n/locales/zh.ts` 和 `en.ts`
  - 添加 `smsQuery` 相关翻译（页面标题、输入提示、按钮文字、状态文字等）
  - 添加 `nav.smsQuery` / `home.smsQuery` 导航文字

## 验证步骤

1. `cd backend && go generate ./ent && go generate ./cmd/server` — 生成代码
2. `cd backend && go build ./...` — 编译通过
3. `cd frontend && npx pnpm run typecheck` — 类型检查通过
4. `docker compose -f deploy/docker-compose.local.yml up -d --build` — Docker 构建成功
5. 访问首页，点击查询入口，输入订单号测试查询
