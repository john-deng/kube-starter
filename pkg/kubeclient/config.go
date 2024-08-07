package kubeclient

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/kube-starter/pkg/kubeconfig"
	"github.com/hidevopsio/kube-starter/pkg/oidc"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Subject = "subject"
)

// NewKubeClient new kube client
func NewKubeClient(scheme *runtime.Scheme, cfg *RestConfig) (k8sClient client.Client, err error) {
	log.Info("Creating kube client")
	k8sClient, err = client.New(cfg.Config, client.Options{Scheme: scheme})
	log.Infof("created kube client: %v", k8sClient)
	return
}

// NewRuntimeKubeClient new runtime kube client
func NewRuntimeKubeClient(scheme *runtime.Scheme, token *oidc.Token, useToken bool, properties *Properties, cluster *kubeconfig.ClusterConfig) (cli client.Client, err error) {
	var cfg *rest.Config
	rcc := &kubeconfig.ClusterConfig{
		ClusterInfo: kubeconfig.ClusterInfo{
			Name:   cluster.Name,
			Config: cluster.Config,
		},
	}
	cfg, err = kubeconfig.Kubeconfig(&rcc.ClusterInfo)
	if err != nil {
		log.Warn(err)
		return
	}
	cfg.QPS = properties.QPS
	cfg.Burst = properties.Burst
	cfg.Timeout = properties.Timeout

	if token != nil && token.Claims != nil && token.Data != "" {
		kubeServiceHost := os.Getenv("KUBERNETES_SERVICE_HOST")
		if kubeServiceHost == "" && useToken {
			cfg.BearerToken = token.Data
			cfg.BearerTokenFile = ""
		} else {
			switch properties.OIDCScope {
			case "email":
				cfg.Impersonate.UserName = token.Claims.Email
			case "profile":
				cfg.Impersonate.UserName = token.Claims.Username
			case "openid":
				cfg.Impersonate.UserName = token.Claims.Issuer + "#" + token.Claims.Subject
			default:
				cfg.Impersonate.UserName = token.Claims.Email
			}
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

	log.Infof("created runtime client with qps: %v, burst: %v", cfg.QPS, cfg.Burst)
	return
}

func GetClusterConfig(clusterName string, token *oidc.Token, prop *Properties) (clusterConfig *kubeconfig.ClusterConfig) {
	clusterConfig = new(kubeconfig.ClusterConfig)

	if clusterName == "" {
		clusterName = "main"
	}

	clusterConfig.Config = prop.Clusters[clusterName].Config
	clusterConfig.Name = clusterName

	clusterConfig.Username = token.Claims.Username

	return
}
