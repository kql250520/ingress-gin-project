package main

import (
	"log"
	"net/http"

	"github.com/kql250520/ingress-gin-project/internal/configutils"
	"github.com/kql250520/ingress-gin-project/internal/eureka"
	"github.com/kql250520/ingress-gin-project/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 健康检查路由
	r.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Spring Cloud Config Server 的 URL
	configServerURL := "http://10.181.0.7:30001"
	// 配置文件名称
	application := "pkg2-ms-log"
	// 配置文件后缀
	profile := "prod"
	// git分支
	label := "miniocean-prod"
	r.GET("/config", func(c *gin.Context) {
		configresult, err := configutils.GetConfigFromSpringServer(configServerURL, application, profile, label)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, configresult)
	})

	// 注册到Eureka获取返回值
	retEurekaResponse, err := eureka.RegisterToEureka()
	if err != nil {
		log.Fatalf("Failed to register to Eureka: %v", err)
	}
	k8sUrl := retEurekaResponse.KUBERNETES_URL
	k8sToken := retEurekaResponse.KUBERNETES_TOKEN
	// 创建Ingress规则的路由
	r.GET("/create-ingress", handlers.HandleCreateIngress(k8sUrl, k8sToken))

	// 查询Ingress规则的路由
	namespace := "migu-p2c"
	ingressName := "test-ingress"
	r.GET("/get-ingress", handlers.HandleGetIngress(k8sUrl, k8sToken, namespace, ingressName))

	// 删除Ingress规则的路由
	r.GET("/delete-ingress", handlers.HandleDeleteIngress(k8sUrl, k8sToken, namespace, ingressName))

	r.Run(":8080")
}
