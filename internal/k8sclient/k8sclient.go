package k8sclient

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewKubernetesClient 创建一个新的Kubernetes客户端
func NewKubernetesClient(k8url, k8token string) (*kubernetes.Clientset, error) {
	// 配置 Kubernetes 客户端
	config := &rest.Config{
		Host:        k8url,
		BearerToken: k8token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true, // 注意：生产环境中建议使用安全的 TLS 配置
		},
	}

	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create Kubernetes client: %v", err)
		return nil, err
	}

	return clientset, nil
}
