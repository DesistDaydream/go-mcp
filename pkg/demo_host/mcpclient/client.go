package mcpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// 启动:
// set -a; source .env-mcp-grafana; set +a
// mcp-grafana -transport streamable-http -address 0.0.0.0:8080 -debug -log-level debug -enabled-tools rendering
func GrafanaMCP() {
	ctx := context.Background()

	client := mcp.NewClient(
		&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"},
		nil,
	)

	// 使用 "标准输入输出" 或者 "HTTP" 与 MCP Server 交互之前，先设置环境变量
	// set -a; source .env-mcp-grafana; set +a
	transport := &mcp.CommandTransport{Command: exec.Command("mcp-grafana")}
	session, err := client.Connect(ctx, transport, nil)
	// 通过 HTTP 连接到 MCP Server。需要指定 MCP Server 的 HTTP 地址。
	// mcp-grafana -transport streamable-http -address 0.0.0.0:8080 -debug -log-level debug -enabled-tools rendering
	// session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: "http://0.0.0.0:8080/mcp"}, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// 列出所有可用的工具
	tools, err := session.ListTools(ctx, nil)
	for _, tool := range tools.Tools {
		fmt.Println(tool.Name)
		fmt.Println(tool.Description)
		b, _ := json.MarshalIndent(tool.InputSchema, "", "  ")
		fmt.Println(string(b))
	}

	// {"dashboardUid":"f5377340-1e2a-4b93-80d3-906dd51fb2e4","height":600,"panelId":17,"theme":"dark","timeRange":{"from":"now-1h","to":"now"},"variables":{"var-custom_house_name":"$__all","var-datasource":"bei-jing","var-instance":"172.16.32.99"},"width":1200}
	args := map[string]any{
		"dashboardUid": "f5377340-1e2a-4b93-80d3-906dd51fb2e4",
		"panelId":      17,
		"width":        1200,
		"height":       600,
		"theme":        "dark",
		// "theme":        "light",
		"timeRange": map[string]any{
			"from": "now-1h",
			"to":   "now",
		},
		"variables": map[string]string{
			"var-datasource":        "bei-jing",
			"var-custom_house_name": "$__all",
			"var-instance":          "172.16.32.99",
		},
	}
	// 调用 MCP Server 中的工具
	params := &mcp.CallToolParams{
		Name:      "get_panel_image",
		Arguments: args,
	}
	res, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("调用工具失败: %v", err)
	}
	if res.IsError {
		log.Fatalf("工具执行失败: %v", res.Content[0].(*mcp.TextContent).Text)
	}
	for _, c := range res.Content {
		switch cType := c.(type) {
		case *mcp.TextContent:
			log.Print(cType.Text)
		case *mcp.ImageContent:
			log.Print(cType.MIMEType)
			// fmt.Println(string(cType.Data))
			// 将 cType.Data 写入文件
			if err := os.WriteFile("image.png", cType.Data, 0644); err != nil {
				log.Fatalf("写入文件失败: %v", err)
			}
		case *mcp.AudioContent:
			log.Print(cType.MIMEType)
			log.Print(string(cType.Data))
		default:
			log.Print("未知内容类型")
		}
	}
}
