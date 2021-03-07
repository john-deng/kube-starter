package kubeclient

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
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

// Client
type Client client.Client
func (c *configuration) Client() (cli Client) {
	cli, _ = KubeClient()
	return
}
