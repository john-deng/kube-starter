package kubeclient

import (
	"os"

	"github.com/hidevopsio/kube-starter/pkg/kubeconfig"
	"github.com/hidevopsio/kube-starter/pkg/oidc"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/hidevopsio/hiboot/pkg/log"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Subject = "subject"
)

// KubeClient new kube client
func KubeClient(scheme *runtime.Scheme, cfg *rest.Config) (k8sClient client.Client, err error)  {
	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	return
}

// RuntimeKubeClient new runtime kube client
func RuntimeKubeClient(scheme *runtime.Scheme, token *oidc.Token, useToken bool, inCluster *bool) (cli client.Client, err error)  {
	var cfg *rest.Config
	cfg, err = kubeconfig.Kubeconfig(inCluster)
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
			cfg.Impersonate.UserName = token.Claims.Email
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
