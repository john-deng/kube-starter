package kubeclient

import (
	"os"
	"path/filepath"

	"github.com/hidevopsio/kube-starter/pkg/oidc"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/hidevopsio/hiboot/pkg/log"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
func Kubeconfig(inCluster *bool) (cfg *rest.Config, err error) {
	if inCluster == nil || *inCluster {
		cfg, err = rest.InClusterConfig()
	}

	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags("", DefaultKubeconfig())
		if err != nil {
			log.Warnf("Error building kubeconfig: %s", err.Error())
		}
	}
	return
}

// KubeClient new kube client
func KubeClient(scheme *runtime.Scheme, inCluster *bool) (k8sClient client.Client, err error)  {
	var cfg *rest.Config
	cfg, err = Kubeconfig(inCluster)
	if err != nil {
		log.Warn(err)
		return
	}
	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	return
}

// RuntimeKubeClient new runtime kube client
func RuntimeKubeClient(scheme *runtime.Scheme, token *oidc.Token, useToken bool, inCluster *bool) (cli client.Client, err error)  {
	var cfg *rest.Config
	cfg, err = Kubeconfig(inCluster)
	if err != nil {
		log.Warn(err)
		return
	}

	if token != nil && token.Claims != nil && token.Data != "" {
		kubeServiceHost := os.Getenv("KUBERNETES_SERVICE_HOST")
		if kubeServiceHost == "" && useToken {
			cfg.BearerToken = token.Data
			cfg.BearerTokenFile = ""
		} else {
			cfg.Impersonate.UserName = token.Claims.Issuer + "#" + token.Claims.Subject
		}
	} else {
		log.Warn("Unauthorized")
		err = errors.NewUnauthorized("Unauthorized")
		return
	}

	cli, err = client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		log.Warn(err)
	}
	return
}
