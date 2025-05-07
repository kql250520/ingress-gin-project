package configutils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// SpringConfig 定义配置结构体
type SpringConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	PropertySources []PropertySource `json:"propertySources"`
}

// PropertySource 定义配置源结构体
type PropertySource struct {
	Name   string         `json:"name"`
	Source map[string]any `json:"source"`
}

// From Spring Cloud Config Server 获取配置
func GetConfigFromSpringServer(configServerURL, application, profile, label string) (*SpringConfig, error) {
	url := fmt.Sprintf("%s/%s/%s/%s", configServerURL, application, profile, label)
	log.Printf("请求配置服务器的URL: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求配置服务器失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("从 SpringCloud Config-Server 获取配置结果: %s", string(body))

	var config SpringConfig
	err = json.Unmarshal(body, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
