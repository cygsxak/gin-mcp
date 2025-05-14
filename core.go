package gin_mcp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCPHandler 是一个Gin中间件，用于集成MCP server
type MCPHandler struct {
	server     *server.MCPServer
	sseServer  *server.SSEServer
	basePath   string
	sseRoute   string
	msgRoute   string
	contextFn  server.HTTPContextFunc
	serverOpts []server.ServerOption
	sseOpts    []server.SSEOption
}

// MCPHandlerOption 是配置MCPHandler的函数选项
type MCPHandlerOption func(*MCPHandler)

// NewMCPHandler 创建一个新的MCP处理器，用于与Gin集成
func NewMCPHandler(name, version string, opts ...MCPHandlerOption) *MCPHandler {
	h := &MCPHandler{
		basePath:   "/mcp",
		sseRoute:   "/sse",
		msgRoute:   "/message",
		serverOpts: []server.ServerOption{},
		sseOpts:    []server.SSEOption{},
	}

	for _, opt := range opts {
		opt(h)
	}

	// 创建MCP服务器
	h.server = server.NewMCPServer(name, version, h.serverOpts...)

	return h
}

// Register 注册MCP路由到Gin引擎
func (h *MCPHandler) Register(r *gin.Engine) {
	// 确保基础路径以/开头且不以/结尾
	if h.basePath == "" {
		h.basePath = "/"
	} else if h.basePath[0] != '/' {
		h.basePath = "/" + h.basePath
	}
	if h.basePath != "/" && h.basePath[len(h.basePath)-1] == '/' {
		h.basePath = h.basePath[:len(h.basePath)-1]
	}

	// 创建SSE服务器
	sseOpts := append(h.sseOpts,
		server.WithStaticBasePath(h.basePath),
		server.WithSSEEndpoint(h.sseRoute),
		server.WithMessageEndpoint(h.msgRoute),
	)

	if h.contextFn != nil {
		sseOpts = append(sseOpts, server.WithHTTPContextFunc(h.contextFn))
	}

	h.sseServer = server.NewSSEServer(h.server, sseOpts...)

	// 注册路由
	r.GET(h.basePath+h.sseRoute, func(c *gin.Context) {
		h.sseServer.ServeHTTP(c.Writer, c.Request)
	})

	r.POST(h.basePath+h.msgRoute, func(c *gin.Context) {
		h.sseServer.ServeHTTP(c.Writer, c.Request)
	})
}

// GetServer 返回底层的MCP服务器实例
func (h *MCPHandler) GetServer() *server.MCPServer {
	return h.server
}

// GetSSEServer 返回底层的SSE服务器实例
func (h *MCPHandler) GetSSEServer() *server.SSEServer {
	return h.sseServer
}

// WithBasePath 设置MCP处理器的基础路径
func WithBasePath(path string) MCPHandlerOption {
	return func(h *MCPHandler) {
		h.basePath = path
	}
}

// WithSSERoute 设置SSE端点的路由
func WithSSERoute(route string) MCPHandlerOption {
	return func(h *MCPHandler) {
		h.sseRoute = route
	}
}

// WithMessageRoute 设置消息端点的路由
func WithMessageRoute(route string) MCPHandlerOption {
	return func(h *MCPHandler) {
		h.msgRoute = route
	}
}

// WithContextFunc 设置HTTP上下文函数
func WithContextFunc(fn server.HTTPContextFunc) MCPHandlerOption {
	return func(h *MCPHandler) {
		h.contextFn = fn
	}
}

// WithServerOptions 添加MCP服务器选项
func WithServerOptions(opts ...server.ServerOption) MCPHandlerOption {
	return func(h *MCPHandler) {
		h.serverOpts = append(h.serverOpts, opts...)
	}
}

// WithSSEOptions 添加SSE服务器选项
func WithSSEOptions(opts ...server.SSEOption) MCPHandlerOption {
	return func(h *MCPHandler) {
		h.sseOpts = append(h.sseOpts, opts...)
	}
}

// AuthHandlerFunc 定义认证处理函数
type AuthHandlerFunc func(c *gin.Context) bool

// WithAuth 添加基本认证中间件
func WithAuth(authFn AuthHandlerFunc) MCPHandlerOption {
	return func(h *MCPHandler) {
		h.contextFn = func(ctx context.Context, r *http.Request) context.Context {
			if gc, exists := r.Context().Value(gin.ContextKey).(*gin.Context); exists && authFn(gc) {
				return ctx
			}
			return ctx
		}
	}
}

// AddTool 向MCP服务器添加工具
func (h *MCPHandler) AddTool(tool mcp.Tool, handler server.ToolHandlerFunc) {
	h.server.AddTool(tool, handler)
}

// AddTools 向MCP服务器添加多个工具
func (h *MCPHandler) AddTools(tools ...server.ServerTool) {
	h.server.AddTools(tools...)
}

// AddSessionTool 向特定会话添加工具
func (h *MCPHandler) AddSessionTool(sessionID string, tool mcp.Tool, handler server.ToolHandlerFunc) error {
	return h.server.AddSessionTool(sessionID, tool, handler)
}

// AddSessionTools 向特定会话添加多个工具
func (h *MCPHandler) AddSessionTools(sessionID string, tools ...server.ServerTool) error {
	return h.server.AddSessionTools(sessionID, tools...)
}

// DeleteSessionTools 从特定会话删除工具
func (h *MCPHandler) DeleteSessionTools(sessionID string, names ...string) error {
	return h.server.DeleteSessionTools(sessionID, names...)
}

// SendNotificationToAllClients 向所有客户端发送通知
func (h *MCPHandler) SendNotificationToAllClients(method string, params map[string]any) {
	h.server.SendNotificationToAllClients(method, params)
}

// SendNotificationToSpecificClient 向特定客户端发送通知
func (h *MCPHandler) SendNotificationToSpecificClient(sessionID string, method string, params map[string]any) error {
	return h.server.SendNotificationToSpecificClient(sessionID, method, params)
}
