package eureka

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/SimonWang00/goeureka"
)

type ResponseStruct struct {
	MASTER_NODE          []string
	API_VERSION          string
	KUBERNETES_URL       string
	KUBERNETES_TOKEN     string
	IS_EDGE_MODEL_DEPLOY bool
}

// 注册服务到Eureka,获取eureka中Region服务发起http请求
func RegisterToEureka() (ResponseStruct, error) {
	//第二项参数为空代表默认使用本机IP地址;nil代表opt不需要登录账户和密码注册
	goeureka.RegisterClient("http://10.181.12.168:30333", "", "MYAPP", "8080", "443", nil)

	//获取所有服务实例
	services, err := goeureka.GetServices()
	if err != nil {
		fmt.Println("Error getting services:", err)

	}

	var regionIpPort string
	for _, service := range services {
		name := service.Name
		instances := service.Instance
		//查找REGION服务
		if name == "REGION" {
			// 遍历 REGION 服务的所有信息
			for _, instance := range instances {
				regionIpPort = fmt.Sprintf("%s%s%d", instance.IpAddr, ":", instance.Port.Port) //:=如果这样就重新申明了变量region_url，会报declared and not used: region_url
				log.Printf("REGION service info regionUrl: %s", regionIpPort)
			}

		}

	}

	//regioUrl := regionIpPort + "/v1/regions/2ba1bb42-dc5b-ae3d-6fdb-e173dc69211d/kubernetes/"

	urlTest := "http://10.181.12.168:30503/v1/regions/2ba1bb42-dc5b-ae3d-6fdb-e173dc69211d/kubernetes/"
	resp2, err := http.Get(urlTest)
	if err != nil {
		return ResponseStruct{}, fmt.Errorf("failed to http.Get request: %w", err)
	}
	// 关闭响应体
	defer resp2.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp2.Body)
	if err != nil {
		return ResponseStruct{}, fmt.Errorf("failed to read response body: %w", err)
	}
	// 解析JSON响应体
	var response ResponseStruct
	if err := json.Unmarshal(body, &response); err != nil {
		return ResponseStruct{}, fmt.Errorf("failed to unmarshall json response: %w", err)
	}
	return response, nil

}
