package main

import (
	"context"
	"fmt"
	"log"

	gin_mcp "github.com/TIANLI0/gin-mcp"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	r := gin.Default()

	// 创建MCP处理器
	mcpHandler := gin_mcp.NewMCPHandler("example-server", "1.0.0",
		gin_mcp.WithBasePath("/api/mcp"),
		gin_mcp.WithServerOptions(
			server.WithToolCapabilities(true),
			server.WithPromptCapabilities(true),
		),
	)

	mcpHandler.AddTool(
		mcp.NewTool("hello",
			mcp.WithDescription("向指定名称问好"),
			mcp.WithString("name", mcp.Description("要问候的名称"), mcp.Required()),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			return mcp.NewToolResultText(fmt.Sprintf("你好, %s!", name)), nil
		},
	)

	// 注册MCP路由到Gin
	mcpHandler.Register(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 启动服务器
	log.Println("服务器运行在 http://localhost:8080")
	log.Println("MCP SSE端点: http://localhost:8080/api/mcp/sse")
	log.Println("MCP 消息端点: http://localhost:8080/api/mcp/message")
	r.Run(":8080")
}
