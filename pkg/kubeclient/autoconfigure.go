// Work flow:
// 1. Create the kube client by the kube config
//   - kube-starter will use the default kube config under config/kubeclient/default.yaml
//   - if the config is empty, will use the default $HOME/.kube/config
//   - if the config is not empty, will use the config file

// Package kubeclient implement the kube client for the applications that use the API Server of kubernetes
package kubeclient

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/kube-starter/pkg/kubeconfig"
	"github.com/hidevopsio/kube-starter/pkg/oidc"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Profile = "kubeclient"
)

type RestConfig struct {
	at.Scope `value:"prototype"`

	*rest.Config
}

// Client is the encapsulation of the default kube client
type Client struct {
	at.Scope `value:"prototype"`

	client.Client
}

// ClientCreation is the encapsulation of the default kube client
type ClientCreation struct {
	at.Scope `value:"prototype"`

	client.Client
}

// RuntimeClient is the client the runtime kube client
type RuntimeClient struct {
	at.Scope `value:"prototype"`

	client.Client

	Context context.Context `json:"context"`
}

// RuntimeClientCreation is the client the runtime kube client
type RuntimeClientCreation struct {
	at.Scope `value:"prototype"`

	client.Client

	Context context.Context `json:"context"`
}

type configuration struct {
	at.AutoConfiguration

	Properties *Properties

	clientFactory        *instantiate.ScopedInstanceFactory[*ClientCreation]
	runtimeClientFactory *instantiate.ScopedInstanceFactory[*RuntimeClientCreation]
}

func newConfiguration(prop *Properties) *configuration {

	return &configuration{Properties: prop}
}

func init() {
	app.Register(
		newConfiguration,
		new(Properties),
		new(instantiate.ScopedInstanceFactory[*ClientCreation]),
		new(instantiate.ScopedInstanceFactory[*RuntimeClientCreation]),
	)
}

func (c *configuration) ClusterConfig(prop *Properties) (cluster *kubeconfig.ClusterConfig, err error) {
	clusterConfig := prop.Clusters["main"]
	cluster = &kubeconfig.ClusterConfig{
		ClusterInfo: kubeconfig.ClusterInfo{
			Name:   clusterConfig.Name,
			Config: clusterConfig.Config, // if config is empty, will use the default $HOME/.kube/config
		},
	}
	log.Infof("ClusterConfig: %+v", cluster)
	return
}

func (c *configuration) RestConfig(cluster *kubeconfig.ClusterConfig) (restConfig *RestConfig, err error) {
	restConfig = new(RestConfig)

	restConfig.Config, err = kubeconfig.Kubeconfig(&cluster.ClusterInfo)
	restConfig.Config.QPS = c.Properties.QPS
	restConfig.Config.Burst = c.Properties.Burst
	restConfig.Config.Timeout = c.Properties.Timeout
	return
}

func (c *configuration) ClientCreation(scheme *runtime.Scheme, cfg *RestConfig) (cli *ClientCreation, err error) {

	cli = &ClientCreation{}
	cli.Client, err = NewKubeClient(scheme, cfg)

	return
}

func (c *configuration) Client(
	cluster *kubeconfig.ClusterConfig,
	clientFactory *instantiate.ScopedInstanceFactory[*ClientCreation],
) (cli *Client, err error) {
	cli = new(Client)
	var kc *ClientCreation
	kc, err = clientFactory.GetInstance(cluster)
	if err == nil {
		cli.Client = kc.Client
	}

	return
}

func (c *configuration) RuntimeClientCreation(
	ctx context.Context,
	scheme *runtime.Scheme,
	token *oidc.Token,
	cluster *kubeconfig.ClusterConfig,
) (cli *RuntimeClientCreation, err error) {

	cli = new(RuntimeClientCreation)
	var newClient client.Client

	newClient, err = NewRuntimeKubeClient(scheme, token, true, c.Properties, cluster)
	if err != nil {
		log.Error(err)
		return
	}

	cli = &RuntimeClientCreation{
		Context: ctx,
		Client:  newClient,
	}

	return
}

func (c *configuration) RuntimeClient(
	ctx context.Context,
	token *oidc.Token,
	cluster *kubeconfig.ClusterConfig,
	runtimeClientFactory *instantiate.ScopedInstanceFactory[*RuntimeClientCreation],
) (cli *RuntimeClient, err error) {

	cli = new(RuntimeClient)
	cluster.Username = token.Claims.Username
	var rc *RuntimeClientCreation
	rc, err = runtimeClientFactory.GetInstance(ctx, cluster)
	if err == nil {
		cli.Context = ctx
		cli.Client = rc.Client
	}

	return
}
