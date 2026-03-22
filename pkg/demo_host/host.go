package main

import (
	"github.com/DesistDaydream/go-mcp/pkg/demo_host/mcpclient"
)

type MCPServerConfig struct {
	Enabled   bool              `json:"enabled"`   // 是否启用该服务器
	Name      string            `json:"name"`      // 唯一标识，作为 map key
	Transport string            `json:"transport"` // "stdio" 或 "sse" 或 "http"
	Command   string            `json:"command"`   // stdio 模式：可执行文件路径
	Args      []string          `json:"args"`      // stdio 模式：命令参数
	Env       map[string]string `json:"env"`       // 额外注入的环境变量
	URL       string            `json:"url"`       // http/sse 模式的远端地址
	Headers   map[string]string `json:"headers"`   // http/sse 模式的请求头
}

type MCPHostConifg struct {
	Servers map[string]MCPServerConfig `json:"servers"`
}

func main() {
	mcpclient.GrafanaMCP()
}
