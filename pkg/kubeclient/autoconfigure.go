package kubeclient

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type configuration struct {
	at.AutoConfiguration

}

func newConfiguration() *configuration {
	return &configuration{}
}

func init() {
	app.Register(newConfiguration)
}

// RestConfig
type RestConfig struct {
	*rest.Config
}

// RestConfig
func (c *configuration) RestConfig(scheme *runtime.Scheme) (cfg *RestConfig) {
	cfg = new(RestConfig)
	cfg.Config, _ = Kubeconfig()
	return
}

// Client
type Client client.Client
func (c *configuration) Client(scheme *runtime.Scheme, cfg *RestConfig) (cli Client) {
	cli, _ = KubeClient(scheme, cfg)
	return
}

// RuntimeClient
type RuntimeClient struct {
	at.ContextAware

	client.Client
}

// RuntimeClient
func (c *configuration) RuntimeClient(ctx context.Context, scheme *runtime.Scheme, cfg *RestConfig) (cli *RuntimeClient) {
	newCli, _ := RuntimeKubeClient(ctx, scheme, cfg)
	cli = &RuntimeClient{
		Client: newCli,
	}
	return
}
