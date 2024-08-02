package kubeconfig

import (
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/utils/crypto/base64"
	"os"
	"path/filepath"

	"github.com/hidevopsio/hiboot/pkg/log"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterInfo struct {
	Name      string `json:"name"`
	Config    string `json:"config"`
	InCluster bool   `json:"inCluster"`
}

type ClusterConfig struct {
	at.Scope              `value:"prototype"`
	at.ConditionalOnField `value:"Name"`

	ClusterInfo
}

type RuntimeClusterConfig struct {
	at.Scope              `value:"prototype"`
	at.ConditionalOnField `value:"Name,Username"`

	ClusterInfo
	Username string `value:"username"`
}

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
func Kubeconfig(clusterInfo *ClusterInfo) (cfg *rest.Config, err error) {
	if clusterInfo.InCluster {
		cfg, err = rest.InClusterConfig()
	} else {
		if clusterInfo.Config != "" {
			decodedConfig, decodeErr := base64.Decode([]byte(clusterInfo.Config))
			if decodeErr != nil {
				err = decodeErr
				log.Warnf("Error decoding base64 kubeconfig: %s", err.Error())
				return
			}
			cfg, err = clientcmd.RESTConfigFromKubeConfig(decodedConfig)
			if err != nil {
				log.Warnf("Error building kubeconfig from decoded content: %s", err.Error())
			}
		} else {
			cfg, err = clientcmd.BuildConfigFromFlags("", DefaultKubeconfig())
			if err != nil {
				log.Warnf("Error building kubeconfig: %s", err.Error())
			}
		}
	}
	return
}
