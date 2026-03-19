package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	// 实例化 MCP Client
	// 使用默认的选项
	client := mcp.NewClient(
		&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"},
		nil,
	)

	// 通过 标准输入/标准输出 连接到 MCP Server。需要指定 MCP Server 二进制文件的路径（在 exec.Command 参数中）。
	// transport := &mcp.CommandTransport{Command: exec.Command("bin/myserver")}
	// session, err := client.Connect(ctx, transport, nil)
	// 通过 HTTP 连接到 MCP Server。需要指定 MCP Server 的 HTTP 地址。
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: "http://localhost:8080/"}, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// 列出所有可用的工具
	tools, err := session.ListTools(ctx, nil)
	for _, tool := range tools.Tools {
		fmt.Println(tool.Name)
		fmt.Println(tool.Description)
		// InputSchema 保存了要使用工具时需要传递的参数的 JSON Schema 定义
		b, _ := json.MarshalIndent(tool.InputSchema, "", "  ")
		fmt.Println(string(b))
	}

	// 调用 MCP Server 中的工具
	params := &mcp.CallToolParams{
		Name:      "greet",
		Arguments: map[string]any{"name": "you"},
	}
	res, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("调用工具失败: %v", err)
	}
	if res.IsError {
		log.Fatal("工具执行失败")
	}
	for _, c := range res.Content {
		log.Print(c.(*mcp.TextContent).Text)
	}
}
