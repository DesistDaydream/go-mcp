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
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
	// 添加工具
	mcp.AddTool(
		server,
		&mcp.Tool{Name: "greet", Description: "工具的描述性信息"}, // InputSchema 由该 MCP 库自动生成
		SayHi,
	)

	// 启动 MCP Server
	getServer := func(req *http.Request) *mcp.Server { return server }
	handler := mcp.NewStreamableHTTPHandler(getServer, nil)
	http.ListenAndServe(":8080", handler)
}
