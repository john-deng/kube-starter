package kubeclient

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/kube-starter/pkg/kubeconfig"
	"github.com/hidevopsio/kube-starter/pkg/oidc"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Profile = "kubeclient"
)

type clientCache struct {
	client client.Client
	uid string
	token string
}

type configuration struct {
	at.AutoConfiguration

	Properties *Properties

	clients cmap.ConcurrentMap
}

func newConfiguration() *configuration {
	return &configuration{clients: cmap.New()}
}

func init() {
	app.Register(newConfiguration)
}

func (c *configuration) RestConfig() (cfg *rest.Config, err error) {
	cfg, err = kubeconfig.Kubeconfig(c.Properties.DefaultInCluster)
	return
}

// Client is the encapsulation of the default kube client
type Client struct {
	//at.ContextAware

	client.Client

	//Context context.Context `json:"context"`
}

func (c *configuration) Client(scheme *runtime.Scheme, cfg *rest.Config) (cli *Client, err error) {

	cli = &Client{}
	cli.Client, err = KubeClient(scheme, cfg)

	return
}

// ImpersonateClient is the client the impersonate kube client
type ImpersonateClient struct {
	at.ContextAware

	client.Client

	Context context.Context `json:"context"`
}

func (c *configuration) ImpersonateClient(ctx context.Context, scheme *runtime.Scheme, token *oidc.Token) (cli *ImpersonateClient) {
	cli = new(ImpersonateClient)

	newCli, _ := RuntimeKubeClient(scheme, token, false, c.Properties.DefaultInCluster)

	cli = &ImpersonateClient{
		Context: ctx,
		Client: newCli,
	}
	return
}

// TokenizeClient is the client the tokenize kube client
type TokenizeClient struct {
	at.ContextAware

	client.Client

	Context context.Context `json:"context"`
}

func (c *configuration) TokenizeClient(ctx context.Context, scheme *runtime.Scheme, token *oidc.Token) (cli *TokenizeClient) {
	cli = new(TokenizeClient)

	newCli, _ := RuntimeKubeClient(scheme, token, true, c.Properties.DefaultInCluster)

	cli = &TokenizeClient{
		Context: ctx,
		Client: newCli,
	}
	return
}

// RuntimeClient is the client the runtime kube client
type RuntimeClient struct {
	at.ContextAware

	client.Client

	Context context.Context `json:"context"`
}

func (c *configuration) RuntimeClient(ctx context.Context, scheme *runtime.Scheme, token *oidc.Token) (cli *RuntimeClient, err error) {
	cli = new(RuntimeClient)
	var newClient client.Client
	var ok bool
	var cachedClient interface{}

	uid := token.Claims.Issuer + "#" + token.Claims.Subject
	cachedClient, ok = c.clients.Get(uid)
	if ok {
		cc := cachedClient.(clientCache)
		if cc.token == token.Data {
			newClient = cc.client
		}
	}

	if newClient == nil {
		newClient, err = RuntimeKubeClient(scheme, token, true, c.Properties.DefaultInCluster)
		if err != nil {
			return
		}

		c.clients.Set(uid, clientCache{client: newClient, uid: token.Claims.Username, token: token.Data})
	}

	cli = &RuntimeClient{
		Context: ctx,
		Client: newClient,
	}
	return
}
