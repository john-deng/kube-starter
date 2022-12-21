package kubeclient

import (
	"errors"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrNilKubeClient = errors.New("kube client is nil, please check if API Server is available")
)

type middleware struct {
	at.Middleware
}

func newMiddleware() *middleware {
	return &middleware{}
}

func init() {
	app.Register(newMiddleware)
}

// CheckKubeClient is the middleware handler,it supports dependency injection, method annotation
// middleware handler can be annotated to specific purpose or general purpose
func (m *middleware) CheckKubeClient(_ struct {
	at.MiddlewareHandler `value:"/" `
},
	kubeclient *Client,
	scheme *runtime.Scheme,
	cfg *rest.Config,
	ctx context.Context) (err error) {

	if kubeclient.Client == nil {
		kubeclient.Client, err = client.New(cfg, client.Options{Scheme: scheme})
		if err == nil && kubeclient.Client != nil {
			app.Register(kubeclient)
			log.Infof("Got kube client by retry %v", kubeclient)
		} else {
			log.Warn(err)
			ctx.StatusCode(500)
			ctx.ResponseBody(err.Error(), nil)
			return
		}
	}
	log.Debug("Got kube client from middleware")
	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}
