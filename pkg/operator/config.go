package operator

import (
	"os"
	"path/filepath"

	"github.com/hidevopsio/hiboot/pkg/log"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	Subject = "subject"
)

// DefaultKubeconfig load default kube config
func DefaultKubeconfig() string {
	fname := os.Getenv("KUBECONFIG")
	if fname != "" {
		return fname
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Warnf("failed to get home directory: %v", err)
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}

// Kubeconfig new kube config
func Kubeconfig() (cfg *rest.Config, err error) {
	cfg, err = rest.InClusterConfig()
	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags("", DefaultKubeconfig())
		if err != nil {
			log.Warnf("Error building kubeconfig: %s", err.Error())
		}
	}
	return
}
