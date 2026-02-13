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

internal/handler (控制层 / 接口层)
职责：“门面担当”。只负责处理 HTTP 请求和响应，不处理业务逻辑。

具体工作：

参数解析：从请求中提取参数（c.ShouldBindJSON, c.Param, c.Query）。

参数校验：检查必填项、格式等（使用 validator 或手动检查）。

调用 Service：将解析后的参数传给 Service 层去处理。

响应格式化：根据 Service 返回的结果或错误，封装成统一的 JSON 格式（如 code, msg, data）返回给前端。

例子：

menu.go: 接收前端发来的“创建菜单”请求，检查参数是否合法，然后调用 menuService.Create()。

internal/service (业务逻辑层)
职责：“大脑”。系统的核心逻辑都在这里，不碰 HTTP，也不直接碰数据库底层。

具体工作：

业务规则：例如“删除菜单前，必须检查是否有子菜单”、“注册时密码需要加密”、“下单前检查库存”。

数据组装：如果需要从多个 Data 来源取数据（比如先查 User 表，再查 Role 表），在这里进行组合。

事务控制：如果一个操作涉及多步数据库变更，在这里开启和提交事务。

例子：

sys_menu.go: Create() 方法里，它会先判断父级菜单是否存在，然后设置默认值，最后调用 data.CreateMenu()。

internal/data (数据访问层 / Repository)
职责：“仓库管理员”。只负责存取数据，不懂业务。

具体工作：

数据库操作：编写 SQL 或使用 GORM 进行 CRUD（增删改查）。

缓存操作：操作 Redis。

模型定义：定义 struct 与数据库表的映射关系（Model）。

例子：

sys_menu.go: 提供 GetList(), Create(), Delete() 等原子方法，直接操作数据库。

2. 其他重要文件夹
   cmd/server
   职责：“启动入口”。

内容：包含 main.go。负责加载配置、初始化日志、连接数据库、注入依赖（Wire 或手动）、启动 HTTP 服务、监听信号（优雅停机）。这里不写任何业务逻辑。

internal/router
职责：“路由中心”。

内容：将 URL 路径（如 /api/v1/menus）与具体的 Handler 函数绑定起来。它也负责挂载中间件（Middleware）。在你的新架构中，它主要负责顶层的路由分组。

internal/middleware
职责：“拦截器 / 插件”。处理横切关注点（Cross-cutting concerns）。

内容：

鉴权 (Auth)：检查 Token 是否有效 (jwt.go)。

日志 (Logger)：记录每个请求的耗时、IP、状态码 (access_log.go)。

限流 (RateLimit)：防止接口被刷爆。

跨域 (CORS)：允许前端跨域访问。

恢复 (Recovery)：防止 panic 导致整个程序崩溃。

internal/conf
职责：“配置定义”。

内容：定义配置文件的结构体（Struct），例如 Database, Redis, Server 等结构。负责读取 config.yaml 文件并映射到结构体中。

internal/response
职责：“统一响应”。

内容：定义统一的返回格式（如 Result 结构体），以及常用的辅助函数 Success(c, data), Fail(c, code, msg)。确保所有接口返回的数据结构一致。

internal/utils (或 pkg)
职责：“工具箱”。

内容：通用的、与业务无关的函数。

例如：MD5 加密、生成 UUID、时间格式化、文件操作等。

注意：不要把业务逻辑塞到这里，否则会变成“垃圾桶”。
