package kubeclient

import (
	"os"
	"path/filepath"

	"github.com/hidevopsio/kube-starter/pkg/oidc"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/log"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Subject = "subject"
)

// DefaultKubeconfig
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

// Kubeconfig
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

// KubeClient
func KubeClient(scheme *runtime.Scheme) (k8sClient client.Client, err error)  {
	var cfg *rest.Config
	cfg, err = Kubeconfig()
	if err != nil {
		log.Warn(err)
		return
	}
	//cfg.Impersonate.UserName = ""
	//cfg.BearerToken = ""
	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	return
}

// RuntimeKubeClient
func RuntimeKubeClient(ctx context.Context, scheme *runtime.Scheme, token *oidc.Token, useToken bool) (cli client.Client, err error)  {
	var cfg *rest.Config
	cfg, err = Kubeconfig()
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
		// unauthorized user
		//ctx.StatusCode(http.StatusUnauthorized) no need to use it as middleware will handle it
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

