package main

import (
	"context"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type Output struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

func SayHi(ctx context.Context, req *mcp.CallToolRequest, input Input) (*mcp.CallToolResult, Output, error) {
	return nil, Output{Greeting: "Hi " + input.Name}, nil
}

func main() {
	// 实例化 MCP Server
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
	// 添加工具
	// 参数 1: MCP Server 实例
	// 参数 2: 工具定义，包含工具的名称、描述、输入/输出参数、etc. 。这些信息会被转为 JSON Schema 定义，当 MCP Client 请求时告诉 Client。
	// 参数 3: 工具被调用时执行的具体逻辑
	mcp.AddTool(
		server,
		&mcp.Tool{Name: "greet", Description: "工具的描述性信息"}, // InputSchema 由该 MCP 库自动生成
		SayHi,
	)

	// 启动 MCP Server
	// 在 标准输入/标准输出 上运行 MCP Server。
	// 以二进制文件的形式提供给 MCP Client，Client 可以通过执行该文件来连接到 MCP Server。
	// if err := server.Run(
	// 	context.Background(),
	// 	&mcp.StdioTransport{},
	// ); err != nil {
	// 	log.Fatal(err)
	// }

	// 启动 MCP Server
	// 在 TCP 端口上运行 MCP Server
	// 将 IP:PORT 信息提供给 MCP Client，使用 IP:PORT 即可连接到 MCP Server。
	getServer := func(req *http.Request) *mcp.Server { return server }
	handler := mcp.NewStreamableHTTPHandler(getServer, nil)
	http.ListenAndServe(":8080", handler)
}
