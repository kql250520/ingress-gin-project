package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/kql250520/ingress-gin-project/internal/k8sclient"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 处理CreateIngress
func HandleCreateIngress(k8url, k8token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//调用k8sclient方法，新建NewKubernetes客户端
		clientset, err := k8sclient.NewKubernetesClient(k8url, k8token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Kubernetes client"})
			return
		}

		ingressClassName := "nginx-02"
		// 构建Ingress规则
		ingress := &v1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "test-ingress",
				Namespace:   "migu-p2c",
				Annotations: map[string]string{"nginx.ingress.kubernetes.io/rewrite-target": "/"},
			},
			Spec: v1.IngressSpec{
				IngressClassName: &ingressClassName,
				Rules: []v1.IngressRule{
					{
						//	Host: "*",
						IngressRuleValue: v1.IngressRuleValue{
							HTTP: &v1.HTTPIngressRuleValue{
								Paths: []v1.HTTPIngressPath{
									{
										Path:     "/nginx-test",
										PathType: &[]v1.PathType{v1.PathTypePrefix}[0],
										Backend: v1.IngressBackend{
											Service: &v1.IngressServiceBackend{
												Name: "nginx-test-service",
												Port: v1.ServiceBackendPort{
													Number: 80,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		result, err := CreateIngress(clientset, "migu-p2c", ingress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ingress"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Ingress created successfully", "ingress": result})
	}
}

// 创建Ingress
func CreateIngress(clientset *kubernetes.Clientset, namespace string, ingress *v1.Ingress) (*v1.Ingress, error) {
	ingressClient := clientset.NetworkingV1().Ingresses(namespace)
	result, err := ingressClient.Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Failed to create ingress: %v", err)
		return nil, err
	}
	return result, nil
}

// 处理GetIngress
func HandleGetIngress(k8url, k8token, namespace, ingressName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//调用k8sclient方法，新建NewKubernetes客户端
		clientset, err := k8sclient.NewKubernetesClient(k8url, k8token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Kubernetes client"})
			return
		}

		result, err := GetIngress(clientset, namespace, ingressName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ingress"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Ingress created successfully", "ingress": result})
	}
}

// 查询Ingress
func GetIngress(clientset *kubernetes.Clientset, namespace, ingressName string) (*v1.Ingress, error) {
	ingressClient := clientset.NetworkingV1().Ingresses(namespace)
	result, err := ingressClient.Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get ingress: %v", err)
		return nil, err
	}
	return result, nil
}
