internal/
├── conf/ # 配置定义
├── data/ # 数据访问层 (Database, Redis) - 实现 Repository 接口
├── service/ # 业务逻辑层 (Business Logic) - 处理核心业务，调用 data
├── handler/ # HTTP 接口层 (Gin Handlers) - 解析参数，调用 service，响应结果
│ ├── auth.go # 原 api/auth.go 移到这里
│ ├── menu.go # 原 handler/sys_menu.go 改名为 menu.go，保持风格统一
│ └── user.go
├── router/ # 路由注册
└── middleware/ # 中间件
