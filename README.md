# gin-mcp

gin-mcp 是一个用于在 [Gin Web 框架](https://github.com/gin-gonic/gin) 中轻松集成 [Model Context Protocol (MCP)](https://github.com/mark3labs/mcp-go) 服务器的库。它允许您在 Gin 应用程序中轻松添加 MCP 功能，包括工具、提示和资源管理。

## 功能特性

- 在 Gin 中轻松集成 MCP 服务器
- 支持 SSE (Server-Sent Events) 传输
- 自定义路由和基础路径
- 会话管理和工具注册
- 基于 Gin 中间件的认证支持
- 自定义上下文函数支持

## 安装

```bash
go get github.com/TIANLI0/gin-mcp
```

## 快速开始

以下是一个简单的示例，展示如何在 Gin 应用程序中设置 MCP Server：

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/TIANLI0/gin-mcp"
    "github.com/gin-gonic/gin"
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    r := gin.Default()

    // 创建 MCP 处理器
    mcpHandler := gin_mcp.NewMCPHandler("example-server", "1.0.0",
        gin_mcp.WithBasePath("/api/mcp"),
        gin_mcp.WithServerOptions(
            server.WithToolCapabilities(true),
        ),
    )

    mcpHandler.AddTool(
        mcp.NewTool("hello",
            mcp.WithDescription("向指定名称问好"),
            mcp.WithString("name", mcp.Description("要问候的名称"), mcp.Required()),
        ),
        func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
            name, ok := request.Params.Arguments["name"].(string)
            if !ok {
                return nil, fmt.Errorf("参数无效: 需要字符串类型的name参数")
            }
            return mcp.NewToolResultText(fmt.Sprintf("你好, %s!", name)), nil
        },
    )

    mcpHandler.Register(r)

    // 启动服务器
    log.Println("服务器运行在 http://localhost:8080")
    log.Println("MCP SSE端点: http://localhost:8080/api/mcp/sse")
    log.Println("MCP 消息端点: http://localhost:8080/api/mcp/message")
    r.Run(":8080")
}
```

## 认证

您可以使用 `WithAuth` 选项为 MCP 端点添加认证：

```go
mcpHandler := gin_mcp.NewMCPHandler("auth-example-server", "1.0.0",
    gin_mcp.WithAuth(func(c *gin.Context) bool {
        // 简单的 API 密钥验证
        apiKey := c.GetHeader("X-API-Key")
        if apiKey != "your-secret-api-key" {
            c.AbortWithStatusJSON(401, gin.H{"error": "未授权"})
            return false
        }
        return true
    }),
)
```

## 会话管理

gin-mcp 支持会话特定工具，允许您向特定用户会话添加专用工具：

```go
// 添加全局工具
mcpHandler.AddTool(...)

// 向特定会话添加工具
mcpHandler.AddSessionTool("user-123", mcp.NewTool("private-tool"), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    return mcp.NewToolResultText("这是一个会话特定的工具"), nil
})
```

## 配置选项

gin-mcp 提供多种配置选项：

- `WithBasePath(path)`: 设置 MCP 端点的基础路径
- `WithSSERoute(route)`: 自定义 SSE 端点路径 
- `WithMessageRoute(route)`: 自定义消息端点路径
- `WithContextFunc(fn)`: 设置自定义上下文函数
- `WithServerOptions(opts...)`: 添加 MCP 服务器选项
- `WithSSEOptions(opts...)`: 添加 SSE 服务器选项
- `WithAuth(fn)`: 添加认证处理程序

## License

MIT License - 详情请参阅 [LICENSE](LICENSE) 文件。