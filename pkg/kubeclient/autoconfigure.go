package kubeclient

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/kube-starter/pkg/oidc"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Profile = "kubeclient"
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

// ImpersonateClient
type Client struct {
	//at.ContextAware

	client.Client

	//Context context.Context `json:"context"`
}

func (c *configuration) Client(scheme *runtime.Scheme) (cli *Client) {

	newCli, _ := KubeClient(scheme)

	cli = &Client{
		Client: newCli,
	}

	return
}

// ImpersonateClient
type ImpersonateClient struct {
	at.ContextAware

	client.Client

	Context context.Context `json:"context"`
}

// ImpersonateClient
func (c *configuration) ImpersonateClient(ctx context.Context, scheme *runtime.Scheme, token *oidc.Token) (cli *ImpersonateClient) {
	cli = new(ImpersonateClient)

	newCli, _ := RuntimeKubeClient(ctx, scheme, token, false)

	cli = &ImpersonateClient{
		Context: ctx,
		Client: newCli,
	}
	return
}

// TokenizeClient
type TokenizeClient struct {
	at.ContextAware

	client.Client

	Context context.Context `json:"context"`
}

// TokenizeClient
func (c *configuration) TokenizeClient(ctx context.Context, scheme *runtime.Scheme, token *oidc.Token) (cli *TokenizeClient) {
	cli = new(TokenizeClient)

	newCli, _ := RuntimeKubeClient(ctx, scheme, token, true)

	cli = &TokenizeClient{
		Context: ctx,
		Client: newCli,
	}
	return
}
