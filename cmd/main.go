package main

import (
	"log"

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
	r.GET("/delete-ingress", handlers.HandleDeleteIngress(k8sUrl, k8sToken, namespace, ingressName))
	r.Run(":8080")
}
