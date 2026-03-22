package mcpclient_test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func getPanelImage(
	grafanaURL string,
	token string,
	dashboardUID string,
	panelID int,
	width, height int,
	theme string,
	from, to string,
	variables map[string]string,
	// 多值 instance 单独传，支持多个
	instances []string,
) ([]byte, error) {
	// Grafana render 路径
	renderPath := fmt.Sprintf("/render/d-solo/%s", dashboardUID)

	q := url.Values{}
	q.Set("panelId", fmt.Sprintf("%d", panelID))
	q.Set("width", fmt.Sprintf("%d", width))
	q.Set("height", fmt.Sprintf("%d", height))
	q.Set("theme", theme)
	q.Set("from", from)
	q.Set("to", to)

	// 普通变量（单值）
	for k, v := range variables {
		q.Set(k, v)
	}

	// 多值变量：var-instance 重复添加
	for _, inst := range instances {
		q.Add("var-instance", inst)
	}

	fullURL := strings.TrimRight(grafanaURL, "/") + renderPath + "?" + q.Encode()

	fmt.Println(fullURL)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	fmt.Println(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("grafana render failed: status=%d body=%s", resp.StatusCode, body)
	}

	return io.ReadAll(resp.Body)
}

func TestGrafanaMCP(t *testing.T) {
	grafanaURL := os.Getenv("GRAFANA_URL")
	token := os.Getenv("GRAFANA_SERVICE_ACCOUNT_TOKEN")
	pngData, err := getPanelImage(
		grafanaURL,
		token,
		"f5377340-1e2a-4b93-80d3-906dd51fb2e4",
		17,
		1200, 600,
		"dark",
		"now-1h", "now",
		map[string]string{
			"var-datasource":        "bei-jing",
			"var-custom_house_name": "$__all",
		},
		[]string{"172.16.31.129", "172.16.32.99"},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("got PNG, size=%d bytes\n", len(pngData))
	// 后续可写入文件或转 base64
	if err := os.WriteFile("panel.png", pngData, 0644); err != nil {
		panic(err)
	}
}
