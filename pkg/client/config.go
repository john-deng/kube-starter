package client

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DefaultKubeconfig
func DefaultKubeconfig() string {
	fname := os.Getenv("KUBECONFIG")
	if fname != "" {
		return fname
	}
	home, err := os.UserHomeDir()
	if err != nil {
		klog.Warningf("failed to get home directory: %v", err)
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}

// Kubeconfig
func Kubeconfig() (cfg *rest.Config, err error) {
	cfg, err = rest.InClusterConfig()
	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags("", DefaultKubeconfig())
		if err != nil {
			klog.Fatalf("Error building kubeconfig: %s", err.Error())
		}
	}
	return
}

// Client
func Client() (k8sClient client.Client, err error)  {
	var cfg *rest.Config
	cfg, err = Kubeconfig()
	if err == nil {
		k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	} else {
		klog.Fatal(err)
	}
	return
}
