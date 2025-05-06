package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 实例化 MCP Server
	mcpServer := server.NewMCPServer(
		"Demo",
		"1.0.0",
	)

	// 声明一个 Tool
	toolDemo := mcp.NewTool("MCP Server 的 Tool 示例",
		// Tool 的描述性信息
		mcp.WithDescription("这是一个 MCP Server 的 Tool 示例。可以返回传入的参数"),
		// Tool 的参数，类似命令行参数。这里定义了一个名为 name 的字符串类型参数，是必填的。
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("在这里定义传入的参数，该工具的逻辑会将参数与其他字符串组合作为返回值"),
		),
		// 如果 Tool 有多个参数，可以继续添加。
	)

	toolCurrentTime := mcp.NewTool("获取当前时间",
		mcp.WithDescription("获取指定时区的当前时间, 默认时区为 Asia/Shanghai"),
		mcp.WithString("timezone",
			mcp.Required(),
			mcp.Description("当前时间的时区"),
		),
	)

	// 将 Tools 添加到 MCP Server
	mcpServer.AddTool(toolDemo, demoHandler)
	mcpServer.AddTool(toolCurrentTime, currentTimeHandler)

	// 启动 MCP Server，通过标准输入/输出与 MCP Client 通信
	// TODO: 怎么让 Server 监听在 TCP 端口上？
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func demoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 从请求中获取参数
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return mcp.NewToolResultError("name 参数必须是字符串类型"), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf(`你好，%s`, name)), nil
}

func currentTimeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	timezone, ok := request.Params.Arguments["timezone"].(string)
	if !ok {
		return mcp.NewToolResultError("timezone 参数必须是字符串类型"), nil
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("解析时区异常: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf(`当前时间是: %s`, time.Now().In(loc))), nil
}
